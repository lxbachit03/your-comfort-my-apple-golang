package usecase

import (
	authCommand "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases/auth/commands"
	userCommand "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases/users/commands"
	userDomainEvent "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases/users/events/domain-events"
	userQuery "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases/users/queries"
)

type Usecase struct {
	Commands Commands
	Queries  Queries
	Events   DomainEvents
}

type Commands struct {
	LoginAccountHandler    authCommand.LoginAccountHandler
	RegisterAccountHandler authCommand.RegisterAccountHandler
	ForgotPasswordHandler  authCommand.ForgotPasswordHandler
	AddAddressHandler      userCommand.AddAddressHandler
	UpdateProfileHandler   userCommand.UpdateProfileHandler
}

type Queries struct {
	GetUsersHandler      userQuery.GetUsersHandler
	GetUserByUUIDHandler userQuery.GetUserByUUIDHandler
}

type DomainEvents struct {
	UserCreatedDomainHandler userDomainEvent.UserCreatedDomainEventHandler
}
