package file

import (
	"image-processing-service/internal/shared/utils"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Handler interface {
	Upload(w http.ResponseWriter, r *http.Request)
	GetOne(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service Service
}

var (
	ErrValidation = utils.NewError(422, "VALIDATION_FAILED", "Error de validación", nil)
	// --- Nuevos errores para Archivos ---
	ErrFileRequired    = utils.NewError(400, "FILE_REQUIRED", "No se ha proporcionado ningún archivo en la petición", nil)
	ErrFileTooLarge    = utils.NewError(413, "FILE_TOO_LARGE", "El archivo excede el tamaño máximo permitido (10MB)", nil)
	ErrInvalidFileType = utils.NewError(415, "UNSUPPORTED_FILE_TYPE", "El tipo de archivo no está permitido (solo JPG, PNG o PDF)", nil)
	ErrFileRead        = utils.NewError(500, "FILE_READ_ERROR", "Error al procesar el archivo en el servidor", nil)
	ErrStorageUpload   = utils.NewError(502, "STORAGE_UPLOAD_FAILED", "No se pudo subir el archivo al almacenamiento remoto", nil)
)

func NewHandler(s Service) Handler {
	return &handler{service: s}
}

func (h *handler) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		if strings.Contains(err.Error(), "too large") {
			utils.HandleError(w, ErrFileTooLarge)
		} else {
			utils.HandleError(w, ErrFileRequired)
		}
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		utils.HandleError(w, ErrFileRequired)
		return
	}
	defer file.Close()

	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" {
		head := make([]byte, 512)
		n, readErr := file.Read(head)
		if readErr != nil && readErr != io.EOF {
			utils.HandleError(w, ErrFileRead)
			return
		}
		mimeType = http.DetectContentType(head[:n])
		if _, seekErr := file.Seek(0, io.SeekStart); seekErr != nil {
			utils.HandleError(w, ErrFileRead)
			return
		}
	}

	allowedTypes := map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"application/pdf": true,
	}

	if !allowedTypes[mimeType] {
		utils.HandleError(w, ErrInvalidFileType)
		return
	}

	req := FileUploadRequest{
		FileName:   fileHeader.Filename,
		StorageKey: utils.GenerateID(),
		MimeType:   mimeType,
		FileSize:   fileHeader.Size,
	}

	if errs := utils.Validate(req); errs != nil {
		utils.HandleError(w, utils.NewError(
			ErrValidation.StatusCode,
			ErrValidation.Code,
			ErrValidation.Message,
			errs,
		))
		return
	}

	if err := h.service.Upload(req); err != nil {
		utils.HandleError(w, ErrStorageUpload)
		return
	}

	utils.Success(w, http.StatusOK, "Archivo subido con éxito")
}

func (h *handler) GetOne(w http.ResponseWriter, r *http.Request) {
	storageKey := chi.URLParam(r, "file")

	file, err := h.service.FindOne(storageKey)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.Success(w, http.StatusOK, file)
}
