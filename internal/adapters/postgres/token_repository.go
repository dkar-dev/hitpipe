package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dkar-dev/hitpipe/internal/domain/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VerificationTokenRepository struct {
	pool *pgxpool.Pool
}

func NewVerificationTokenRepository(pool *pgxpool.Pool) *VerificationTokenRepository {
	return &VerificationTokenRepository{pool}
}

func (r *VerificationTokenRepository) CreateToken(ctx context.Context, userID uuid.UUID, token string, expires time.Time) error {
	id := uuid.New()
	_, err := r.pool.Exec(ctx, `
		INSERT INTO email_verification_tokens
		VALUES ($1, $2, $3, $4)
	`, id, userID, token, expires)

	if err != nil {
		return err
	}

	if pgErr := new(pgconn.PgError); errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return user.ErrTokenAlreadyExists
		}
	}

	return nil
}

func (r *VerificationTokenRepository) GetUserIDByToken(ctx context.Context, token string) (uuid.UUID, error) {

	row := r.pool.QueryRow(ctx, `
		SELECT id, userid, expires_at from email_verification_tokens WHERE token=$1
		`, token)

	var id, userID uuid.UUID
	var expiresAt time.Time

	err := row.Scan(
		&id,
		&userID,
		&expiresAt)

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("scan values in token structure: %w", err)
	}

	if expiresAt.Before(time.Now()) {
		return uuid.UUID{}, fmt.Errorf("token expired")
	}

	return userID, nil
}

func (r *VerificationTokenRepository) DeleteToken(ctx context.Context, token string) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM email_verification_tokens WHERE token=$1`, token)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("token not found")
	}
	return nil
}
