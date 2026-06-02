package main

import (
	"log"

	"github.com/gin-gonic/gin"
	v1handler "github.com/lxbachit03/your-comfort-my-apple-golang/internal/api/v1/handler"
	v2handler "github.com/lxbachit03/your-comfort-my-apple-golang/internal/api/v2/handler"
)

func main() {
	r := gin.Default()

	v1Group := r.Group("api/v1")
	{
		userHandlerV1 := v1handler.NewUserHandler()
		productHandlerV1 := v1handler.NewProductHandler()

		userGroup := v1Group.Group("users")
		{
			userGroup.GET("/", userHandlerV1.GetUsersV1)
			userGroup.GET(("/:uuid"), userHandlerV1.GetUserByUUID)
		}

		productGroup := v1Group.Group("products")
		{
			productGroup.GET("/", productHandlerV1.GetProductsV1)
			productGroup.GET(("/:id"), productHandlerV1.GetProductById)
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
