package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgConf struct {
	User              string
	Password          string
	Host              string
	DBname            string
	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
}

type PostgresDB struct {
	pool *pgxpool.Pool
}

func NewPgStorage(ctx context.Context, cfg *PgConf) (*PostgresDB, error) {
	config, err := pgxpool.ParseConfig(getConnStr(cfg))
	if err != nil {
		return &PostgresDB{}, fmt.Errorf("Unable to parse config: %v", err)
	}

	config.MaxConns = cfg.MaxConns
	config.MinConns = cfg.MinConns
	config.MaxConnLifetime = cfg.MaxConnLifetime
	config.MaxConnIdleTime = cfg.MaxConnIdleTime
	config.HealthCheckPeriod = cfg.HealthCheckPeriod

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return &PostgresDB{}, fmt.Errorf("Unable to create connection pool: %v", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return &PostgresDB{}, fmt.Errorf("Unable to ping database: %v", err)
	}

	return &PostgresDB{
		pool: pool,
	}, nil
}

func getConnStr(cfg *PgConf) string {
	return fmt.Sprintf("postgres://%s:%s@%s:5432/%s", cfg.User, cfg.Password, cfg.Host, cfg.DBname)
}

func (s *PostgresDB) Close() {
	s.pool.Close()
}

func (s *PostgresDB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return s.pool.Query(ctx, sql, args...)
}

func (s *PostgresDB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return s.pool.QueryRow(ctx, sql, args...)
}

func (s *PostgresDB) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return s.pool.Exec(ctx, sql, args...)
}
