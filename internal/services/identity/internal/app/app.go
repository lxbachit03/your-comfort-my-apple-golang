package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/jwt"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/logger"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/security/hash"
	validator "github.com/lxbachit03/ygz-microservices-golang/internal/pkg/validator"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/config"
	v1handler "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/endpoints/handlers/v1"
	routes "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/endpoints/routes"
	v1routes "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/endpoints/routes/v1"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/db"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/db/repository"
	identity_validator "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/utils/validation"
	usecase "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases"
	command "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases/auth/commands"
	usersCommand "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases/users/commands"
	query "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases/users/queries"
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

	// external services
	jwtService := jwt.NewJwtService()
	hashService := hash.NewHashService()

	// repositories
	userRepository := repository.NewUserRepository(db.DB)

	usecases := &usecase.Usecase{
		Commands: usecase.Commands{
			LoginAccountHandler:    command.NewLoginAccountHandler(userRepository, jwtService, hashService),
			RegisterAccountHandler: command.NewRegisterAccountHandler(userRepository),
			AddAddressHandler:      usersCommand.NewAddAddressHandler(),
			UpdateProfileHandler:   usersCommand.NewUpdateProfileHandler(),
		},
		Queries: usecase.Queries{
			GetUsersHandler:      query.NewGetUsersHandler(),
			GetUserByUUIDHandler: query.NewGetUserByUUIDHandler(userRepository),
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
