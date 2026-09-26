package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	ConnString        string
	MaxConns          int32
	MinConns          int32
	MaxConnIdleTime   time.Duration
	MaxConnLifetime   time.Duration
	HealthCheckPeriod time.Duration
}

func NewPool(ctx context.Context, cfg DBConfig) (*pgxpool.Pool, error) {
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
