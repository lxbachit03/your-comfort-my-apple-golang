package v1handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	utils "github.com/lxbachit03/your-comfort-my-apple-golang/utils"
)

type CategoryHandler struct{}

// validators
type GetCategoryByPathParam struct {
	Path string `uri:"path" binding:"oneof=men women electronics"`
}

type CreateCategoryRequestBody struct {
	Name string `form:"name" binding:"required"`
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{}
}

func (categoryHandler *CategoryHandler) CreateCategory(ctx *gin.Context) {

	var body CreateCategoryRequestBody

	if err := ctx.ShouldBind(&body); err != nil {

		log.Print("==" + err.Error() + "==")

		ctx.JSON(http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}
}

func (categoryHandler *CategoryHandler) GetCategoryByPath(ctx *gin.Context) {

	var categoryPath GetCategoryByPathParam

	if err := ctx.ShouldBindUri(&categoryPath); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"path": categoryPath.Path,
	})
}
