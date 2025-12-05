// services/User/internal/adapters/postgres/user_repository.go
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/dkar-dev/hitpipe/internal/domain/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Save(ctx context.Context, agg *user.UserAggregate) error {

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	_, err = tx.Exec(ctx, `
		INSERT INTO users (id, email)
		VALUES ($1, $2)
		`, agg.User.ID, agg.User.Email)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("failed to insert into users: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO user_auth_identities (id, user_id, provider, provider_user_id, password_hash)
		VALUES ($1, $2, $3, $4, $5)
		`,
		agg.Identities.ID,
		agg.User.ID,
		agg.Identities.Provider,
		agg.Identities.ProviderUserID,
		agg.Identities.PasswordHash)

	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("failed to insert into user_auth_identitie: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

//
//func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
//
//	row := r.pool.QueryRow(ctx, `
//		SELECT email, status, created_at, updated_at FROM users WHERE id = $1
//		`, id)
//
//	var u domain.User
//
//	err := row.Scan(
//		&u.Email,
//		&u.Status,
//		&u.CreatedAt,
//		&u.UpdatedAt,
//	)
//	if err != nil {
//		if errors.Is(err, pgx.ErrNoRows) {
//			return nil, domain.ErrUserNotFound
//		}
//	}
//
//	return &u, nil
//}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.UserAggregate, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, status, created_at, updated_at FROM users WHERE email = $1
		`, email)

	var a = user.UserAggregate{
		User: user.User{
			Email: email,
		},
		Identities: user.AuthIdentities{
			Provider:       "local",
			ProviderUserID: "email",
		},
	}

	err := row.Scan(
		&a.User.ID,
		&a.User.Status,
		&a.User.CreatedAt,
		&a.User.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	row = r.pool.QueryRow(ctx, `
		SELECT id, password_hash, is_verified, created_at FROM user_auth_identities WHERE user_id = $1
		`, a.User.ID)

	err = row.Scan(
		&a.Identities.ID,
		&a.Identities.PasswordHash,
		&a.Identities.IsVerified,
		&a.Identities.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no users_auth_identies found with email %s: %w", email, err)
		}
	}

	return &a, nil
}

func (r *UserRepository) MarkVerified(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE user_auth_identities 
		SET is_verified=true 
		WHERE user_id=$1
		`, userID)

	if err != nil {
		return fmt.Errorf("error marking user verified: %w", err)
	}

	return nil
}
