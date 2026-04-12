package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/periasamy/school-api/internal/config"
	"github.com/periasamy/school-api/internal/db"
	"github.com/periasamy/school-api/internal/server"
	"github.com/periasamy/school-api/pkg/logger"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGINT,
	)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		os.Stderr.WriteString("failed to load config: " + err.Error() + "\n")
		return err
	}

	log := logger.New(cfg.LogLevel)
	log.Info().Str("env", cfg.Environment).Msg("school-api starting")

	gormDB, err := db.NewDB(cfg.PostgresDSN(), cfg.LogLevel, log)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to database")
		return err
	}
	defer func() {
		sqlDB, err := gormDB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}()

	if err = server.Start(ctx, cfg, gormDB, log); err != nil {
		log.Error().Err(err).Msg("server error")
		return err
	}

	log.Info().Msg("school-api stopped")
	return nil
}
