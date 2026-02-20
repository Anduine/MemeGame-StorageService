package server

import (
	"context"
	"fmt"
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
		slog.Info("Storage service running on port: " + port)
		err := server.ListenAndServe() //"certs/cert.pem", "certs/key.pem"
		if err != nil && err != http.ErrServerClosed {
			slog.Error("Server error: ", "Error", err)
		}
	}()

	// Ожидание сигнала завершения
	<-stop
	slog.Info("Shutting down server...")

	// Завершение работы сервера
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	slog.Info("Shutdown ", "stopcode", server.Shutdown(ctx))
}
