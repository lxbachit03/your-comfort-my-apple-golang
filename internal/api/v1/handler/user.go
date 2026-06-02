package v1handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	utils "github.com/lxbachit03/your-comfort-my-apple-golang/utils"
)

type UserHandler struct{}

// validators
type GetUserByUUIDPathParam struct {
	UUID string `uri:"uuid" binding:"uuid"`
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (userHandler *UserHandler) GetUsersV1(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "get users v1",
	})
}

func (userHandler *UserHandler) GetUserByUUID(ctx *gin.Context) {

	var param GetUserByUUIDPathParam

	if err := ctx.ShouldBindUri(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"uuid": param.UUID,
	})
}
