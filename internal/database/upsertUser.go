package database

import (
	"context"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pardnchiu/go-podrun/internal/model"
)

func (s *SQLite) UpsertUser(ctx context.Context, d *model.User) error {
	slog.Info("info", slog.Any("data", d))
	_, err := s.db.ExecContext(ctx, `
  INSERT INTO users (email)
  VALUES (?)
  ON CONFLICT(email) DO UPDATE SET
    dismiss = 0
  `,
		d.Email,
	)
	return err
}
