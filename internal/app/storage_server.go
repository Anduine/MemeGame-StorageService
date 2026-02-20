package app

import (
	"storage-service/internal/config"
	"storage-service/internal/delivery/http_handlers"
	"storage-service/internal/server"
	"storage-service/internal/service"
)

func Run(cfg *config.Config) {
	storageService := service.NewStorageService(cfg.StoragePath)
	storageHandler := http_handlers.NewImagesHandler(storageService)

	handler := server.NewRouter(storageHandler)

	server.StartServer(handler, cfg.Port, cfg.Timeout)
}
