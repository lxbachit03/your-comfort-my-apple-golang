package main

import (
	"log"

	"github.com/gin-gonic/gin"
	v1handler "github.com/lxbachit03/your-comfort-my-apple-golang/internal/api/v1/handler"
	v2handler "github.com/lxbachit03/your-comfort-my-apple-golang/internal/api/v2/handler"
)

func main() {
	r := gin.Default()

	userHandlerV1 := v1handler.UserHandler{}
	userHandlerV1 = *userHandlerV1.NewUserHandler()

	r.GET("api/v1/users", userHandlerV1.GetUsersV1)

	userHandlerV2 := v2handler.UserHandler{}
	userHandlerV2 = *userHandlerV2.NewUserHandler()

	r.GET("api/v2/users", userHandlerV2.GetUsersV2)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
