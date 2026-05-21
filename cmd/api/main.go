package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/NursultanKoshoev11/VRTrainingServer/internal/config"
	"github.com/NursultanKoshoev11/VRTrainingServer/internal/db"
	"github.com/NursultanKoshoev11/VRTrainingServer/internal/httpserver"
	"github.com/NursultanKoshoev11/VRTrainingServer/internal/training"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel()}))

	database, err := db.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	if err := db.RunMigrations(context.Background(), database.Pool, cfg.MigrationsPath); err != nil {
		logger.Error("database migrations failed", "error", err)
		os.Exit(1)
	}
	logger.Info("database migrations applied")

	trainingRepo := training.NewPostgresRepository(database.Pool)
	trainingService := training.NewService(trainingRepo)

	srv := httpserver.New(cfg, logger, trainingService)

	logger.Info("starting VRTrainingServer", "address", cfg.HTTPAddress())
	if err := srv.Start(); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
