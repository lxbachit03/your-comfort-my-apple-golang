package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/middleware"
	usecase "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, uc *usecase.Usecase, routes ...Route) {

	// middlewares
	r.Use(
		middleware.ApiKeyMiddleware(),
		middleware.CORSMiddleware(),
	)

	v1ApiGroup := r.Group("/api/v1")

	for _, route := range routes {
		route.Register(v1ApiGroup)
	}
}
