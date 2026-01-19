package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pardnchiu/go-podrun/internal/database"
	"github.com/pardnchiu/go-podrun/internal/routes"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		// 開發環境：使用工作目錄
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dbPath = filepath.Join(wd, "podrun.db")
	}

	db, err := database.NewSQLite(dbPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := routes.New(db); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
