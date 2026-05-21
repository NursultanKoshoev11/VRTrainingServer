package main

import (
    "log/slog"
    "os"

    "github.com/NursultanKoshoev11/VRTrainingServer/internal/config"
    "github.com/NursultanKoshoev11/VRTrainingServer/internal/httpserver"
)

func main() {
    cfg := config.Load()
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel()}))

    srv := httpserver.New(cfg, logger)

    logger.Info("starting VRTrainingServer", "address", cfg.HTTPAddress())
    if err := srv.Start(); err != nil {
        logger.Error("server stopped", "error", err)
        os.Exit(1)
    }
}
