package server

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServer(router http.Handler, port string, timeout time.Duration) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	// Канал для остановки сервера
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	go func() {
		log.Println("Storage service running on port:", port)
		err := server.ListenAndServe() //"certs/cert.pem", "certs/key.pem"
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error: ", slog.Any("Error", err))
		}
	}()

	// Ожидание сигнала завершения
	<-stop
	log.Println("Shutting down server...")

	// Завершение работы сервера
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	log.Println("Shutdown ", slog.Any("stopcode", server.Shutdown(ctx)))
}
