package config

import (
    "log/slog"
    "os"
    "strings"
)

type Config struct {
    AppEnv             string
    ServerPort         string
    DatabaseURL        string
    JWTSecret          string
    CORSAllowedOrigins string
    ReportStoragePath  string
    LogLevelValue      string
}

func Load() Config {
    return Config{
        AppEnv:             getEnv("APP_ENV", "development"),
        ServerPort:         getEnv("SERVER_PORT", "8080"),
        DatabaseURL:        getEnv("DATABASE_URL", ""),
        JWTSecret:          getEnv("JWT_SECRET", "development-only-change-me"),
        CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
        ReportStoragePath:  getEnv("REPORT_STORAGE_PATH", "./reports"),
        LogLevelValue:      getEnv("LOG_LEVEL", "info"),
    }
}

func (c Config) HTTPAddress() string {
    return ":" + c.ServerPort
}

func (c Config) IsProduction() bool {
    return strings.EqualFold(c.AppEnv, "production")
}

func (c Config) LogLevel() slog.Level {
    switch strings.ToLower(c.LogLevelValue) {
    case "debug":
        return slog.LevelDebug
    case "warn":
        return slog.LevelWarn
    case "error":
        return slog.LevelError
    default:
        return slog.LevelInfo
    }
}

func getEnv(key string, fallback string) string {
    value := strings.TrimSpace(os.Getenv(key))
    if value == "" {
        return fallback
    }
    return value
}
