package file

import (
	"errors"
	"image"
	"image-processing-service/internal/shared/auth"
	"image-processing-service/internal/shared/utils"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	_ "golang.org/x/image/webp"

	"github.com/go-chi/chi/v5"
)

type Handler interface {
	Upload(w http.ResponseWriter, r *http.Request)
	GetOne(w http.ResponseWriter, r *http.Request)
	ListMine(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service Service
}

var (
	ErrFileRequired    = utils.NewError(400, "FILE_REQUIRED", "No se ha proporcionado ningún archivo en la petición", nil)
	ErrFileTooLarge    = utils.NewError(413, "FILE_TOO_LARGE", "El archivo excede el tamaño máximo permitido (10MB)", nil)
	ErrInvalidFileType = utils.NewError(415, "UNSUPPORTED_FILE_TYPE", "El tipo de archivo no está permitido (solo JPG, PNG, GIF o WEBP)", nil)
	ErrFileRead        = utils.NewError(500, "FILE_READ_ERROR", "Error al procesar el archivo en el servidor", nil)
	ErrStorageUpload   = utils.NewError(502, "STORAGE_UPLOAD_FAILED", "No se pudo subir el archivo al almacenamiento remoto", nil)
	ErrUnauthorized    = utils.NewError(401, "UNAUTHORIZED", "Debes iniciar sesión para subir imágenes", nil)
)

type uploadFileResponse struct {
	ID           string `json:"id"`
	OriginalName string `json:"originalName"`
	FileName     string `json:"filename"`
	MimeType     string `json:"mimeType"`
	Size         int64  `json:"size"`
	Width        int64  `json:"width"`
	Height       int64  `json:"height"`
	Format       string `json:"format"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
	UserID       string `json:"userId,omitempty"`
	CreatedAt    string `json:"createdAt"`
}

func NewHandler(s Service) Handler {
	return &handler{service: s}
}

func (h *handler) Upload(w http.ResponseWriter, r *http.Request) {
	const maxFileSize = 10 << 20
	const maxRequestSize = maxFileSize + (1 << 20)
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)

	if err := r.ParseMultipartForm(maxRequestSize); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) || strings.Contains(err.Error(), "request body too large") {
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

	config, format, err := image.DecodeConfig(file)
	if err != nil {
		utils.HandleError(w, ErrInvalidFileType)
		return
	}

	if _, seekErr := file.Seek(0, io.SeekStart); seekErr != nil {
		utils.HandleError(w, ErrFileRead)
		return
	}

	if fileHeader.Size > maxFileSize {
		utils.HandleError(w, ErrFileTooLarge)
		return
	}

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
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	if !allowedTypes[mimeType] {
		utils.HandleError(w, ErrInvalidFileType)
		return
	}

	req := FileUploadRequest{
		FileName: fileHeader.Filename,
		MimeType: mimeType,
		FileSize: fileHeader.Size,
		Format:   strings.ToUpper(format),
		Width:    int64(config.Width),
		Height:   int64(config.Height),
	}

	authUser, ok := auth.GetAuthUser(r.Context())
	if !ok || authUser.UserID == "" {
		utils.HandleError(w, ErrUnauthorized)
		return
	}
	req.UserID = authUser.UserID

	if errs := utils.Validate(req); errs != nil {
		utils.HandleError(w, utils.ValidationError(errs))
		return
	}

	uploadedFile, err := h.service.Upload(file, req)
	if err != nil {
		utils.HandleError(w, ErrStorageUpload)
		return
	}

	utils.Success(w, http.StatusCreated, mapUploadFileResponse(uploadedFile, r))
}

func (h *handler) GetOne(w http.ResponseWriter, r *http.Request) {
	authUser, ok := auth.GetAuthUser(r.Context())
	if !ok || authUser.UserID == "" {
		utils.HandleError(w, ErrUnauthorized)
		return
	}

	storageKey := chi.URLParam(r, "*")
	if storageKey == "" {
		storageKey = chi.URLParam(r, "file")
	}

	decodedStorageKey, err := url.PathUnescape(storageKey)
	if err != nil || strings.TrimSpace(decodedStorageKey) == "" {
		utils.HandleError(w, ErrNotFound)
		return
	}

	fileMetadata, file, err := h.service.GetFile(decodedStorageKey, authUser.UserID)
	if err != nil {
		utils.HandleError(w, err)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", fileMetadata.MimeType)

	io.Copy(w, file)
}

func (h *handler) ListMine(w http.ResponseWriter, r *http.Request) {
	authUser, ok := auth.GetAuthUser(r.Context())
	if !ok || authUser.UserID == "" {
		utils.HandleError(w, ErrUnauthorized)
		return
	}

	if !utils.IsValidID(authUser.UserID) {
		utils.HandleError(w, ErrUnauthorized)
		return
	}

	files, err := h.service.ListByUserID(authUser.UserID)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	response := make([]uploadFileResponse, 0, len(files))
	for i := range files {
		response = append(response, mapUploadFileResponse(&files[i], r))
	}

	utils.Success(w, http.StatusOK, response)
}

func mapUploadFileResponse(file *File, r *http.Request) uploadFileResponse {
	fileURL := buildFileURL(r, file.StorageKey)
	thumbnailURL := ""
	if file.ThumbnailStorageKey != "" {
		thumbnailURL = buildFileURL(r, file.ThumbnailStorageKey)
	}

	return uploadFileResponse{
		ID:           file.ID,
		OriginalName: file.FileName,
		FileName:     path.Base(file.StorageKey),
		MimeType:     file.MimeType,
		Size:         file.FileSize,
		Width:        file.Width,
		Height:       file.Height,
		Format:       file.Format,
		URL:          fileURL,
		ThumbnailURL: thumbnailURL,
		UserID:       file.UserID,
		CreatedAt:    file.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func buildFileURL(r *http.Request, storageKey string) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if forwardedProto := r.Header.Get("X-Forwarded-Proto"); forwardedProto != "" {
		scheme = forwardedProto
	}

	return scheme + "://" + r.Host + "/api/v1/files/" + storageKey
}
