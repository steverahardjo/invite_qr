// Package config handles environment initialization including database
// connection setup, connection pooling, and configuration parsing.
package config

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

// Config holds all configurable application parameters including database
// connection details, SSL settings, and connection pool limits.
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

	AdminPasswordHash string
}

// DB wraps a *sql.DB connection for use throughout the application.
type DB struct {
	Conn *sql.DB
}

// NewDBFromEnv opens a PostgreSQL connection using the pgx driver via
// database/sql, configures connection pool parameters from the Config struct,
// and verifies connectivity with a ping. Returns the DB wrapper or an error.
func NewDBFromEnv(
	ctx context.Context,
	cfg *Config,
	log *zap.Logger,
) (*DB, error) {

	conn, err := sql.Open("pgx", dbDSN(cfg))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if cfg.PoolMaxConnections > 0 {
		conn.SetMaxOpenConns(cfg.PoolMaxConnections)
	}
	if cfg.PoolMinConnections > 0 {
		conn.SetMaxIdleConns(cfg.PoolMinConnections)
	}
	if cfg.PoolMaxConnLife > 0 {
		conn.SetConnMaxLifetime(cfg.PoolMaxConnLife)
	}
	if cfg.PoolMaxConnIdle > 0 {
		conn.SetConnMaxIdleTime(cfg.PoolMaxConnIdle)
	}

	if err := conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("database connection initialized")
	return &DB{Conn: conn}, nil
}

// Close shuts down the underlying database connection pool.
func (db *DB) Close(log *zap.Logger) {
	log.Info("closing database connection")
	db.Conn.Close()
}

// dbDSN builds a space-separated key=value DSN string from the Config struct.
func dbDSN(cfg *Config) string {
	vals := dbValues(cfg)

	parts := make([]string, 0, len(vals))
	for k, v := range vals {
		parts = append(parts, k+"="+v)
	}

	return strings.Join(parts, " ")
}

// setIfNotEmpty adds the key-value pair to the map only if value is non-empty.
func setIfNotEmpty(
	p map[string]string,
	key string,
	value string,
) {
	if value != "" {
		p[key] = value
	}
}

// setIfPositive adds the key-value pair as a string only if value > 0.
func setIfPositive(
	p map[string]string,
	key string,
	value int,
) {
	if value > 0 {
		p[key] = strconv.Itoa(value)
	}
}

// setIfPositiveDuration adds the key-value pair as milliseconds only if value > 0.
func setIfPositiveDuration(
	p map[string]string,
	key string,
	value time.Duration,
) {
	if value > 0 {
		p[key] = strconv.FormatInt(value.Milliseconds(), 10)
	}
}

// dbValues builds a map of libpq-compatible connection parameters from the Config.
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

	return p
}
