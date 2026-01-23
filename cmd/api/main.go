package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pardnchiu/go-podrun/internal/database"
	"github.com/pardnchiu/go-podrun/internal/model"
	"github.com/pardnchiu/go-podrun/internal/routes"
)

func init() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("failed to load .env",
			slog.String("error", err.Error()))
	}
}

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		folderPath := filepath.Join(home, ".podrun")
		os.MkdirAll(folderPath, 0755)
		dbPath = filepath.Join(folderPath, "database.db")
	}
	db, err := database.NewSQLite(dbPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if emails := os.Getenv("ALLOW_EMAILS"); emails != "" {
		ctx := context.Background()
		for email := range strings.SplitSeq(emails, ",") {
			email = strings.TrimSpace(email)
			if email == "" {
				continue
			}
			if err := db.UpsertUser(ctx, &model.User{Email: email}); err != nil {
				fmt.Printf("failed to insert email %s: %v\n", email, err)
			}
		}
	}

	if err := routes.New(db); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
