package repository

import (
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/go-api-template/core/pgdb"
)

type Repository struct{
	DB                *pgxpool.Pool // For health checks and other operations
	
	// Example repositories - replace with your actual repositories
	ExampleRepository ExampleRepository
}

func NewRepository() (*Repository, error) {
	readPgPool, err := pgdb.GetReadPgPool()
	if err != nil {
		return nil, fmt.Errorf("error getting read pool: %w", err)
	}

	writePgPool, err := pgdb.GetWritePgPool()
	if err != nil {
		return nil, fmt.Errorf("error getting write pool: %w", err)
	}

	slog.Info("Repository initialized", "readPgPool", readPgPool!=nil, "writePgPool", writePgPool!=nil)
	// Initialize all repositories here
	return &Repository{
		DB: readPgPool, // Use read pool for health checks
		
		// Example repositories - replace with your actual repositories
		ExampleRepository: NewExampleRepository(readPgPool, writePgPool),
	}, nil
}
