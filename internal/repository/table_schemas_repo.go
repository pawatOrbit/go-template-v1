package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	db_sqlc "github.com/pawatOrbit/ai-mock-data-service/go/internal/sqlc/db"
)

type TableSchemasRepository interface {
	GetDatabaseSchemaTableNames(ctx context.Context) ([]string, error)
	GetDatabaseSchemaByTableName(ctx context.Context, tableName string) (db_sqlc.DatabaseSchema, error)
}

type TableSchemasRepositoryImpl struct {
	readPgPool  *pgxpool.Pool
	writePgPool *pgxpool.Pool
}

func NewTableSchemasRepository(readPgPool *pgxpool.Pool, writePgPool *pgxpool.Pool) TableSchemasRepository {
	return &TableSchemasRepositoryImpl{
		readPgPool:  readPgPool,
		writePgPool: writePgPool,
	}
}

func (r *TableSchemasRepositoryImpl) GetDatabaseSchemaTableNames(ctx context.Context) ([]string, error) {
	qtx := db_sqlc.New(r.readPgPool)
	tableNames, err := qtx.GetDatabaseSchemaTableNames(ctx)
	if err != nil {
		return nil, err
	}
	return tableNames, nil
}

func (r *TableSchemasRepositoryImpl) GetDatabaseSchemaByTableName(ctx context.Context, tableName string) (db_sqlc.DatabaseSchema, error) {
	qtx := db_sqlc.New(r.readPgPool)
	schema, err := qtx.GetDatabaseSchemaByTableName(ctx, tableName)
	if err != nil {
		return db_sqlc.DatabaseSchema{}, err
	}
	return schema, nil
}
