package sqllib

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pawatOrbit/ai-mock-data-service/go/core/pgdb"
)

func Execute[R any](dbModel R, query string, args pgx.NamedArgs, isQueryWrite bool) ([]R, *int, error) {
	var dbPool *pgxpool.Pool
	var err error

	if isQueryWrite {
		dbPool, err = pgdb.GetWritePgPool()
	} else {
		dbPool, err = pgdb.GetReadPgPool()
	}

	if err != nil {
		return nil, nil, fmt.Errorf("error getting database pool: %w", err)
	}

	if dbPool == nil {
		return nil, nil, fmt.Errorf("dbPool is nil")
	}

	// Create a context with a timeout to avoid long-running queries
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if isQueryWrite {
		// Execute the query for write operations (INSERT, UPDATE, DELETE)
		rowsEffective, err := dbPool.Exec(ctx, query, args)
		if err != nil {
			return nil, nil, fmt.Errorf("error executing query: %w", err)
		}

		rows := int(rowsEffective.RowsAffected())
		// Check the number of rows affected
		if rows == 0 {
			slog.InfoContext(ctx, "No rows affected")
		} else {
			slog.InfoContext(ctx, "Rows affected", slog.Any("rowsAffected", rowsEffective.RowsAffected()))
		}

		return nil, &rows, nil
	}

	// Execute the query for read operations
	rows, err := dbPool.Query(ctx, query, args)
	if err != nil {
		return nil, nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	// Process the result
	var result []R
	result, err = pgx.CollectRows(rows, pgx.RowToStructByNameLax[R])
	if err != nil {
		return nil, nil, fmt.Errorf("error processing rows: %w", err)
	}

	rowLen := len(result)

	return result, &rowLen, nil
}
