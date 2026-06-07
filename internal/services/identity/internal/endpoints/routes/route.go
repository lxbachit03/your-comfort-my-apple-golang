package route

import (
	"github.com/gin-gonic/gin"
	usecase "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, uc *usecase.Usecase, routes ...Route) {

	// middlewares

	v1ApiGroup := r.Group("/api/v1")

	for _, route := range routes {
		route.Register(v1ApiGroup)
	}
}
