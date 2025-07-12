package server

import (
	"log"
	"net/http"
	"storage-service/internal/delivery/http_handlers"

	"github.com/gorilla/mux"
)

func NewRouter(imageHandler *http_handlers.ImagesHandler) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/api/storage/get_meme/{filename}", imageHandler.GetImage).Methods("GET")
	router.HandleFunc("/api/storage/post_meme", imageHandler.SaveImage).Methods("POST")

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Маршрут не найден: %s %s", r.Method, r.URL.Path)
		http.Error(w, "Маршрут не найден", http.StatusNotFound)
	})

	return router
}

