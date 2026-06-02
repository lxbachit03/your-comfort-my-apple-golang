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

type ProductImage struct {
	ImageName string `json:"image_name" binding:"required"`
	ImageUrl  string `json:"image_url" binding:"required"`
}

type ProductAttribute struct {
	Name  string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type ProductInfo struct {
	InfoKey   string `json:"info_key" binding:"required"`
	InfoValue string `json:"info_value" binding:"required"`
}

type PostProductPostBody struct {
	Name              string                 `json:"name" binding:"required"`
	Price             float64                `json:"price" binding:"required,gte=0"`
	Display           *bool                  `json:"display" binding:"omitempty"`
	ProductImage      ProductImage           `json:"product_image" binding:"required"`
	Tags              []string               `json:"tags" binding:"required,gt=3"`
	ProductAttributes []ProductAttribute     `json:"product_attributes" binding:"required,gt=0,dive"`
	ProductInfo       map[string]ProductInfo `json:"product_info" binding:"required,gt=0,dive"`
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

func (ProductHandler *ProductHandler) PostProduct(ctx *gin.Context) {
	var body PostProductPostBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatValidationError(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"body": body,
	})
}
