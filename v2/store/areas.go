package store

import (
	"context"
	"github.com/jackc/pgx/v4"
)

var (
	// createAreasTableSQL is an SQL statement to create the areas table
	createAreasTableSQL = "CREATE TABLE IF NOT EXISTS areas (code VARCHAR (50) PRIMARY KEY NOT NULL, name VARCHAR (100) NOT NULL);"

	// insertAreaSQL is an SQL query to insert a new area - requires area code and name.
	insertAreaSQL = "INSERT INTO areas (code, name) VALUES ($1, $2) RETURNING code;"
)

// NewArea insert a new area, returns the area code.
func (s *AreaProfileStore) AddArea(code, name string) (string, error) {
	var areaCode string
	err := s.conn.QueryRow(context.Background(), insertAreaSQL, code, name).Scan(&areaCode)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return areaCode, nil
}
