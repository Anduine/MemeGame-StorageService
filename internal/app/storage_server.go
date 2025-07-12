package app

import (
	"storage-service/internal/delivery/http_handlers"
	"storage-service/internal/server"
	"storage-service/internal/service"
	"time"
)

func Run(port string, timeout time.Duration) {
	storageService := service.NewStorageService("storage/memes/")
	storageHandler := http_handlers.NewImagesHandler(storageService)

	handler := server.NewRouter(storageHandler)

	server.StartServer(handler, port, timeout)
}