package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud_storage/config"
	pgadapter "cloud_storage/internal/adapter/postgres"
	"cloud_storage/internal/adapter/storage"
	httpcontroller "cloud_storage/internal/controller/http"
	v1 "cloud_storage/internal/controller/http/v1"
	"cloud_storage/internal/usecase"
	"cloud_storage/pkg/httpserver"
	pkgjwt "cloud_storage/pkg/jwt"
	"cloud_storage/pkg/logger"
	pkgpostgres "cloud_storage/pkg/postgres"

	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func Run() error {
	// 1. Config
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	// 2. Logger
	log := logger.New(cfg.LogLevel)
	log.Info("starting cloud-storage")

	// 3. Postgres pool
	ctx := context.Background()
	pool, err := pkgpostgres.New(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("postgres: %w", err)
	}
	defer pool.Close()
	log.Info("connected to postgres")

	// 4. Migrations
	if err = runMigrations(cfg.Postgres, log); err != nil {
		return fmt.Errorf("migrations: %w", err)
	}

	// 5. Adapters
	pgRepo := pgadapter.New(pool)

	localStore, err := storage.New(cfg.Storage.BasePath)
	if err != nil {
		return fmt.Errorf("storage: %w", err)
	}

	// 6. JWT Manager
	jwtManager := pkgjwt.New(cfg.JWT.Secret, cfg.JWT.TTL)

	// 7. UseCase
	uc := usecase.New(pgRepo, pgRepo, pgRepo, localStore, jwtManager)

	// 8. HTTP handler & router
	handler := v1.New(uc, localStore, log)
	router := httpcontroller.NewRouter(handler, jwtManager, log)

	// 9. HTTP server
	srv := httpserver.New(router, cfg.HTTP)

	// 10. Start server
	serverErr := make(chan error, 1)
	go func() {
		log.Info("http server listening",
			"addr", fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port),
		)
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	// 11. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err = <-serverErr:
		return fmt.Errorf("server error: %w", err)
	case sig := <-quit:
		log.Info("shutting down", "signal", sig.String())
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	log.Info("server stopped gracefully")
	return nil
}

func runMigrations(cfg config.PostgresConfig, log *slog.Logger) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	poolCfg, _ := pgxpool.ParseConfig(dsn)
	db := stdlib.OpenDB(*poolCfg.ConnConfig)
	defer db.Close()

	driver, err := migratepostgres.WithInstance(db, &migratepostgres.Config{})
	if err != nil {
		return fmt.Errorf("migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migration/postgres", "postgres", driver)
	if err != nil {
		return fmt.Errorf("migrate instance: %w", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}

	log.Info("migrations applied")
	return nil
}
