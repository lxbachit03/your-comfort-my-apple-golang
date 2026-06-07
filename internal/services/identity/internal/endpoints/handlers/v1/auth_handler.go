package v1handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	usecase "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases"
	command "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases/auth/commands/login_account"
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
	// var request v1dto.LoginAccount
	// if err := ctx.ShouldBindJSON(&input); err != nil {
	// 	utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
	// 	return
	// }

	cmd := command.LoginAccountCommand{
		UserName: "test",
		Password: "test",
	}

	result, err := arh.uc.Commands.LoginAccountHandler.Handle(ctx.Request.Context(), cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
