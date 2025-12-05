package user

import (
	"github.com/google/uuid"
)

func NewLocalUserAggregate(email string, passwordHash string) *UserAggregate {
	userId := uuid.New()

	return &UserAggregate{
		User: User{
			ID:     userId,
			Email:  email,
			Status: "pending_verification",
		},
		Identities: AuthIdentities{
			ID:             uuid.New(),
			Provider:       "local",
			ProviderUserID: "email",
			PasswordHash:   passwordHash,
		},
	}
}
