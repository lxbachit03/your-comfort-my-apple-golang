package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	v1handler "github.com/lxbachit03/your-comfort-my-apple-golang/internal/api/v1/handler"
	v2handler "github.com/lxbachit03/your-comfort-my-apple-golang/internal/api/v2/handler"
	"github.com/lxbachit03/your-comfort-my-apple-golang/middlewares"
	"github.com/lxbachit03/your-comfort-my-apple-golang/utils"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	if err := utils.RegisterCustomValidators(); err != nil {
		panic(err)
	}

	r := gin.Default()

	r.Use(
		middlewares.RateLimiterMiddleware(),
		middlewares.LoggerMiddleware(),
	)

	go middlewares.CleanUpRateLimiter()

	v1Group := r.Group("api/v1")
	{
		userHandlerV1 := v1handler.NewUserHandler()
		productHandlerV1 := v1handler.NewProductHandler()
		categoryHandlerV1 := v1handler.NewCategoryHandler()

		userGroup := v1Group.Group("users").Use(middlewares.ApiKeyMiddleware())
		{
			userGroup.GET("/", userHandlerV1.GetUsersV1)
			userGroup.GET("/:uuid", userHandlerV1.GetUserByUUID)
		}

		productGroup := v1Group.Group("products")
		{
			productGroup.GET("/", productHandlerV1.GetProductsV1)
			productGroup.GET("/test/:id", productHandlerV1.GetProductById)
			productGroup.GET("/details/:slug", productHandlerV1.GetProductBySlug)
			productGroup.POST("/", productHandlerV1.PostProduct)
		}

		categoryGroup := v1Group.Group("categories")
		{
			categoryGroup.GET("/:path", categoryHandlerV1.GetCategoryByPath)
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
