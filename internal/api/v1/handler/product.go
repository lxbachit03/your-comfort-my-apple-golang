package v1handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	utils "github.com/lxbachit03/your-comfort-my-apple-golang/utils"
)

type ProductHandler struct{}

// validators
type GetProductByIdPathParam struct {
	Id int `uri:"id" binding:"gt=0"`
}

type GetProductBySlugPathParam struct {
	Slug string `uri:"slug" binding:"slug"`
}

type GetProductsV1QueryParam struct {
	Search string `form:"search" binding:"omitempty,search"`
	Limit  int    `form:"limit" binding:"omitempty,lte=50"`
	Offset int    `form:"offset" binding:"omitempty,gte=1"`
	Date   string `form:"date" binding:"omitempty,datetime=2006-01-02"`
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (ProductHandler *ProductHandler) GetProductsV1(ctx *gin.Context) {

	var queryParams GetProductsV1QueryParam

	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	if queryParams.Limit == 0 {
		queryParams.Limit = 10
	}

	if queryParams.Offset == 0 {
		queryParams.Offset = 1
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "get products v1",
		"params":  queryParams,
	})
}

func (ProductHandler *ProductHandler) GetProductById(ctx *gin.Context) {
	var param GetProductByIdPathParam

	if err := ctx.ShouldBindUri(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": param.Id,
	})
}

func (ProductHandler *ProductHandler) GetProductBySlug(ctx *gin.Context) {
	var param GetProductBySlugPathParam

	if err := ctx.ShouldBindUri(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"slug": param.Slug,
	})
}
