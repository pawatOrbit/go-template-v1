package pgdb

import (
	"context"
	"fmt"
	"sync"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	readPgPool  *pgxpool.Pool
	writePgPool *pgxpool.Pool
	m           sync.Mutex
)

type Postgres struct {
	Read  PostgresConfig `mapstructure:"read"`
	Write PostgresConfig `mapstructure:"write"`
}

type PostgresConfig struct {
	Host                     string `mapstructure:"host"`
	Port                     int    `mapstructure:"port"`
	Username                 string `mapstructure:"username"`
	Password                 string `mapstructure:"password"`
	Database                 string `mapstructure:"database"`
	Schema                   string `mapstructure:"schema"`
	MaxConnections           int32  `mapstructure:"maxConnections"`
	EnableQueryParamsTracing bool   `mapstructure:"enableQueryParamsTracing"`
}

func InitPgConnectionPool(ctx context.Context, cfg Postgres) error {
	m.Lock()
	defer m.Unlock()

	// If no read config is provided
	// OR both configs point to the same host,
	// use a single connection pool.
	if cfg.Read.Host == "" || cfg.Read.Host == cfg.Write.Host {
		singlePool, err := initSinglePool(ctx, cfg.Write)
		if err != nil {
			return err
		}

		readPgPool = singlePool
		writePgPool = singlePool
		return nil
	}

	// For distinct read/write databases
	readPool, err := initSinglePool(ctx, cfg.Read)
	if err != nil {
		return err
	}
	writePool, err := initSinglePool(ctx, cfg.Write)
	if err != nil {
		return err
	}

	readPgPool = readPool
	writePgPool = writePool
	return nil
}

func GetReadPgPool() (*pgxpool.Pool, error) {
	if readPgPool == nil {
		return nil, fmt.Errorf("readPgPool is nil")
	}
	return readPgPool, nil
}

func GetWritePgPool() (*pgxpool.Pool, error) {
	if writePgPool == nil {
		return nil, fmt.Errorf("writePgPool is nil")
	}
	return writePgPool, nil
}

// initSinglePool initializes a single pool without acquiring a lock
func initSinglePool(ctx context.Context, postgresConfig PostgresConfig) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s search_path=%s",
		postgresConfig.Host,
		postgresConfig.Port,
		postgresConfig.Username,
		postgresConfig.Password,
		postgresConfig.Database,
		postgresConfig.Schema,
	)

	connConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		fmt.Println("Failed to parse config:", err)
		return nil, err
	}

	opts := []otelpgx.Option{}
	//if postgresConfig.EnableQueryParamsTracing {
	opts = append(opts, otelpgx.WithIncludeQueryParameters())
	//}

	connConfig.ConnConfig.Tracer = otelpgx.NewTracer(opts...)

	// Set maximum number of connections
	connConfig.MaxConns = postgresConfig.MaxConnections

	pgxPool, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		return nil, err
	}

	// Test the connection
	conn, err := pgxPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	return pgxPool, nil
}

func InitSchema(ctx context.Context, writePgPool *pgxpool.Pool, schema string) (err error) {
	// Create schema if it doesn't exist
	// Ignore error if schema already exists or if the user doesn't have permission to create schema
	_, err = writePgPool.Exec(
		ctx,
		fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s`, schema),
	)
	if err != nil {
		return err
	}

	return nil
}

func InitExtentions(ctx context.Context, writePgPool *pgxpool.Pool, schema string) (err error) {
	// Create extensions if they don't exist
	// Ignore error if extension already exists or if the user doesn't have permission to do so
	sql := `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp" schema pg_catalog;
	`
	_, err = writePgPool.Exec(ctx, sql)
	if err != nil {
		return err
	}

	return nil
}

func ClosePgPool() {
	m.Lock()
	defer m.Unlock()

	if readPgPool != nil {
		readPgPool.Close()
		readPgPool = nil
	}

	if writePgPool != nil {
		writePgPool.Close()
		writePgPool = nil
	}
}
