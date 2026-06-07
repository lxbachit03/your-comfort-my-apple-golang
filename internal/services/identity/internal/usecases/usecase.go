package usecase

import command "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases/auth/commands/login_account"

type Usecase struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	LoginAccountHandler command.LoginAccountHandler
}

type Queries struct {
}
