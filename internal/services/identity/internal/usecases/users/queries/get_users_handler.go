package query

import (
	"context"

	apiresponse "github.com/lxbachit03/ygz-microservices-golang/internal/pkg/api-response"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/decorator"
	v1dto "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/endpoints/dtos/v1"
)

type GetUsersQuery struct {
	Search string
	Page   int32
	Limit  int32
	Order  string
	Sort   string
}

type getUsersHandler struct {
}

type GetUsersHandler decorator.QueryHandler[GetUsersQuery, apiresponse.PaginationApiResponse[v1dto.UserDTO]]

func NewGetUsersHandler() GetUsersHandler {
	return decorator.ApplyQueryDecorators[GetUsersQuery, apiresponse.PaginationApiResponse[v1dto.UserDTO]](
		getUsersHandler{},
	)
}

func (h getUsersHandler) Handle(ctx context.Context, q GetUsersQuery) (apiresponse.PaginationApiResponse[v1dto.UserDTO], error) {
	// TODO: replace with real database query
	fakeUsers := []v1dto.UserDTO{
		{UUID: "uuid-001", Name: "Alice Johnson", Email: "alice@example.com"},
		{UUID: "uuid-002", Name: "Bob Smith", Email: "bob@example.com"},
		{UUID: "uuid-003", Name: "Charlie Brown", Email: "charlie@example.com"},
	}

	page := q.Page
	if page == 0 {
		page = 1
	}
	limit := q.Limit
	if limit == 0 {
		limit = 10
	}

	totalRecords := int32(len(fakeUsers))
	totalPages := (totalRecords + limit - 1) / limit

	usersDTO := v1dto.MapUsersToDTO(fakeUsers)

	// apiresponse.NewPagiantionResponse
	result := apiresponse.NewPagination[v1dto.UserDTO](usersDTO, page, limit, totalRecords, totalPages)
	return *result, nil
}
