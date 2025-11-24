// Package main запускает HTTP-сервис назначения ревьюеров для Pull Request'ов.
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dixitix/pr-reviewer-service/internal/config"
	httpserver "github.com/dixitix/pr-reviewer-service/internal/http"
	"github.com/dixitix/pr-reviewer-service/internal/logger"
	"github.com/dixitix/pr-reviewer-service/internal/repository/postgres"
	"github.com/dixitix/pr-reviewer-service/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// main - точка входа в сервис назначения ревьюеров.
func main() {
	log := logger.New()

	if err := run(log); err != nil {
		log.Error("application exited with error", slog.Any("err", err))
		os.Exit(1)
	}
}

// run настраивает конфиг, подключение к БД, сервисы и HTTP-сервер и
// блокируется до завершения контекста (SIGINT/SIGTERM).
func run(log *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	db, err := newDB(ctx, cfg, log)
	if err != nil {
		return fmt.Errorf("init db: %w", err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			log.Error("failed to close db", slog.Any("err", cerr))
		}
	}()

	teamRepo := postgres.NewTeamRepository(db)
	userRepo := postgres.NewUserRepository(db)
	prRepo := postgres.NewPullRequestRepository(db)

	svc := service.NewService(teamRepo, userRepo, prRepo)

	httpHandler := httpserver.NewHandler(svc, log.With("layer", "http"))

	mux := http.NewServeMux()
	httpHandler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:         cfg.HTTP.Addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("http server starting", slog.String("addr", cfg.HTTP.Addr))

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server failed", slog.Any("err", err))
			stop()
		}
	}()

	<-ctx.Done()

	shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Info("http server shutting down")

	if err := srv.Shutdown(shutDownCtx); err != nil {
		return fmt.Errorf("http server shutdown: %w", err)
	}

	log.Info("application stopped cleanly")

	return nil
}

// newDB создаёт и настраивает пул подключений к БД и проверяет соединение.
func newDB(ctx context.Context, cfg config.Config, log *slog.Logger) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DB.DSN)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.DB.ConnMaxLifetime)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		log.Error("failed to ping database", slog.Any("err", err))

		return nil, fmt.Errorf("ping db: %w", err)
	}

	log.Info("database connection established")

	return db, nil
}
