package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"

	"os"
	"url_shortener_svc/internal/config"
	"url_shortener_svc/internal/lib/logger/handlers/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log, err := setupLogger(cfg.Env)
	if err != nil {
		fmt.Printf("failed to setup logger: %v", err)
		os.Exit(1)
	}

	log.Info("starting url shortener service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	//storage, err := sqlite.New(cfg.StoragePath)
	//if err != nil {
	//	log.Error("failed to init storage", sl.Err(err))
	//	os.Exit(1)
	//}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	//router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
}

func setupLogger(env string) (*slog.Logger, error) {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		return nil, fmt.Errorf("invalid environment: %s", env)
	}

	return log, nil
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
