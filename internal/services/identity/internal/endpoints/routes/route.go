package route

import (
	"path"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/logger"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/middleware"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/config"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/utils"
	usecase "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases"
	"github.com/rs/zerolog"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, uc *usecase.Usecase, routes ...Route) {

	rateLimterLogger, recoveryLogger := initLoggers()

	// middlewares

	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(
		middleware.CORSMiddleware(),
		middleware.RateLimiterMiddleware(rateLimterLogger),
		// middleware.HttpLoggerMiddleware
		middleware.RecoveryMiddleware(recoveryLogger),
		// middleware.ApiKeyMiddleware(),
		// AuthMiddleware
	)

	v1ApiGroup := r.Group("/api/v1")

	for _, route := range routes {
		route.Register(v1ApiGroup)
	}
}

func initLoggers() (*zerolog.Logger, *zerolog.Logger) {
	rootDir := utils.GetWorkingDir()

	rateLimiterPath := path.Join(rootDir, "logs/identity/rate_limiter.log")
	recoveryLoggerPath := path.Join(rootDir, "logs/identity/recovery.log")

	rateLimterLogger := logger.NewLogger(logger.LoggerConfig{
		Level:      "warning",
		Filename:   rateLimiterPath,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		AppEnv:     config.AppConfig.Server.AppEnv,
	})

	recoveryLogger := logger.NewLogger(logger.LoggerConfig{
		Level:      "error",
		Filename:   recoveryLoggerPath,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		AppEnv:     config.AppConfig.Server.AppEnv,
	})

	return rateLimterLogger, recoveryLogger
}
