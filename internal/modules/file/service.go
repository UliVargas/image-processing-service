package file

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image-processing-service/internal/shared/utils"
	"image/jpeg"
	"io"
	"math"
	"path/filepath"
	"strings"

	xdraw "golang.org/x/image/draw"
	"gorm.io/gorm"
)

type Service interface {
	Upload(content io.Reader, req FileUploadRequest) (*File, error)
	GetFile(storageKey string, userID string) (*File, io.ReadCloser, error)
	ListByUserID(userID string) ([]File, error)
}

type service struct {
	repo    Repository
	storage StorageProvider
}

var (
	ErrNotFound = utils.NewError(404, "FILE_NOT_FOUND", "Archivo no encontrado", nil)
)

func NewService(r Repository, s StorageProvider) Service {
	return &service{repo: r, storage: s}
}

func (s *service) Upload(content io.Reader, req FileUploadRequest) (*File, error) {
	contentBytes, err := io.ReadAll(content)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(req.UserID) == "" {
		return nil, ErrInvalidFileType
	}

	fileID := utils.GenerateID()
	originalKey := buildOriginalObjectKey(req.UserID, fileID, req.MimeType, req.FileName)
	thumbnailKey := fmt.Sprintf("%s/thumbnails/%s.jpg", req.UserID, fileID)

	storageKey, err := s.storage.Save(bytes.NewReader(contentBytes), originalKey, req.MimeType)
	if err != nil {
		return nil, err
	}

	thumbnailBytes, err := generateThumbnail(contentBytes, 200, 200)
	if err != nil {
		return nil, err
	}

	storedThumbnailKey, err := s.storage.Save(bytes.NewReader(thumbnailBytes), thumbnailKey, "image/jpeg")
	if err != nil {
		return nil, err
	}

	file := &File{
		ID:                  fileID,
		FileName:            req.FileName,
		StorageKey:          storageKey,
		ThumbnailStorageKey: storedThumbnailKey,
		MimeType:            req.MimeType,
		FileSize:            req.FileSize,
		UserID:              req.UserID,
		Format:              req.Format,
		Width:               req.Width,
		Height:              req.Height,
	}

	if err := s.repo.Create(file); err != nil {
		return nil, err
	}

	return file, nil
}

func buildOriginalObjectKey(userID string, fileID string, mimeType string, fileName string) string {
	if ext := extensionFromMimeType(mimeType); ext != "" {
		return fmt.Sprintf("%s/images/%s%s", userID, fileID, ext)
	}

	ext := strings.ToLower(filepath.Ext(fileName))
	if ext == "" {
		ext = ".bin"
	}

	return fmt.Sprintf("%s/images/%s%s", userID, fileID, ext)
}

func extensionFromMimeType(mimeType string) string {
	switch strings.ToLower(strings.TrimSpace(mimeType)) {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ""
	}
}

func generateThumbnail(content []byte, maxWidth int, maxHeight int) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()
	if originalWidth <= 0 || originalHeight <= 0 {
		return nil, fmt.Errorf("dimensiones de imagen inválidas")
	}

	scale := math.Min(float64(maxWidth)/float64(originalWidth), float64(maxHeight)/float64(originalHeight))
	if scale > 1 {
		scale = 1
	}

	thumbWidth := int(math.Round(float64(originalWidth) * scale))
	thumbHeight := int(math.Round(float64(originalHeight) * scale))
	if thumbWidth < 1 {
		thumbWidth = 1
	}
	if thumbHeight < 1 {
		thumbHeight = 1
	}

	thumb := image.NewRGBA(image.Rect(0, 0, thumbWidth, thumbHeight))
	xdraw.CatmullRom.Scale(thumb, thumb.Bounds(), img, bounds, xdraw.Over, nil)

	var output bytes.Buffer
	if err := jpeg.Encode(&output, thumb, &jpeg.Options{Quality: 80}); err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

func (s *service) GetFile(storageKey string, userID string) (*File, io.ReadCloser, error) {
	fileMetadata, err := s.repo.FindOneByAnyKeyAndUserID(storageKey, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrNotFound
		}
		return nil, nil, err
	}

	if fileMetadata == nil {
		return nil, nil, ErrNotFound
	}

	file, err := s.storage.Get(storageKey)
	if err != nil {
		return nil, nil, ErrStorageUpload
	}

	if storageKey == fileMetadata.ThumbnailStorageKey {
		fileMetadata.MimeType = "image/jpeg"
	}

	return fileMetadata, file, nil
}

func (s *service) ListByUserID(userID string) ([]File, error) {
	return s.repo.FindByUserID(userID)
}
