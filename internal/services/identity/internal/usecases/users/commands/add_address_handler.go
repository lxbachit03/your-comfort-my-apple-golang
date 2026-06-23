package command

import (
	"context"

	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/decorator"
)

type AddAddressCommand struct {
	Label              string `json:"label" binding:"required"`
	ContactName        string `json:"contact_name" binding:"required"`
	ContactPhoneNumber string `json:"contact_phone_number" binding:"required"`
	AddressLine        string `json:"address_line" binding:"required"`
	District           string `json:"district" binding:"required"`
	Province           string `json:"province" binding:"required"`
	Country            string `json:"country" binding:"required"`
}

type AddAddressHandler decorator.CommandHandler[AddAddressCommand, bool]

type addAddressHandler struct {
}

func NewAddAddressHandler() AddAddressHandler {
	return decorator.ApplyCommandDecorators(
		addAddressHandler{},
	)
}

func (h addAddressHandler) Handle(ctx context.Context, cmd AddAddressCommand) (bool, error) {

	return false, nil
}
