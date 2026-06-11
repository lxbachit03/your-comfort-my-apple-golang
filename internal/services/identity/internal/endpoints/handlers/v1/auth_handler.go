package v1handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apiresponse "github.com/lxbachit03/ygz-microservices-golang/internal/pkg/api-response"
	v1dto "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/endpoints/dtos/v1"
	identity_validator "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/utils/validation"
	usecase "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases"
	command "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases/auth/commands"
)

type AuthRouteHandler struct {
	uc *usecase.Usecase
}

func NewAuthRouteHandler(uc *usecase.Usecase) *AuthRouteHandler {
	return &AuthRouteHandler{
		uc: uc,
	}
}

func (arh *AuthRouteHandler) LoginAccount(ctx *gin.Context) {
	var request v1dto.LoginAccountRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		apiresponse.ValidationErrorResp(ctx, identity_validator.HandleValidationError(err))
		return
	}

	cmd := command.LoginAccountCommand{
		Email:    request.Email,
		Password: request.Password,
	}

	result, err := arh.uc.Commands.LoginAccountHandler.Handle(ctx.Request.Context(), cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
