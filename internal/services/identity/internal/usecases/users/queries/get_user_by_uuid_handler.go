package query

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/decorator"
	v1dto "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/endpoints/dtos/v1"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/db/repository"
)

type GetUserByUUIDQuery struct {
	UUID string
}

type getUserByUUIDHandler struct {
	ur repository.UserRepository
}

type GetUserByUUIDHandler decorator.QueryHandler[GetUserByUUIDQuery, v1dto.UserDTO]

func NewGetUserByUUIDHandler(ur repository.UserRepository) GetUserByUUIDHandler {
	return decorator.ApplyQueryDecorators[GetUserByUUIDQuery, v1dto.UserDTO](
		getUserByUUIDHandler{
			ur: ur,
		},
	)
}

func (h getUserByUUIDHandler) Handle(ctx context.Context, q GetUserByUUIDQuery) (v1dto.UserDTO, error) {
	if q.UUID == "" {
		return v1dto.UserDTO{}, errors.New("uuid is required")
	}

	uuid, err := uuid.Parse(q.UUID)
	if err != nil {
		return v1dto.UserDTO{}, errors.New("invalid uuid")
	}

	user, err := h.ur.GetUserByUUID(ctx, uuid)
	if err != nil {
		return v1dto.UserDTO{}, err
	}

	// Returns fake data for now based on the requested UUID
	return v1dto.UserDTO{
		UUID:  user.UserUuid.String(),
		Name:  user.UserFullname,
		Email: user.UserEmail,
	}, nil
}
