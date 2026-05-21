package db

import (
    "context"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
    Pool *pgxpool.Pool
}

func Connect(ctx context.Context, databaseURL string) (*Postgres, error) {
    if databaseURL == "" {
        return nil, fmt.Errorf("DATABASE_URL is required")
    }

    cfg, err := pgxpool.ParseConfig(databaseURL)
    if err != nil {
        return nil, fmt.Errorf("parse database config: %w", err)
    }

    cfg.MaxConns = 10
    cfg.MinConns = 1
    cfg.MaxConnLifetime = time.Hour
    cfg.MaxConnIdleTime = 30 * time.Minute
    cfg.HealthCheckPeriod = time.Minute

    pool, err := pgxpool.NewWithConfig(ctx, cfg)
    if err != nil {
        return nil, fmt.Errorf("connect postgres: %w", err)
    }

    if err := pool.Ping(ctx); err != nil {
        pool.Close()
        return nil, fmt.Errorf("ping postgres: %w", err)
    }

    return &Postgres{Pool: pool}, nil
}

func (p *Postgres) Close() {
    if p != nil && p.Pool != nil {
        p.Pool.Close()
    }
}
