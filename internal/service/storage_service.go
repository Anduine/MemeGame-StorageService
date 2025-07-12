package service

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

type StorageService struct {
	storagePath string
}

func NewStorageService(path string) *StorageService {
	return &StorageService{
		storagePath: path,
	}
}

func (s *StorageService) GetImage(filename string) string {
	fileFullPath := fmt.Sprint(s.storagePath, filename)	

	log.Printf("Запрашиваемый файл: %s", fileFullPath)

	if _, err := os.Stat(fileFullPath); os.IsNotExist(err) {
		return ""
	}

	ext := filepath.Ext(fileFullPath)
	if ext != ".webp" && ext != ".png" && ext != ".jpg" && ext != ".jpeg"{
		return ""
	}

	return fileFullPath
}

func (s *StorageService) SaveImage(file multipart.File, fileName string) (string, error) {
	// Убедимся, что директория существует
	err := os.MkdirAll(s.storagePath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("не удалось создать папку для загрузки: %v", err)
	}

	dstPath := filepath.Join(s.storagePath, fileName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("не удалось создать файл: %v", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("не удалось сохранить файл: %v", err)
	}

	return fileName, nil
}