package v1routes

import (
	"github.com/gin-gonic/gin"
	v1handler "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/endpoints/handlers/v1"
)

type authRoutes struct {
	handler *v1handler.AuthRouteHandler
}

func NewAuthRoutes(handler *v1handler.AuthRouteHandler) *authRoutes {
	return &authRoutes{
		handler: handler,
	}
}

func (ar *authRoutes) Register(r *gin.RouterGroup) {
	// public routes
	auth := r.Group("/auth")
	{
		auth.POST("/login", ar.handler.LoginAccount)
		auth.POST("/register", ar.handler.RegisterAccount)
	}
}
