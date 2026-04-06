package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"image-processing-service/internal/shared/auth"
	"image-processing-service/internal/shared/utils"

	"github.com/go-chi/chi/v5"
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
		utils.HandleError(w, utils.ErrInvalidIDFormat)
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
	const defaultPage = 1
	const defaultLimit = 10
	const maxLimit = 100

	page := defaultPage
	limit := defaultLimit

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			utils.HandleError(w, utils.ValidationError(nil))
			return
		}
		page = p
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l < 1 || l > maxLimit {
			utils.HandleError(w, utils.ValidationError(nil))
			return
		}
		limit = l
	}

	users, total, err := h.service.GetAll(page, limit)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	if users == nil {
		users = make([]*User, 0)
	}

	utils.Success(w, http.StatusOK, utils.PaginatedResult[*User]{
		Data: users,
		Meta: utils.PaginatedMeta{Total: total, Page: page, Limit: limit},
	})
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if !utils.IsValidID(id) {
		utils.HandleError(w, utils.ErrInvalidIDFormat)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, utils.ErrInvalidJSON)
		return
	}

	if errs := utils.Validate(req); errs != nil {
		utils.HandleError(w, utils.ValidationError(errs))
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
		utils.HandleError(w, utils.ErrInvalidIDFormat)
		return
	}

	var req UpdatePasswordUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, utils.ErrInvalidJSON)
		return
	}

	if errs := utils.Validate(req); errs != nil {
		utils.HandleError(w, utils.ValidationError(errs))
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
		utils.HandleError(w, utils.ErrInvalidIDFormat)
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.Success(w, http.StatusOK, map[string]string{"message": "Usuario eliminado correctamente"})
}
