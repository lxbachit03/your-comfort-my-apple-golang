package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/joho/godotenv"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/config"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/app"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/utils"
)

func main() {
	rootDir := utils.GetWorkingDir()

	// `make start-identity ENV=local` -> env = "local"
	// `make start-identity ENV=dev` -> env = "dev"
	// `make start-identity ENV=test` -> env = "test"
	// `make start-identity ENV=staging` -> env = "staging"
	// `make start-identity ENV=prod` -> env = "prod"
	// `make start-identity` -> env = "local" (default)
	var env = os.Getenv("ENV")
	if env == "" {
		env = "local"
	}
	var envFile = fmt.Sprintf(".env.%s", env)
	var envPath = path.Join(rootDir, envFile)

	// Load local environment variables
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("❌ Unable to load env file: %s", err)
	}

	// Load config
	var configPath = path.Join(rootDir, "internal/services/identity/config")
	if err := config.LoadConfig(configPath, env); err != nil {
		log.Fatalf("❌ Unable to load config: %v", err)
	}

	log.Print("test:", config.AppConfig.Server.Port)
	log.Print("test 2:", config.AppConfig.Database.Name)

	app, err := app.NewApplication()
	if err != nil {
		log.Fatal("❌ Unable to create application: ", err)
	}

	if err := app.Run(); err != nil {
		log.Fatal("❌ Unable to run application: ", err)
	}
}
