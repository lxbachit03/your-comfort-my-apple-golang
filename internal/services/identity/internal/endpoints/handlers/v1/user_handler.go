package v1handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	usecase "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases"
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

	ctx.JSON(http.StatusOK, nil)
}
