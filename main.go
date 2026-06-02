package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	v1handler "github.com/lxbachit03/your-comfort-my-apple-golang/internal/api/v1/handler"
	v2handler "github.com/lxbachit03/your-comfort-my-apple-golang/internal/api/v2/handler"
)

func main() {
	r := gin.Default()

	v1Group := r.Group("api/v1")
	{
		userHandlerV1 := v1handler.NewUserHandler()

		userGroup := v1Group.Group("users")
		{
			userGroup.GET("/", userHandlerV1.GetUsersV1)
		}

		productGroup := v1Group.Group("products")
		{
			productGroup.GET("/", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{
					"message": "get products v1",
				})
			})
		}
	}

	v2Group := r.Group("/api/v2")
	{
		userHandlerV2 := v2handler.NewUserHandler()

		userGroup := v2Group.Group("/users")
		{
			userGroup.GET("/", userHandlerV2.GetUsersV2)
		}

	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
