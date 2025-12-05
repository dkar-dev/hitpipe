package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/dkar-dev/hitpipe/internal/adapters/notifiers"
	"github.com/dkar-dev/hitpipe/internal/domain/user"
	"github.com/dkar-dev/hitpipe/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo  ports.UserRepository
	tokenRepo ports.VerificationTokenRepository
	notifier  ports.Notifier
	log       *slog.Logger
}

func NewUserService(r ports.UserRepository, t ports.VerificationTokenRepository, n ports.Notifier, l *slog.Logger) *UserService {
	return &UserService{userRepo: r, tokenRepo: t, notifier: n, log: l}
}

func (s *UserService) Register(ctx context.Context, email, password string) (*user.UserAggregate, error) {

	res, err := s.userRepo.FindByEmail(ctx, email)
	if res != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}
	if !errors.Is(err, user.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("generating password hash: %w", err)
	}

	var aggregate = user.NewLocalUserAggregate(email, string(passwordHash))

	err = s.userRepo.Save(ctx, aggregate)
	if err != nil {
		return nil, err
	}

	token, err := user.GenerateVerificationToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token: %w", err)
	}

	exp := time.Now().Add(20 * time.Minute)
	err = s.tokenRepo.CreateToken(ctx, aggregate.User.ID, token, exp)
	if err != nil {
		return nil, fmt.Errorf("failed to create verification token: %w", err)
	}

	req := notifiers.NewWelcomeSendRequest(email, token)

	go func() {
		err := s.notifier.Send(context.TODO(), *req)
		if err != nil {
			log.Println("fail during the message sending", "error", err)
		}
	}()

	return aggregate, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*user.UserAggregate, error) {
	res, err := s.userRepo.FindByEmail(ctx, email)
	if errors.Is(err, user.ErrUserNotFound) {
		return nil, fmt.Errorf("user with email %s not found", email)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if !res.Identities.IsVerified {
		return nil, fmt.Errorf("user with email %s is not verified", email)
	}

	err = bcrypt.CompareHashAndPassword([]byte(res.Identities.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	return res, nil
}

func (s *UserService) VerifyEmail(ctx context.Context, token string) error {
	userID, err := s.tokenRepo.GetUserIDByToken(ctx, token)
	if err != nil {
		return err
	}

	err = s.userRepo.MarkVerified(ctx, userID)
	if err != nil {
		return err
	}

	_ = s.tokenRepo.DeleteToken(ctx, token)

	return nil
}
