package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	auth_pkg "github.com/lxbachit03/ygz-microservices-golang/internal/pkg/auth/jwt"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/cache"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/events"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/logger"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/queue/rabbitmq"
	hash_pkg "github.com/lxbachit03/ygz-microservices-golang/internal/pkg/security/hash"
	validator "github.com/lxbachit03/ygz-microservices-golang/internal/pkg/validator"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/config"
	v1handler "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/endpoints/handlers/v1"
	routes "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/endpoints/routes"
	v1routes "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/endpoints/routes/v1"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/db"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/db/repository"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/external/mail"
	ext_redis "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/external/redis"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/utils"
	identity_validator "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/utils/validation"
	usecase "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases"
	authCommands "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases/auth/commands"
	userCommands "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases/users/commands"
	userEvents "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases/users/events/domain-events"
	userQueries "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases/users/queries"
)

type Application struct {
	r *gin.Engine
}

func NewApplication() (*Application, error) {
	if v, err := validator.InitValidator(); err != nil {
		logger.Log.Fatal().Err(err).Msg("❌ Validator init failed")
		return nil, err
	} else {
		identity_validator.AddCustomValidation(v)
	}

	r := gin.Default()

	if err := db.InitDB(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Database init failed")
	}

	redisClient := ext_redis.NewRedisClient()

	// external services
	cacheService := cache.NewRedisCacheService(redisClient)
	jwtService := auth_pkg.NewJwtService(
		cacheService,
		config.AppConfig.Security.Jwt.Secret,
		config.AppConfig.Security.Jwt.EncryptionKey,
		config.AppConfig.Security.Jwt.Issuer,
		time.Duration(config.AppConfig.Security.Jwt.AccessTokenTtl)*time.Second,
		time.Duration(config.AppConfig.Security.Jwt.RefreshTokenTtl)*time.Second,
	)
	hashService := hash_pkg.NewHashService()

	rootDir := utils.GetWorkingDir()

	mailPath := path.Join(rootDir, "logs/identity/mail.log")
	mailLogger := logger.NewLogger(logger.LoggerConfig{
		Level:      "info",
		Filename:   mailPath,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		AppEnv:     config.AppConfig.Server.AppEnv,
	})
	mailFactory, err := mail.NewMailFactory(mail.GoogleMailProvider)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("❌ Mail factory init failed")
	}
	mailService := mail.NewMailService(mailFactory, mailLogger)

	// rabbitmq
	messageQueueService := rabbitmq.NewRabbitMQService(
		rabbitmq.MessageQueueConfig{
			Host:     config.AppConfig.MessageQueue.RabbitMQ.Host,
			Port:     config.AppConfig.MessageQueue.RabbitMQ.Port,
			User:     config.AppConfig.MessageQueue.RabbitMQ.User,
			Password: config.AppConfig.MessageQueue.RabbitMQ.Password,
		}, logger.Log)

	// repositories
	userRepository := repository.NewUserRepository(db.DB)

	// event bus
	var domainEventBus = events.NewDomainEventBus()

	usecases := &usecase.Usecase{
		Commands: usecase.Commands{
			LoginAccountHandler:    authCommands.NewLoginAccountHandler(userRepository, jwtService, hashService),
			RegisterAccountHandler: authCommands.NewRegisterAccountHandler(userRepository, domainEventBus),
			ForgotPasswordHandler:  authCommands.NewForgotPasswordHandler(messageQueueService, cacheService, mailService),
			AddAddressHandler:      userCommands.NewAddAddressHandler(),
			UpdateProfileHandler:   userCommands.NewUpdateProfileHandler(),
		},
		Queries: usecase.Queries{
			GetUsersHandler:      userQueries.NewGetUsersHandler(),
			GetUserByUUIDHandler: userQueries.NewGetUserByUUIDHandler(userRepository),
		},
		Events: usecase.DomainEvents{
			UserCreatedDomainHandler: userEvents.NewUserCreatedDomainEventHandler(domainEventBus),
		},
	}

	v1authHandler := v1handler.NewAuthRouteHandler(usecases)
	v1userHandler := v1handler.NewUserRouteHandler(usecases)

	routeList := []routes.Route{
		v1routes.NewAuthRoutes(v1authHandler),
		v1routes.NewUserRoutes(v1userHandler),
	}

	routes.RegisterRoutes(r, usecases, routeList...)

	return &Application{
		r: r,
	}, nil

}

func (app *Application) Run() error {
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.AppConfig.Server.Port),
		Handler: app.r,
	}

	osSignalChannel := make(chan os.Signal, 1)
	// syscall.SIGINT -> Ctrl + C
	// syscall.SIGTERM -> Kill command
	// syscall.SIGHUP -> Terminal closed
	signal.Notify(osSignalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		logger.Log.Info().Msgf("🚀 Server is running on port %d", config.AppConfig.Server.Port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Log.Error().Err(err).Msg("⛔️ Failed to start server")
		}
	}()

	<-osSignalChannel
	logger.Log.Warn().Msg("⚠️  Shutdown signal received ...")

	// Create a context with a timeout to allow graceful shutdown
	// Set the timeout duration to 15 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Error().Err(err).Msg("⛔️ Server forced to shutdown")
	}

	logger.Log.Info().Msg("🍺 Server exited gracefully")

	return nil
}
