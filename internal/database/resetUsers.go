package database

import (
	"context"

	_ "github.com/mattn/go-sqlite3"
)

func (s *SQLite) ResetUsers(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
	UPDATE users
	SET dismiss = 1
  `)
	return err
}
