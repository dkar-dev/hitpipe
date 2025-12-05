package user

import (
	"time"

	"github.com/google/uuid"
)

type UserAggregate struct {
	User       User           `json:"user"`
	Identities AuthIdentities `json:"identities"`
}

type User struct {
	ID        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type AuthIdentities struct {
	ID             uuid.UUID `db:"id"`
	Provider       string    `db:"provider"`
	ProviderUserID string    `db:"provider_user_id"`
	PasswordHash   string    `db:"password_hash"`
	IsVerified     bool      `db:"is_verified"`
	CreatedAt      time.Time `db:"created_at"`
}
