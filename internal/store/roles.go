package store

import (
	"context"
	"database/sql"
)

type Roles struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Description string `json:"description"`
}

type RolesStore struct {
	db *sql.DB
}

func (s *RolesStore) GetByName(ctx context.Context, roleName string) (*Roles, error) {
	query := `SELECT id, name, level, description FROM roles WHERE name = $1`

	var role Roles
	err := s.db.QueryRowContext(ctx, query, roleName).Scan(&role.ID, &role.Name, &role.Level, &role.Description)

	if err != nil {
		return nil, err
	}

	return &role, nil
}
