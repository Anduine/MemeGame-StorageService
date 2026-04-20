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
		slog.Debug("Файл не найден:", "filepath", filePath)
		responseHTTP.JSONError(w, "Файл не найден", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}

func (h *ImagesHandler) SaveImages(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(30 << 20)
	if err != nil {
		slog.Debug("Ошибка парсинга формы", "error", err.Error())
		responseHTTP.JSONError(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		slog.Debug("Файлы отсутсвуют в форме")
		responseHTTP.JSONError(w, "Файлы отсутсвуют в форме", http.StatusBadRequest)
		return
	}

	var savedFilenames []string

	for _, fileHeader := range files {

		file, err := fileHeader.Open()
		if err != nil {
			slog.Debug("Ошибка чтения файла:", "filename", fileHeader.Filename, "error", err)
			continue
		}

		filename := fileHeader.Filename

		savedName, err := h.service.SaveImage(file, "images/", filename)
		file.Close()

		if err != nil {
			slog.Debug("Ошибка при сохранении файла:", "filename", filename, "error", err.Error())
			responseHTTP.JSONError(w, "Ошибка при сохранении файла", http.StatusInternalServerError)
			return
		}
		savedFilenames = append(savedFilenames, savedName)
	}

	slog.Debug("Файлы успешно сохранены:", "count", len(savedFilenames))

	responseHTTP.JSONResp(w, savedFilenames, http.StatusCreated)
}

func (h *ImagesHandler) DeleteImages(w http.ResponseWriter, r *http.Request) {
	var req domain.DeleteImagesRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Debug("Неправильный формат запроса:", "error", err.Error())
		responseHTTP.JSONError(w, "Неправильный формат запроса", http.StatusBadRequest)
		return
	}

	deletedCount := 0
	for _, filename := range req.Filenames {
		err := h.service.DeleteImage("images/", filename)
		if err != nil {
			slog.Debug("Не удалося удалить файл:", "filename", filename, "error", err.Error())
		} else {
			deletedCount++
		}
	}

	slog.Info("Удалено файлов: " + strconv.Itoa(deletedCount))
	responseHTTP.JSONError(w, "Удалено файлов: "+strconv.Itoa(deletedCount), http.StatusOK)
}

func (h *ImagesHandler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["filename"]
	filePath := h.service.GetImage("avatars/", filename)

	if filePath == "" {
		slog.Debug("Файл не найден:", "filepath", filePath)
		http.Error(w, "Файл не найден", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}

func (h *ImagesHandler) SaveAvatar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(30 << 20)
	if err != nil {
		slog.Debug("Ошибка парсинга формы:", "error", err.Error())
		http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		slog.Debug("Файл отсутсвует в форме:", "error", err.Error())
		http.Error(w, "Файл отсутсвует в форме", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := r.FormValue("filename")
	if fileName == "" {
		fileName = handler.Filename
	}

	savedName, err := h.service.SaveImage(file, "avatars/", fileName)
	if err != nil {
		slog.Info("Ошибка при сохранении файла:", "error", err.Error())
		http.Error(w, "Ошибка при сохранении файла", http.StatusInternalServerError)
		return
	}

	slog.Debug("Файл успешно сохранён:", "filename", savedName)

	responseHTTP.JSONResp(w, savedName, http.StatusCreated)
}

func (h *ImagesHandler) DeleteAvatar(w http.ResponseWriter, r *http.Request) {
	var req domain.DeleteAvatarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Debug("Неправильный формат запроса:", "error", err.Error())
		http.Error(w, "Неправильный формат запроса", http.StatusBadRequest)
		return
	}

	slog.Debug("Файл успешно удалён:", "filename", req.Filename)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(req.Filename))
}
