package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/pardnchiu/go-podrun/internal/command"
	"github.com/pardnchiu/go-podrun/internal/utils"
)

func main() {
	if err := utils.CheckRelyPackages(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	args, err := command.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	slog.Info("parsed args", slog.Any("args", args))
}
