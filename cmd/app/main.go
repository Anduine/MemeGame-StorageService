package main

import (
	"log/slog"
	"os"
	"storage-service/internal/app"
	"storage-service/internal/config"
	"storage-service/internal/lib/logger"
)

func main() {
	config := config.MustLoadConfig()

	logger.InitGlobalLogger(os.Stdout, slog.LevelDebug)

	app.Run(config)
}
