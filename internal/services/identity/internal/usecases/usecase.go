package usecase

import (
	command "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases/auth/commands/login_account"
	query "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases/users/queries"
)

type Usecase struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	LoginAccountHandler command.LoginAccountHandler
}

type Queries struct {
	GetUsersHandler      query.GetUsersHandler
	GetUserByUUIDHandler query.GetUserByUUIDHandler
}
