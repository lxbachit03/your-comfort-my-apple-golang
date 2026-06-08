package query

import (
	"context"
	"errors"

	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/decorator"
	v1dto "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/endpoints/dtos/v1"
)

type GetUserByUUIDQuery struct {
	UUID string
}

type getUserByUUIDHandler struct {
}

type GetUserByUUIDHandler decorator.QueryHandler[GetUserByUUIDQuery, *v1dto.UserDTO]

func NewGetUserByUUIDHandler() GetUserByUUIDHandler {
	return decorator.ApplyQueryDecorators[GetUserByUUIDQuery, *v1dto.UserDTO](
		getUserByUUIDHandler{},
	)
}

func (h getUserByUUIDHandler) Handle(ctx context.Context, q GetUserByUUIDQuery) (*v1dto.UserDTO, error) {
	if q.UUID == "" {
		return nil, errors.New("uuid is required")
	}

	// Returns fake data for now based on the requested UUID
	return &v1dto.UserDTO{
		UUID:  q.UUID,
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}, nil
}
