package server

import (
	"log/slog"
	"net/http"
	"storage-service/internal/delivery/http_handlers"

	"github.com/gorilla/mux"
)

func NewRouter(imageHandler *http_handlers.ImagesHandler) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/api/storage/image/{filename}", imageHandler.GetImage).Methods("GET")
	router.HandleFunc("/api/storage/upload_images", imageHandler.SaveImages).Methods("POST")
	router.HandleFunc("/api/storage/delete_images", imageHandler.DeleteImages).Methods("POST")

	router.HandleFunc("/api/storage/avatar/{filename}", imageHandler.GetAvatar).Methods("GET")
	router.HandleFunc("/api/storage/upload_avatar", imageHandler.SaveAvatar).Methods("POST")
	router.HandleFunc("/api/storage/delete_avatar", imageHandler.DeleteAvatar).Methods("POST")

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Маршрут не найден", "method", r.Method, "path", r.URL.Path)
		http.Error(w, "Маршрут не найден", http.StatusNotFound)
	})

	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Метод запрещён", "method", r.Method, "path", r.URL.Path)
		http.Error(w, "Запрещённый метод", http.StatusMethodNotAllowed)
	})

	return router
}
