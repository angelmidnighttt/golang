package utils

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// HandleValidationErrors chuyển lỗi bind/validate của gin thành JSON dễ đọc.
//
//	Lỗi validator  -> {"errors": {"ID": "ID must be greater than 0"}}
//	Lỗi khác       -> {"error": "<message gốc>"}
func HandleValidationErrors(err error) gin.H {
	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		return gin.H{"error": err.Error()}
	}

	errs := make(map[string]string, len(validationErrs))
	for _, e := range validationErrs {
		errs[e.Field()] = messageForTag(e)
	}
	return gin.H{"errors": errs}
}

func messageForTag(e validator.FieldError) string {
	field := e.Field()
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, e.Param())
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, e.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
