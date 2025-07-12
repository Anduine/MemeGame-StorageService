package http_handlers

import (
	"log"
	"net/http"
	"storage-service/internal/service"

	"github.com/gorilla/mux"
)

type ImagesHandler struct {
	service *service.StorageService
}

func NewImagesHandler(service *service.StorageService) *ImagesHandler {
	return &ImagesHandler{
		service: service,
	}
}

func (h *ImagesHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	fileName := mux.Vars(r)["filename"]
	filePath := h.service.GetImage(fileName)

	if filePath == "" {
		http.Error(w, "Файл не найден", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}

func (h *ImagesHandler) SaveImage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(20 << 20) // 20MB
	if err != nil {
		http.Error(w, "Ошибка разбора формы", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Файл не найден в форме", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := r.FormValue("filename")
	if fileName == "" {
		// Если имя не передали, использовать оригинальное
		fileName = handler.Filename
	}

	savedName, err := h.service.SaveImage(file, fileName)
	if err != nil {
		log.Printf("Ошибка при сохранении файла: %v", err)
		http.Error(w, "Ошибка при сохранении файла", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(savedName))
}