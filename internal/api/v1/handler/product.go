package v1handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (ProductHandler *ProductHandler) GetProductsV1(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "get products v1",
	})
}
