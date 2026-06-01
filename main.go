package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/demo", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "demo routes",
		})
	})

	r.GET("/languages/:lang", func(ctx *gin.Context) {
		pathParam := ctx.Param("lang")

		ctx.JSON(200, gin.H{
			"message": "languages route",
			"path":    pathParam,
		})
	})

	r.GET("/language/golang", func(ctx *gin.Context) {

		price := ctx.Query("price")
		level := ctx.Query("level")

		ctx.JSON(200, gin.H{
			"message": "golang route",
			"path":    "golang",
			"price":   price,
			"level":   level,
		})
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
