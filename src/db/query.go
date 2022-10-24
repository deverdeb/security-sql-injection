package db

import (
    "database/sql"
    "fmt"
)

// ExecuteQuery lance une requête SQL et retourne la liste des résultats.
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
