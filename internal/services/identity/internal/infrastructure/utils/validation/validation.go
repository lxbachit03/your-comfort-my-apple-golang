package identity_validator

import (
	"fmt"
	"log"
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
				log.Printf("part: %v", part)
				if strings.Contains(part, "[") {
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
				errors[fieldPath] = "Is required"
			case "gt":
				errors[fieldPath] = fmt.Sprintf("Must be greater than %s", e.Param())
			case "lt":
				errors[fieldPath] = fmt.Sprintf("Must be less than %s", e.Param())
			case "gte":
				errors[fieldPath] = fmt.Sprintf("Must be greater than or equal to %s", e.Param())
			case "lte":
				errors[fieldPath] = fmt.Sprintf("Must be less than or equal to %s", e.Param())
			case "oneof":
				allowedValues := strings.Join(strings.Split(e.Param(), " "), ",")
				errors[fieldPath] = fmt.Sprintf("Must be one of: %s", allowedValues)
			case "min":
				errors[fieldPath] = fmt.Sprintf("Must be at least %s characters long", e.Param())
			case "max":
				errors[fieldPath] = fmt.Sprintf("Must be at most %s characters long", e.Param())
			case "min_int":
				errors[fieldPath] = fmt.Sprintf("Must be at least %s", e.Param())
			case "max_int":
				errors[fieldPath] = fmt.Sprintf("Must be at most %s", e.Param())
			case "email":
				errors[fieldPath] = "Must be a valid email address"
			case "uuid":
				errors[fieldPath] = "Must be a valid UUID"
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
