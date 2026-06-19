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
		ClientIP: ctx.ClientIP(),
	}

	result, err := arh.uc.Commands.LoginAccountHandler.Handle(ctx.Request.Context(), cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// name, value, maxAge, path, domain, secure, httpOnly
	ctx.SetCookie("refresh_token", result.RefreshToken, 7*24*60*60, "/", "", false, true)

	ctx.JSON(http.StatusOK, result)
}

func (arh *AuthRouteHandler) RegisterAccount(ctx *gin.Context) {
	var request v1dto.RegisterAccountRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		apiresponse.ValidationErrorResp(ctx, identity_validator.HandleValidationError(err))
		return
	}

	cmd := command.RegisterAccountCommand{
		Email:           request.Email,
		Password:        request.Password,
		ConfirmPassword: request.ConfirmPassword,
		FirstName:       request.FirstName,
		LastName:        request.LastName,
	}

	result, err := arh.uc.Commands.RegisterAccountHandler.Handle(ctx.Request.Context(), cmd)
	if err != nil {
		apiresponse.ErrorResponse(ctx, err)
	}

	apiresponse.Response(ctx, http.StatusOK, "User registered successfully", result)
}
