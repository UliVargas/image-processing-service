package utils

// AppError define la estructura de error de toda la aplicación
type AppError struct {
	StatusCode int         `json:"-"`
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewError(status int, code, msg string, details interface{}) *AppError {
	return &AppError{
		StatusCode: status,
		Code:       code,
		Message:    msg,
		Details:    details,
	}
}

// ValidationError crea un error 422 con details de validación.
func ValidationError(details interface{}) *AppError {
	return NewError(422, "VALIDATION_FAILED", "Error de validación", details)
}

// Errores HTTP genéricos compartidos por todos los handlers.
var (
	ErrInvalidJSON     = NewError(400, "INVALID_JSON", "El cuerpo de la petición no es un JSON válido", nil)
	ErrInvalidIDFormat = NewError(400, "INVALID_ID", "El formato del identificador proporcionado es incorrecto", nil)
	ErrAlreadyExists   = NewError(409, "USER_ALREADY_EXISTS", "No es posible registrar este correo. Intenta con otro.", nil)
)
