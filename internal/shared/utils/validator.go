package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func Validate(s interface{}) map[string]string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	formattedErrors := make(map[string]string)
	if valErrors, ok := err.(validator.ValidationErrors); ok {
		for _, f := range valErrors {
			formattedErrors[f.Field()] = msgForTag(f)
		}
	}

	return formattedErrors
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Este campo es obligatorio"
	case "email":
		return "El formato del correo electrónico no es válido"
	case "min":
		return fmt.Sprintf("Debe tener al menos %s caracteres", fe.Param())
	case "max":
		return fmt.Sprintf("No puede tener más de %s caracteres", fe.Param())
	}
	return "Valor inválido"
}
