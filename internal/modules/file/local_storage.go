package file

import (
	"io"
	"os"
	"path/filepath"
)

type localStorage struct {
	uploadDir string
}

func NewLocalStorage(dir string) StorageProvider {
	return &localStorage{uploadDir: dir}
}

func (l *localStorage) Save(content io.Reader, objectKey string, _ string) (string, error) {
	fullPath := filepath.Join(l.uploadDir, filepath.FromSlash(objectKey))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", err
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}

	defer dst.Close()

	if _, err := io.Copy(dst, content); err != nil {
		return "", err
	}

	return objectKey, nil
}

func (l *localStorage) Get(storageKey string) (io.ReadCloser, error) {
	fullPath := filepath.Join(l.uploadDir, filepath.FromSlash(storageKey))

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}
