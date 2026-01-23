package database

import (
	_ "github.com/mattn/go-sqlite3"
)

// # NOT THIS PROJECT POINT, REMOVE IT FOR NOW
// func (s *SQLite) ResetUsers(ctx context.Context) error {
// 	_, err := s.db.ExecContext(ctx, `
// 	UPDATE users
// 	SET dismiss = 1
//   `)
// 	return err
// }
