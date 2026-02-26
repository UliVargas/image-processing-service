package utils

// AppError define la estructura de error de toda la aplicación
type AppError struct {
	StatusCode int         `json:"-"`       // Para el header HTTP
	Code       string      `json:"code"`    // Para lógica del frontend
	Message    string      `json:"message"` // Para el usuario
	Details    interface{} `json:"details,omitempty"`
}

// Implementación de la interfaz 'error' de Go
func (e *AppError) Error() string {
	return e.Message
}

// Fábrica de errores (Simple y eficiente)
func NewError(status int, code, msg string, details interface{}) *AppError {
	return &AppError{
		StatusCode: status,
		Code:       code,
		Message:    msg,
		Details:    details,
	}
}
