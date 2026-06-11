package v1routes

import (
	"github.com/gin-gonic/gin"
	v1handler "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/endpoints/handlers/v1"
)

type userRoutes struct {
	handler *v1handler.UserRouteHandler
}

func NewUserRoutes(handler *v1handler.UserRouteHandler) *userRoutes {
	return &userRoutes{
		handler: handler,
	}
}

func (ur *userRoutes) Register(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.GET("", ur.handler.GetUsers)
		users.GET("/:uuid", ur.handler.GetUserByUUID)
		users.POST("/address", ur.handler.AddAddress)
	}
}
