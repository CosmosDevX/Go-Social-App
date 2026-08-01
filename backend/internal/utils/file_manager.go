package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type FileManagerInterface interface {
	SaveFile(file multipart.File, header *multipart.FileHeader, saveDirectory string) (string, error)
	DeleteFile(path, filename string) error
}

type FileManager struct {
	RootDir string
}

func NewFileManager() FileManagerInterface {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	rootDir := filepath.Join(currentDir, "../..")

	return FileManager{
		RootDir: rootDir,
	}
}

func (m FileManager) SaveFile(file multipart.File, header *multipart.FileHeader, saveDirectory string) (string, error) {
	if file == nil || header == nil {
		return "", nil
	}

	fileExt := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), fileExt)
	savePath := filepath.Join(saveDirectory, filename)

	dst, err := os.Create(savePath)
	if err != nil {
		return "", fmt.Errorf("error during create empty file in %s", saveDirectory)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("error during write data to file")
	}

	return filename, nil
}

func (m FileManager) DeleteFile(path, filename string) error {
	if filename == "" {
		return nil
	}

	filePath := filepath.Join(m.RootDir+path, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file %s not exist in %s", filename, path)
	}

	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("error during file deleting %s", filename)
	}

	return nil
}
