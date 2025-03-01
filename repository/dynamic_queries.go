package repository

import (
	"context"
	"fmt"
)

func (q *Queries) CreateTags(ctx context.Context, db DBTX, tags []string) ([]int64, error) {
	// Dynamically build the query
	query := "INSERT INTO tags (name) VALUES "
	var values []interface{}
	for i, tag := range tags {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("($%d)", i+1)
		values = append(values, tag)
	}
	query += "ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id"

	// Execute the query
	rows, err := db.Query(ctx, query, values...)
	if err != nil {
		return nil, fmt.Errorf("failed to insert tags: %w", err)
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan ID: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}
