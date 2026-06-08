package route

import (
	"path"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/logger"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/middleware"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/utils"
	usecase "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases"
	"github.com/rs/zerolog"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, uc *usecase.Usecase, routes ...Route) {

	var rateLimterLogger = initLoggers()

	// middlewares

	r.Use(
		middleware.RateLimiterMiddleware(rateLimterLogger),
		middleware.ApiKeyMiddleware(),
		middleware.CORSMiddleware(),
	)
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	v1ApiGroup := r.Group("/api/v1")

	for _, route := range routes {
		route.Register(v1ApiGroup)
	}
}

func initLoggers() (ratelimiterLogger *zerolog.Logger) {
	rootDir := utils.GetWorkingDir()

	rateLimiterPath := path.Join(rootDir, "logs/identity/rate_limiter.log")

	rateLimterLogger := logger.NewLogger(logger.LoggerConfig{
		Level:      "warning",
		Filename:   rateLimiterPath,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		IsDev:      "local",
	})

	return rateLimterLogger
}
