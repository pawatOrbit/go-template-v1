package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ExampleRepository demonstrates how to implement a repository in this template
// This is just an example - replace with your actual data access interfaces
type ExampleRepository interface {
	GetExampleByID(ctx context.Context, id string) (*ExampleData, error)
	CreateExample(ctx context.Context, data *ExampleData) error
	UpdateExample(ctx context.Context, id string, data *ExampleData) error
	DeleteExample(ctx context.Context, id string) error
}

// ExampleData represents data structure for examples
type ExampleData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type exampleRepositoryImpl struct {
	readPgPool  *pgxpool.Pool
	writePgPool *pgxpool.Pool
}

// NewExampleRepository creates a new example repository
func NewExampleRepository(readPgPool *pgxpool.Pool, writePgPool *pgxpool.Pool) ExampleRepository {
	return &exampleRepositoryImpl{
		readPgPool:  readPgPool,
		writePgPool: writePgPool,
	}
}

// GetExampleByID retrieves an example by its ID
func (r *exampleRepositoryImpl) GetExampleByID(ctx context.Context, id string) (*ExampleData, error) {
	// Example implementation - replace with your actual SQL queries
	// You can use sqlc generated code here like:
	// qtx := db_sqlc.New(r.readPgPool)
	// result, err := qtx.GetExampleByID(ctx, id)
	
	// For now, return mock data
	return &ExampleData{
		ID:          id,
		Name:        "Example Item",
		Description: "This is an example from the template",
		CreatedAt:   "2024-01-01T00:00:00Z",
		UpdatedAt:   "2024-01-01T00:00:00Z",
	}, nil
}

// CreateExample creates a new example in the database
func (r *exampleRepositoryImpl) CreateExample(ctx context.Context, data *ExampleData) error {
	// Example implementation - replace with your actual SQL queries
	// You can use sqlc generated code here like:
	// qtx := db_sqlc.New(r.writePgPool)
	// err := qtx.CreateExample(ctx, db_sqlc.CreateExampleParams{...})
	
	// For now, just return nil (success)
	return nil
}

// UpdateExample updates an existing example in the database
func (r *exampleRepositoryImpl) UpdateExample(ctx context.Context, id string, data *ExampleData) error {
	// Example implementation - replace with your actual SQL queries
	// qtx := db_sqlc.New(r.writePgPool)
	// err := qtx.UpdateExample(ctx, db_sqlc.UpdateExampleParams{...})
	
	return nil
}

// DeleteExample removes an example from the database
func (r *exampleRepositoryImpl) DeleteExample(ctx context.Context, id string) error {
	// Example implementation - replace with your actual SQL queries
	// qtx := db_sqlc.New(r.writePgPool)
	// err := qtx.DeleteExample(ctx, id)
	
	return nil
}