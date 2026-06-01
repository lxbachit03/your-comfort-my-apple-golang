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

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
