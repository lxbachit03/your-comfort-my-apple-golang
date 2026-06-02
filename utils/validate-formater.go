package utils

import (
	"fmt"
	"strings"

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
			case "lte":
				errorsMap[e.Field()] = fmt.Sprintf("Must be less than or equal to %v", e.Param())
			case "gte":
				errorsMap[e.Field()] = fmt.Sprintf("Must be greater than or equal to %v", e.Param())
			case "datetime":
				errorsMap[e.Field()] = fmt.Sprintf("Invalid date format (must be %v)", e.Param())
			case "slug":
				errorsMap[e.Field()] = "Invalid slug format"
			case "oneof":
				allowedList := strings.Join(strings.Split(e.Param(), " "), ", ")
				errorsMap[e.Field()] = fmt.Sprintf("Value must be one of: %v", allowedList)
			case "search":
				errorsMap[e.Field()] = "Search query only accepts letters, numbers, and spaces"
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
