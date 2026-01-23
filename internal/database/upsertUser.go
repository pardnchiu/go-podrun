package database

import (
	_ "github.com/mattn/go-sqlite3"
)

// # NOT THIS PROJECT POINT, REMOVE IT FOR NOW
// func (s *SQLite) UpsertUser(ctx context.Context, d *model.User) error {
// 	_, err := s.db.ExecContext(ctx, `
//   INSERT INTO users (email)
//   VALUES (?)
//   ON CONFLICT(email) DO UPDATE SET
//     dismiss = 0
//   `,
// 		d.Email,
// 	)
// 	return err
// }
