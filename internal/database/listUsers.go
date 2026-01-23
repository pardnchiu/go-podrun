package database

import (
	_ "github.com/mattn/go-sqlite3"
)

// # NOT THIS PROJECT POINT, REMOVE IT FOR NOW
// func (s *SQLite) ListUsers(ctx context.Context) ([]model.User, error) {
// 	rows, err := s.db.QueryContext(ctx, `
// 	SELECT email, create_at
// 	FROM users
// 	WHERE dismiss = 0
// 	`)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var results []model.User
// 	for rows.Next() {
// 		var c model.User
// 		if err := rows.Scan(
// 			&c.Email,
// 			&c.CreatedAt,
// 		); err != nil {
// 			return nil, err
// 		}
// 		results = append(results, c)
// 	}

// 	return results, rows.Err()
// }
