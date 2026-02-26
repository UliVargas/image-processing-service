package user

import (
	"encoding/json"
	"net/http"

	"image-processing-service/internal/shared/auth"
	"image-processing-service/internal/shared/utils"

	"github.com/go-chi/chi/v5"
)

var (
	ErrInvalidJSON     = utils.NewError(400, "INVALID_JSON", "El cuerpo de la petición no es un JSON válido", nil)
	ErrInvalidIDFormat = utils.NewError(400, "INVALID_ID", "El formato del identificador proporcionado es incorrecto", nil)
	ErrValidation      = utils.NewError(422, "VALIDATION_FAILED", "Error de validación", nil)
)

type Handler interface {
	GetByID(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	UpdatePassword(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service Service
}

func NewHandler(s Service) Handler {
	return &handler{service: s}
}

func (h *handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if !utils.IsValidID(id) {
		utils.HandleError(w, ErrInvalidIDFormat)
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.Success(w, http.StatusOK, user)
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll()
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.Success(w, http.StatusOK, users)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if !utils.IsValidID(id) {
		utils.HandleError(w, ErrInvalidIDFormat)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, ErrInvalidJSON)
		return
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

	user, err := h.service.Update(id, req)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.Success(w, 201, user)
}

func (h *handler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	authUser, _ := auth.GetAuthUser(r.Context())

	if !utils.IsValidID(authUser.UserID) {
		utils.HandleError(w, ErrInvalidIDFormat)
		return
	}

	var req UpdatePasswordUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, ErrInvalidJSON)
		return
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

	user, err := h.service.UpdatePassword(authUser.UserID, req)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.Success(w, 201, user)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if !utils.IsValidID(id) {
		utils.HandleError(w, ErrInvalidIDFormat)
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.Success(w, http.StatusOK, map[string]string{"message": "Usuaurio eliminado correctamente"})
}
