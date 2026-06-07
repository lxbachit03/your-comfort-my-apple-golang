package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1handler "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/endpoints/handlers/v1"
	route "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/endpoints/routes"
	v1routes "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/endpoints/routes/v1"
	usecase "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases"
	command "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases/auth/commands/login_account"
)

type Application struct {
	r *gin.Engine
}

func NewApplication() (*Application, error) {

	r := gin.Default()

	// init dependencies
	// logger
	// database

	usecases := &usecase.Usecase{
		Commands: usecase.Commands{
			LoginAccountHandler: command.NewLoginAccountHandler(),
		},
	}

	v1authHandler := v1handler.NewAuthRouteHandler(usecases)
	v1userHandler := v1handler.NewUserRouteHandler(usecases)

	routeList := []route.Route{
		v1routes.NewAuthRoutes(v1authHandler),
		v1routes.NewUserRoutes(v1userHandler),
	}

	route.RegisterRoutes(r, usecases, routeList...)

	return &Application{
		r: r,
	}, nil

}

func (app *Application) Run() error {

	server := http.Server{
		Addr:    ":8080",
		Handler: app.r,
	}

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}
