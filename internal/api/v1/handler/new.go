package v1handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	utils "github.com/lxbachit03/your-comfort-my-apple-golang/utils"
)

type NewHandler struct{}

type CreateNewsRequestBody struct {
	Title  string `json:"title" binding:"required"`
	Status *bool  `json:"status" binding:"required"`
	Image  string `json:"image" binding:"required"`
}

func NewNewsHandler() *NewHandler {
	return &NewHandler{}
}

func (newHandler *NewHandler) CreateNews(ctx *gin.Context) {

	var body CreateNewsRequestBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}
}
