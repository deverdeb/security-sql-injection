package db

import (
	"database/sql"
	"fmt"
)

// ExecuteQuery function executes a SQL query on database and converts query (with 'converterFunc' parameter function).
func ExecuteQuery[T any](converterFunc func(rows *sql.Rows) ([]*T, error), sqlRequestTemplate string, args ...any) ([]*T, error) {
	conn := GetConnection()
	statement, err := conn.Prepare(sqlRequestTemplate)
	if err != nil {
		return nil, fmt.Errorf("query preparation failed: %s\nerror: %w", sqlRequestTemplate, err)
	}
	row, err := statement.Query(args...)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %s\nerror: %w", sqlRequestTemplate, err)
	}
	defer row.Close()
	return converterFunc(row)
}

// ExecuteCreate function executes a SQL creation query on database and return new id.
func ExecuteCreate(sqlRequestTemplate string, args ...any) (int64, error) {
	conn := GetConnection()
	statement, err := conn.Prepare(sqlRequestTemplate)
	if err != nil {
		return 0, fmt.Errorf("query preparation failed: %s\nerror: %w", sqlRequestTemplate, err)
	}
	result, err := statement.Exec(args...)
	if err != nil {
		return 0, fmt.Errorf("query execution failed: %s\nerror: %w", sqlRequestTemplate, err)
	}
	return result.LastInsertId()
}

// Execute function executes a SQL query on database.
func Execute(sqlRequestTemplate string, args ...any) error {
	conn := GetConnection()
	statement, err := conn.Prepare(sqlRequestTemplate)
	if err != nil {
		return fmt.Errorf("query preparation failed: %s\nerror: %w", sqlRequestTemplate, err)
	}
	_, err = statement.Exec(args...)
	if err != nil {
		return fmt.Errorf("query execution failed: %s\nerror: %w", sqlRequestTemplate, err)
	}
	return nil
}
