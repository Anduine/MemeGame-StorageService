package service

import (
	"fmt"
	"io"
	"log/slog"
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

func (s *StorageService) GetImage(packageName, fileName string) string {
	fileFullPath := fmt.Sprint(s.storagePath, packageName, fileName)

	if _, err := os.Stat(fileFullPath); os.IsNotExist(err) {
		slog.Debug("Не знайдено файл", "filename", fileFullPath)
		return ""
	}

	ext := filepath.Ext(fileFullPath)
	if ext != ".webp" && ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		return ""
	}

	slog.Debug("Знайдено файл", "filename", fileFullPath)

	return fileFullPath
}

func (s *StorageService) SaveImage(file multipart.File, packageName, fileName string) (string, error) {
	fullDir := filepath.Join(s.storagePath, packageName)

	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return "", fmt.Errorf("не вдалося створити папку: %v", err)
	}

	dstPath := filepath.Join(fullDir, fileName)

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("не вдалося створити файл: %v", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, file)
	if err != nil {
		return "", fmt.Errorf("не вдалося зберігти файл: %v", err)
	}

	return fileName, nil
}

func (s *StorageService) DeleteImage(packageName, fileName string) error {
	safeName := filepath.Base(fileName)
	fullPath := filepath.Join(s.storagePath, packageName, safeName)

	err := os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}
