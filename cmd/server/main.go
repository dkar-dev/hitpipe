package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dkar-dev/hitpipe/internal/adapters"
	"github.com/dkar-dev/hitpipe/internal/adapters/notifiers"
	"github.com/dkar-dev/hitpipe/internal/adapters/postgres"
	"github.com/dkar-dev/hitpipe/internal/config"
	"github.com/dkar-dev/hitpipe/internal/service"
	"github.com/dkar-dev/hitpipe/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v4"
)

func main() {

	// Loading configuration from config.yaml and .env file
	// Implemented in loader.go
	cfg, err := config.Load("./config/")
	if err != nil {
		log.Fatalf("ERROR: failed to load config.yaml file: %v", err)
	}

	log := logger.NewLogger(cfg.App.Env, cfg.Logger.Level)
	slog.SetDefault(log)

	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer cancel()

	pool, err := initDatabase(ctx, cfg)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	userRepository := postgres.NewUserRepository(pool)
	tokenRepo := postgres.NewVerificationTokenRepository(pool)
	emailNotifier := notifiers.NewBetaEmailNotifier(
		cfg.Notifier.EmailAdapter.From,
		cfg.Notifier.EmailAdapter.Password,
		cfg.Notifier.EmailAdapter.Host,
		cfg.Notifier.EmailAdapter.Port)

	userService := service.NewUserService(userRepository, tokenRepo, emailNotifier, log)
	userAPI := adapters.NewUserAPI(userService, log)

	e := echo.New()

	go func() {
		<-ctx.Done()
		log.Info("shutting down the server")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			log.Error("failed to shutdown the server", "error", err)
		}
	}()

	e.POST("/register", userAPI.Register)
	e.POST("/login", userAPI.Login)
	e.GET("/verify", userAPI.VerifyEmail)

	err = e.Start(":" + cfg.App.Port)
	if err != nil {
		log.Error("failed to start the server", "error", err)
		return
	}
}

func initDatabase(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DB,
		cfg.Postgres.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}
