package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	apiresponse "github.com/lxbachit03/ygz-microservices-golang/internal/pkg/api-response"
	errorcode "github.com/lxbachit03/ygz-microservices-golang/internal/pkg/error_code"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/db/sqlc"
)

type UserRepository interface {
	GetUserByUUID(ctx context.Context, userUUID uuid.UUID) (sqlc.User, error)
}

type userRepository struct {
	db sqlc.Querier
}

func NewUserRepository(db sqlc.Querier) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur userRepository) GetUserByUUID(ctx context.Context, userUUID uuid.UUID) (sqlc.User, error) {
	user, err := ur.db.GetUser(ctx, userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, &apiresponse.ApiError{
				StatusCode:   http.StatusBadRequest,
				ErrorCode:    errorcode.ErrCodeBadRequest,
				Message:      "User not found",
				ErrorDetails: err,
			}
		}

		return sqlc.User{}, fmt.Errorf("GetUserByUUID : query error: %w", err)
	}

	return user, nil
}
