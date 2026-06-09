package v1handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apiresponse "github.com/lxbachit03/ygz-microservices-golang/internal/pkg/api-response"
	v1dto "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/endpoints/dtos/v1"
	identity_validator "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/utils/validation"
	usecase "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases"
	query "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases/users/queries"
)

type UserRouteHandler struct {
	uc *usecase.Usecase
}

func NewUserRouteHandler(uc *usecase.Usecase) *UserRouteHandler {
	return &UserRouteHandler{
		uc: uc,
	}
}

func (arh *UserRouteHandler) GetUsers(ctx *gin.Context) {
	var request v1dto.GetUsersParams
	if err := ctx.ShouldBindQuery(&request); err != nil {
		apiresponse.ValidationErrorResp(ctx, identity_validator.HandleValidationError(err))
		return
	}

	query := query.GetUsersQuery{
		Search: request.Search,
		Page:   request.Page,
		Limit:  request.Limit,
		Order:  request.Order,
		Sort:   request.Sort,
	}

	result, err := arh.uc.Queries.GetUsersHandler.Handle(ctx.Request.Context(), query)
	if err != nil {
		apiresponse.ErrorResponse(ctx, err)
		return
	}

	apiresponse.PaginationResponse(ctx, http.StatusOK, "Test GetUsers", result)
}

func (arh *UserRouteHandler) GetUserByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
		return
	}

	query := query.GetUserByUUIDQuery{
		UUID: uuid,
	}

	result, err := arh.uc.Queries.GetUserByUUIDHandler.Handle(ctx.Request.Context(), query)
	if err != nil {
		apiresponse.ErrorResponse(ctx, err)
		return
	}

	apiresponse.Response(ctx, http.StatusOK, "User retrieved successfully", result)
}
