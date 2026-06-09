package identity_validator

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	apiresponse "github.com/lxbachit03/ygz-microservices-golang/internal/pkg/api-response"
	errorcode "github.com/lxbachit03/ygz-microservices-golang/internal/pkg/error_code"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/utils"
)

func HandleValidationError(err error) apiresponse.ValidationErrorResponse {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)

		for _, e := range validationErrors {
			root := strings.Split(e.Namespace(), ".")[0]           // LoginAccountRequest.Email -> LoginAccountRequest
			rawPath := strings.TrimPrefix(e.Namespace(), root+".") // Email
			parts := strings.Split(rawPath, ".")                   // [Email]

			for i, part := range parts {
				if strings.Contains("part", "[") {
					// idx := strings.Index(part, "[")
					// base := utils.CamelToSnake(part[:idx])
					// index := part[idx:]
					// parts[i] = base + index
				} else {
					parts[i] = utils.CamelToSnake(part)
				}
			}

			fieldPath := strings.Join(parts, ".")

			switch e.Tag() {
			case "required":
				errors[fieldPath] = fmt.Sprintf("%s is required", fieldPath)
			case "gt":
				errors[fieldPath] = fmt.Sprintf("%s must be greater than %s", fieldPath, e.Param())
			case "lt":
				errors[fieldPath] = fmt.Sprintf("%s must be less than %s", fieldPath, e.Param())
			case "gte":
				errors[fieldPath] = fmt.Sprintf("%s must be greater than or equal to %s", fieldPath, e.Param())
			case "lte":
				errors[fieldPath] = fmt.Sprintf("%s must be less than or equal to %s", fieldPath, e.Param())
			case "oneof":
				allowedValues := strings.Join(strings.Split(e.Param(), " "), ",")
				errors[fieldPath] = fmt.Sprintf("%s must be one of: %s", fieldPath, allowedValues)
			case "min":
				errors[fieldPath] = fmt.Sprintf("%s must be at least %s characters long", fieldPath, e.Param())
			case "max":
				errors[fieldPath] = fmt.Sprintf("%s must be at most %s characters long", fieldPath, e.Param())
			case "min_int":
				errors[fieldPath] = fmt.Sprintf("%s must be at least %s", fieldPath, e.Param())
			case "max_int":
				errors[fieldPath] = fmt.Sprintf("%s must be at most %s", fieldPath, e.Param())
			case "email":
				errors[fieldPath] = fmt.Sprintf("%s must be a valid email address", fieldPath)
			case "uuid":
				errors[fieldPath] = fmt.Sprintf("%s must be a valid UUID", fieldPath)
			}
		}

		return apiresponse.ValidationErrorResponse{
			StatusCode:   http.StatusBadRequest,
			ErrorCode:    errorcode.ErrCodeValidation,
			ErrorDetails: errors,
			Message:      "Validation error",
		}
	}

	return apiresponse.ValidationErrorResponse{
		StatusCode:   http.StatusBadRequest,
		ErrorCode:    errorcode.ErrCodeValidation,
		ErrorDetails: err.Error(),
		Message:      "Unpredicted validation error",
	}
}
