package main

import (
	"log/slog"
	"os"
	"storage-service/internal/app"
	"storage-service/internal/config"
	"storage-service/internal/lib/logger"
)

func main() {
	logger.InitGlobalLogger(os.Stdout, slog.LevelDebug)

	config := config.MustLoadConfig()

	app.Run(config)
}
