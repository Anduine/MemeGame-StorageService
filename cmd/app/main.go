package main

import (
	"os"
	"storage-service/internal/app"
	"storage-service/internal/config"
	"storage-service/internal/lib/logger"
)

func main() {
	config := config.MustLoadConfig()

	logger.SetupPlusLogger(os.Stdout)

	app.Run(config.Port, config.Timeout)
}