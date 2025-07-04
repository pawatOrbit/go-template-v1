package repository

import (
	"fmt"
	"log/slog"

	"github.com/pawatOrbit/ai-mock-data-service/go/core/pgdb"
)

type Repository struct{
	TableSchemasRepository TableSchemasRepository
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
		TableSchemasRepository: NewTableSchemasRepository(readPgPool, writePgPool),
	}, nil
}
