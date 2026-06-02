package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) gin.H {
	if errors, ok := err.(validator.ValidationErrors); ok {
		errorsMap := make(map[string]string)

		for _, e := range errors {
			switch e.Tag() {
			case "uuid":
				errorsMap[e.Field()] = "Invalid format (must be UUID)"
			case "gt":
				errorsMap[e.Field()] = fmt.Sprintf("Must be greater than %v", e.Param())
			case "slug":
				errorsMap[e.Field()] = "Invalid slug format"
			default:
				errorsMap[e.Field()] = fmt.Sprintf("Invalid value %v on %s tag", e.Value(), e.Tag())
			}
		}

		return gin.H{
			"errors": errorsMap,
		}
	}

	return gin.H{
		"error": "Something went wrong: " + err.Error(),
	}
}
