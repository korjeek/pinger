package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/korjeek/pinger/backend/config"
)

func NewPool(ctx context.Context, cfg config.DBConfig) (*pgxpool.Pool, error) {
	dbCfg, err := pgxpool.ParseConfig(cfg.ConnString)
	if err != nil {
		//TODO: handling error
		return nil, err
	}

	dbCfg.MaxConns = cfg.MaxConns
	dbCfg.MinConns = cfg.MinConns
	dbCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	dbCfg.MaxConnLifetime = cfg.MaxConnLifetime
	dbCfg.HealthCheckPeriod = cfg.HealthCheckPeriod

	pool, err := pgxpool.NewWithConfig(ctx, dbCfg)
	if err != nil {
		//TODO: handling error
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		//TODO: handling error
		return nil, err
	}

	return pool, nil
}
