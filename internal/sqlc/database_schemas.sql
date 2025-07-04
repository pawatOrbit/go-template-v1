CREATE TABLE database_schemas (
    table_name varchar(100) PRIMARY KEY,
    table_script TEXT
);

-- name: GetDatabaseSchemaByTableName :one
SELECT * FROM database_schemas WHERE table_name = @table_name;

-- name: GetDatabaseSchemaTableNames :many
SELECT table_name FROM database_schemas;