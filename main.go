package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func GetUsersV1(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "get users v1",
	})
}

func GetUsersV2(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "get users v2",
	})
}

func main() {
	r := gin.Default()

	r.GET("api/v1/users", GetUsersV1)
	r.GET("api/v2/users", GetUsersV2)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
