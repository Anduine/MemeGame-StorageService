package http_handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"storage-service/internal/domain"
	"storage-service/internal/lib/responseHTTP"
	"storage-service/internal/service"
	"strconv"

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
	filePath := h.service.GetImage("images/", fileName)

	if filePath == "" {
		responseHTTP.JSONError(w, http.StatusNotFound, "Файл не знайдено")
		return
	}

	http.ServeFile(w, r, filePath)
}

func (h *ImagesHandler) SaveImages(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(40 * 1024 * 1024)
	if err != nil {
		responseHTTP.JSONError(w, http.StatusBadRequest, "Помилка у парсингу формі")
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		responseHTTP.JSONError(w, http.StatusBadRequest, "Файли не знайдено")
		return
	}

	var savedFilenames []string

	for _, fileHeader := range files {

		file, err := fileHeader.Open()
		if err != nil {
			slog.Info("Помилка читання файлу", "filename", fileHeader.Filename, "err", err)
			continue
		}

		fileName := fileHeader.Filename

		savedName, err := h.service.SaveImage(file, "images/", fileName)
		file.Close()

		if err != nil {
			slog.Info("Помилка при збереженні файлу", "err", err.Error())
			http.Error(w, "Помилка при збережені файлу", http.StatusInternalServerError)
			return
		}
		savedFilenames = append(savedFilenames, savedName)
	}

	slog.Debug("Успішно збережено файлів", "count", len(savedFilenames))

	responseHTTP.JSONResp(w, http.StatusCreated, savedFilenames)
}

func (h *ImagesHandler) DeleteImages(w http.ResponseWriter, r *http.Request) {
	var req domain.DeleteImagesRequest
	// Декодируем JSON массив имен файлов
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseHTTP.JSONError(w, http.StatusBadRequest, "Невірний формат запиту")
		return
	}

	deletedCount := 0
	for _, filename := range req.Filenames {
		err := h.service.DeleteImage("images/", filename)
		if err != nil {
			slog.Info("Не вдалося видалити файл: ", "Filename", filename, "Error", err)
		} else {
			deletedCount++
		}
	}

	slog.Info("Видалено файлів: " + strconv.Itoa(deletedCount))
	responseHTTP.JSONError(w, http.StatusOK, "Видалено файлів: "+strconv.Itoa(deletedCount))
}

func (h *ImagesHandler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	fileName := mux.Vars(r)["filename"]
	filePath := h.service.GetImage("avatars/", fileName)

	if filePath == "" {
		http.Error(w, "Файл не знайдено", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}

func (h *ImagesHandler) SaveAvatar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		http.Error(w, "Помилка у формі", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		slog.Debug("Файл не знайдено у формі:", "error", err.Error())
		http.Error(w, "Файл не знайдено у формі", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := r.FormValue("filename")
	if fileName == "" {
		fileName = handler.Filename
	}

	savedName, err := h.service.SaveImage(file, "avatars/", fileName)
	if err != nil {
		slog.Info("Помилка при збережені файлу: ", "error", err)
		http.Error(w, "Помилка при збережені файлу", http.StatusInternalServerError)
		return
	}

	slog.Debug("Файл успішно збережено", "filename", savedName)

	responseHTTP.JSONResp(w, http.StatusCreated, savedName)
}

func (h *ImagesHandler) DeleteAvatar(w http.ResponseWriter, r *http.Request) {
	var req domain.DeleteAvatarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Невірний формат запиту", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteImage("images/", req.Filename)
	if err != nil {
		slog.Info("Не вдалося видалити файл: ", "Filename", req.Filename, "Error", err)
	}

	slog.Info("Видалено файл", "filename", req.Filename)

	w.WriteHeader(http.StatusOK)
}
