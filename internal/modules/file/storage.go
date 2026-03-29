package file

import "io"

type StorageProvider interface {
	Save(content io.Reader, objectKey string, contentType string) (string, error)
	Get(storageKey string) (io.ReadCloser, error)
}
