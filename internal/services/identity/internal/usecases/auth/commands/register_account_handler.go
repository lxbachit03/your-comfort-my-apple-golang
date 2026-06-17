package command

import (
	"context"
	"log"

	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/decorator"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/db/repository"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/db/sqlc"
)

type RegisterAccountCommand struct {
	Email           string
	Password        string
	ConfirmPassword string
	FirstName       string
	LastName        string
}

type RegisterAccountHandler decorator.CommandHandler[RegisterAccountCommand, bool]

type registerAccountHandler struct {
	userRepo repository.UserRepository
}

func NewRegisterAccountHandler(userRepo repository.UserRepository) RegisterAccountHandler {

	if userRepo == nil {
		panic("nil userRepo repository")
	}

	return decorator.ApplyCommandDecorators(
		registerAccountHandler{
			userRepo: userRepo,
		},
	)
}

func (h registerAccountHandler) Handle(ctx context.Context, cmd RegisterAccountCommand) (bool, error) {

	age := int32(1)
	result, err := h.userRepo.CreateUser(ctx, sqlc.CreateUserParams{
		UserEmail:    cmd.Email,
		UserPassword: cmd.Password,
		UserFullname: cmd.FirstName + " " + cmd.LastName,
		UserAge:      &age,
		UserStatus:   1,
		UserLevel:    1,
	})

	log.Printf("Result: %+v", result)

	if err != nil {
		return false, err
	}

	return true, nil
}
