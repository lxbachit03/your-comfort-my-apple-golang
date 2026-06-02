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

			root := strings.Split(e.Namespace(), ".")[0]

			path := strings.TrimPrefix(e.Namespace(), root+".")

			parts := strings.Split(path, ".")

			for i, part := range parts {
				if strings.Contains(part, "[") {
					idx := strings.Index(part, "[")
					base := part[:idx]
					numberIdx := part[idx:]

					parts[i] = CamelCaseToSnakeCase(base) + numberIdx
				} else {
					parts[i] = CamelCaseToSnakeCase(part)
				}
			}

			fieldName := strings.Join(parts, ".")

			switch e.Tag() {
			case "uuid":
				errorsMap[fieldName] = "Invalid format (must be UUID)"
			case "required":
				errorsMap[fieldName] = "This field is required"
			case "gt":
				errorsMap[fieldName] = fmt.Sprintf("Must be greater than %v", e.Param())
			case "lte":
				errorsMap[fieldName] = fmt.Sprintf("Must be less than or equal to %v", e.Param())
			case "gte":
				errorsMap[fieldName] = fmt.Sprintf("Must be greater than or equal to %v", e.Param())
			case "datetime":
				errorsMap[fieldName] = fmt.Sprintf("Invalid date format (must be %v)", e.Param())
			case "slug":
				errorsMap[fieldName] = "Invalid slug format"
			case "oneof":
				allowedList := strings.Join(strings.Split(e.Param(), " "), ", ")
				errorsMap[fieldName] = fmt.Sprintf("Value must be one of: %v", allowedList)
			case "search":
				errorsMap[fieldName] = "Search query only accepts letters, numbers, and spaces"
			default:
				errorsMap[fieldName] = fmt.Sprintf("Invalid value %v on %s tag", e.Value(), e.Tag())
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
