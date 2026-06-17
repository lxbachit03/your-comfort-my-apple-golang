package command

import (
	"context"

	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/decorator"
	v1dto "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/endpoints/dtos/v1"
)

type LoginAccountCommand struct {
	Email    string
	Password string
}

type LoginAccountHandler decorator.CommandHandler[LoginAccountCommand, v1dto.LoginAccountResponse]

type loginAccountHandler struct {
}

func NewLoginAccountHandler() LoginAccountHandler {
	return decorator.ApplyCommandDecorators(
		loginAccountHandler{},
	)
}

func (h loginAccountHandler) Handle(ctx context.Context, cmd LoginAccountCommand) (v1dto.LoginAccountResponse, error) {

	return v1dto.LoginAccountResponse{}, nil
}
