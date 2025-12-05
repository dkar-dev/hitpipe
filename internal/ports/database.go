package ports

import (
	"context"
	"time"

	"github.com/dkar-dev/hitpipe/internal/domain/user"
	"github.com/google/uuid"
)

type UserRepository interface {
	Save(ctx context.Context, agg *user.UserAggregate) error
	//FindByID(ctx context.Context, id uuid.UUID) (*domain.UserAggregate, error)
	FindByEmail(ctx context.Context, email string) (*user.UserAggregate, error)
	MarkVerified(ctx context.Context, userID uuid.UUID) error
}

type VerificationTokenRepository interface {
	CreateToken(ctx context.Context, userID uuid.UUID, token string, expires time.Time) error
	GetUserIDByToken(ctx context.Context, token string) (uuid.UUID, error)
	DeleteToken(ctx context.Context, token string) error
}
