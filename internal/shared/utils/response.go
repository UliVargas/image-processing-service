package utils

import (
	"encoding/json"
	"net/http"
	"time"
)

type successResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data"`
	Timestamp string      `json:"timestamp"`
}

type errorResponse struct {
	Success   bool        `json:"success"`
	Error     interface{} `json:"error"`
	Timestamp string      `json:"timestamp"`
}

func Success(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := successResponse{
		Success:   true,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}

func HandleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var appErr *AppError
	var ok bool

	if appErr, ok = err.(*AppError); !ok {
		appErr = NewError(http.StatusInternalServerError, "INTERNAL_ERROR", "Algo salió mal de nuestro lado", nil)
	}

	w.WriteHeader(appErr.StatusCode)

	response := errorResponse{
		Success:   false,
		Error:     appErr,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}
