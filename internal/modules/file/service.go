package file

import (
	"errors"
	"image-processing-service/internal/shared/utils"

	"gorm.io/gorm"
)

type Service interface {
	Upload(req FileUploadRequest) error
	FindOne(req string) (*File, error)
}

type service struct {
	repo Repository
}

var (
	ErrNotFound = utils.NewError(404, "FILE_NOT_FOUND", "Archivo no encontrado", nil)
)

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Upload(req FileUploadRequest) error {
	file := &File{
		ID:         utils.GenerateID(),
		FileName:   req.FileName,
		StorageKey: req.StorageKey,
		MimeType:   req.MimeType,
		FileSize:   req.FileSize,
	}

	if err := s.repo.Create(file); err != nil {
		return err
	}

	return nil
}

func (s *service) FindOne(req string) (*File, error) {
	file, err := s.repo.FindOne(req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return file, nil
}
