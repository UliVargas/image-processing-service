package auth

import (
	"encoding/json"
	"fmt"
	"image-processing-service/internal/shared/auth"
	"image-processing-service/internal/shared/utils"
	"log"
	"net/http"
)

var (
	ErrInvalidJSON     = utils.NewError(400, "INVALID_JSON", "El cuerpo de la petición no es un JSON válido", nil)
	ErrInvalidIDFormat = utils.NewError(400, "INVALID_ID", "El formato del identificador proporcionado es incorrecto", nil)
	ErrValidation      = utils.NewError(422, "VALIDATION_FAILED", "Error de validación", nil)
)

type Handler interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	SignOut(w http.ResponseWriter, r *http.Request)
	RenewSession(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service Service
}

func NewHandler(s Service) Handler {
	return &handler{service: s}
}

func (h *handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

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

	user, err := h.service.SignUp(req)
	log.Println(err)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.Success(w, http.StatusCreated, user)
}

func (h *handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

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

	result, err := h.service.SignIn(req)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.Success(w, http.StatusOK, result)
}

func (h *handler) SignOut(w http.ResponseWriter, r *http.Request) {
	authUser, _ := auth.GetAuthUser(r.Context())
	fmt.Println("authUser", authUser)

	if err := h.service.SignOut(authUser.JTI); err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.Success(w, 200, map[string]string{"message": "Sesión cerrada correctamente"})
}

func (h *handler) RenewSession(w http.ResponseWriter, r *http.Request) {
	var req RenewSessionRequest

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

	result, err := h.service.RenewSession(req.RefreshToken)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.Success(w, http.StatusOK, result)
}
