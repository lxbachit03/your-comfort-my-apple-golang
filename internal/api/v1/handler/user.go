package v1handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct{}

func (userHandler *UserHandler) NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (userHandler *UserHandler) GetUsersV1(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "get users v1",
	})
}
