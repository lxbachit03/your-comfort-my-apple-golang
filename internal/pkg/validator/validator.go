package validator

import (
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidator() (*validator.Validate, error) {
	v, ok := binding.Validator.Engine().(*validator.Validate)

	if !ok {
		return nil, fmt.Errorf("failed to get validator engine")
	}

	return v, nil
}
