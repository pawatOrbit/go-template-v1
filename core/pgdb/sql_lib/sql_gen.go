package sqllib

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

// GenerateSelect generates a SELECT SQL query.
func GenerateSelect(table string, columns []string, conditions map[string]interface{}, logicalOperators []string) (string, pgx.NamedArgs, bool) {
	sql := fmt.Sprintf("SELECT %s FROM %s", strings.Join(columns, ", "), table)
	args := pgx.NamedArgs{}
	conditionStr := []string{}

	i := 0
	for key, value := range conditions {
		var condition string
		switch v := value.(type) {
		case string:
			if strings.Contains(v, "()") {
				condition = fmt.Sprintf("%s = %s", key, v)
			} else {
				condition = fmt.Sprintf("%s = @%s", key, key)
				args[key] = v
			}
		case []interface{}:
			placeholders := []string{}
			for j, val := range v {
				placeholder := fmt.Sprintf("@%s_%d", key, j)
				placeholders = append(placeholders, placeholder)
				args[fmt.Sprintf("%s_%d", key, j)] = val
			}
			condition = fmt.Sprintf("%s IN (%s)", key, strings.Join(placeholders, ", "))
		default:
			condition = fmt.Sprintf("%s = @%s", key, key)
			args[key] = v
		}

		if i > 0 && len(logicalOperators) > 0 {
			conditionStr = append(conditionStr, logicalOperators[i-1])
		}
		conditionStr = append(conditionStr, condition)
		i++
	}

	if len(conditionStr) > 0 {
		sql += " WHERE " + strings.Join(conditionStr, " ")
	}

	return sql, args, false

}

// GenerateInsert generates an INSERT SQL query.
func GenerateInsert(table string, data map[string]interface{}) (string, pgx.NamedArgs, bool) {
	columns := []string{}
	values := []string{}
	args := pgx.NamedArgs{}

	for key, value := range data {
		columns = append(columns, key)
		values = append(values, fmt.Sprintf("@%s", key))
		args[key] = value
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(columns, ", "), strings.Join(values, ", "))

	return sql, args, true
}

// GenerateUpdate generates an UPDATE SQL query.
func GenerateUpdate(table string, data map[string]interface{}, conditions map[string]interface{}, logicalOperators []string) (string, pgx.NamedArgs, bool) {
	setStr := []string{}
	args := pgx.NamedArgs{}

	for key, value := range data {
		setStr = append(setStr, fmt.Sprintf("%s = @%s", key, key))
		args[key] = value
	}

	conditionStr := []string{}
	i := 0
	for key, value := range conditions {
		var condition string
		switch v := value.(type) {
		case string:
			if strings.Contains(v, "()") {
				condition = fmt.Sprintf("%s = %s", key, v)
			} else {
				condition = fmt.Sprintf("%s = @%s", key, key)
				args[key] = v
			}
		case []interface{}:
			placeholders := []string{}
			for j, val := range v {
				placeholder := fmt.Sprintf("@%s_%d", key, j)
				placeholders = append(placeholders, placeholder)
				args[fmt.Sprintf("%s_%d", key, j)] = val
			}
			condition = fmt.Sprintf("%s IN (%s)", key, strings.Join(placeholders, ", "))
		default:
			condition = fmt.Sprintf("%s = @%s", key, key)
			args[key] = v
		}

		if i > 0 && len(logicalOperators) > 0 {
			conditionStr = append(conditionStr, logicalOperators[i-1])
		}
		conditionStr = append(conditionStr, condition)
		i++
	}

	sql := fmt.Sprintf("UPDATE %s SET %s", table, strings.Join(setStr, ", "))
	if len(conditionStr) > 0 {
		sql += " WHERE " + strings.Join(conditionStr, " ")
	}

	return sql, args, true
}

// GenerateDelete generates a DELETE SQL query.
func GenerateDelete(table string, conditions map[string]interface{}, logicalOperators []string) (string, pgx.NamedArgs, bool) {
	sql := fmt.Sprintf("DELETE FROM %s", table)
	args := pgx.NamedArgs{}
	conditionStr := []string{}

	i := 0
	for key, value := range conditions {
		var condition string
		switch v := value.(type) {
		case string:
			if strings.Contains(v, "()") {
				condition = fmt.Sprintf("%s = %s", key, v)
			} else {
				condition = fmt.Sprintf("%s = @%s", key, key)
				args[key] = v
			}
		case []interface{}:
			placeholders := []string{}
			for j, val := range v {
				placeholder := fmt.Sprintf("@%s_%d", key, j)
				placeholders = append(placeholders, placeholder)
				args[fmt.Sprintf("%s_%d", key, j)] = val
			}
			condition = fmt.Sprintf("%s IN (%s)", key, strings.Join(placeholders, ", "))
		default:
			condition = fmt.Sprintf("%s = @%s", key, key)
			args[key] = v
		}

		if i > 0 && len(logicalOperators) > 0 {
			conditionStr = append(conditionStr, logicalOperators[i-1])
		}
		conditionStr = append(conditionStr, condition)
		i++
	}

	if len(conditionStr) > 0 {
		sql += " WHERE " + strings.Join(conditionStr, " ")
	}

	return sql, args, true
}
