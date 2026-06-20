// file to setup all necessary environment init: db, observability, web launcher
package config

import (
	"context"
	"strconv"
	"strings"
	"time"

	pgxpool "github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Config struct {
	Name              string
	User              string
	Host              string
	Port              string
	SSLMode           string
	ConnectionTimeout int
	Password          string

	SSLCertPath     string
	SSLKeyPath      string
	SSLRootCertPath string

	PoolMinConnections int
	PoolMaxConnections int

	PoolMaxConnLife time.Duration
	PoolMaxConnIdle time.Duration
	PoolHealthCheck time.Duration
}

// DB wraps the pgx connection pool.
type DB struct {
	Pool *pgxpool.Pool
}

func NewDBFromEnv(
	ctx context.Context,
	cfg *Config,
	log *zap.Logger,
) (*DB, error) {

	pgxConfig, err := pgxpool.ParseConfig(dbDSN(cfg))
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, err
	}
	log.Info("database pool initialized")
	return &DB{
		Pool: pool,
	}, nil
}

func (db *DB) Close(log *zap.Logger) {
	log.Info("closing database pool")
	db.Pool.Close()
}

func dbDSN(cfg *Config) string {
	vals := dbValues(cfg)

	parts := make([]string, 0, len(vals))
	for k, v := range vals {
		parts = append(parts, k+"="+v)
	}

	return strings.Join(parts, " ")
}

func setIfNotEmpty(
	p map[string]string,
	key string,
	value string,
) {
	if value != "" {
		p[key] = value
	}
}

func setIfPositive(
	p map[string]string,
	key string,
	value int,
) {
	if value > 0 {
		p[key] = strconv.Itoa(value)
	}
}

func setIfPositiveDuration(
	p map[string]string,
	key string,
	value time.Duration,
) {
	if value > 0 {
		p[key] = strconv.FormatInt(value.Milliseconds(), 10)
	}
}

func dbValues(cfg *Config) map[string]string {
	p := map[string]string{}

	setIfNotEmpty(p, "dbname", cfg.Name)
	setIfNotEmpty(p, "user", cfg.User)
	setIfNotEmpty(p, "host", cfg.Host)
	setIfNotEmpty(p, "port", cfg.Port)
	setIfNotEmpty(p, "sslmode", cfg.SSLMode)

	setIfPositive(
		p,
		"connect_timeout",
		cfg.ConnectionTimeout,
	)

	setIfNotEmpty(p, "password", cfg.Password)
	setIfNotEmpty(p, "sslcert", cfg.SSLCertPath)
	setIfNotEmpty(p, "sslkey", cfg.SSLKeyPath)
	setIfNotEmpty(p, "sslrootcert", cfg.SSLRootCertPath)

	setIfPositive(
		p,
		"pool_min_conns",
		cfg.PoolMinConnections,
	)

	setIfPositive(
		p,
		"pool_max_conns",
		cfg.PoolMaxConnections,
	)

	setIfPositiveDuration(
		p,
		"pool_max_conn_lifetime",
		cfg.PoolMaxConnLife,
	)

	setIfPositiveDuration(
		p,
		"pool_max_conn_idle_time",
		cfg.PoolMaxConnIdle,
	)

	setIfPositiveDuration(
		p,
		"pool_health_check_period",
		cfg.PoolHealthCheck,
	)

	return p
}
