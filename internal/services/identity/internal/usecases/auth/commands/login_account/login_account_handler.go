package command

import (
	"context"

	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/decorator"
)

type LoginAccountCommand struct {
	UserName string
	Password string
}

type LoginAccountResult struct {
	AccessToken  string
	RefreshToken string
}

type LoginAccountHandler decorator.CommandHandler[LoginAccountCommand, LoginAccountResult]

type loginAccountHandler struct {
}

func NewLoginAccountHandler() LoginAccountHandler {
	return decorator.ApplyCommandDecorators[LoginAccountCommand, LoginAccountResult](
		loginAccountHandler{},
	)
}

func (h loginAccountHandler) Handle(ctx context.Context, cmd LoginAccountCommand) (LoginAccountResult, error) {

	return LoginAccountResult{}, nil
}
