package v1handler

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	utils "github.com/lxbachit03/your-comfort-my-apple-golang/utils"
)

type NewHandler struct{}

type CreateNewsRequestBody struct {
	Title  string                `form:"title" binding:"required"`
	Status *bool                 `form:"status" binding:"required"`
	Image  *multipart.FileHeader `form:"image" binding:"required"`
}

func NewNewsHandler() *NewHandler {
	return &NewHandler{}
}

func (newHandler *NewHandler) CreateNews(ctx *gin.Context) {

	var body CreateNewsRequestBody

	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"title":      body.Title,
		"status":     body.Status,
		"image_name": body.Image.Filename,
		"image_size": body.Image.Size,
	})
}
