package command

import (
	"context"
	"log"

	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/decorator"
)

type UpdateProfileCommand struct {
	ProfileId       string
	FirstName       string
	LastName        string
	PhoneNumber     string
	BirthDay        string
	Gender          string
	ProfileImageUrl string
}

type UpdateProfileHandler decorator.CommandHandler[UpdateProfileCommand, bool]

type updateProfileHandler struct {
}

func NewUpdateProfileHandler() UpdateProfileHandler {
	return decorator.ApplyCommandDecorators[UpdateProfileCommand, bool](
		updateProfileHandler{},
	)
}

func (h updateProfileHandler) Handle(ctx context.Context, cmd UpdateProfileCommand) (bool, error) {

	log.Printf("cmd: %+v", cmd)

	return true, nil
}
