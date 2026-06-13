package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/joho/godotenv"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/logger"
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

	log.Printf("ENV: %s", env)

	if env == "" {
		env = "local"
	}
	var envFile = fmt.Sprintf(".env.%s", env)
	var envPath = path.Join(rootDir, envFile)

	log.Printf("envFile: %s", envFile)
	log.Printf("envPath: %s", envPath)

	// Init Application Logger
	logPath := path.Join(rootDir, "logs/identity/app.log")
	logger.NewApplicationLogger(logger.LoggerConfig{
		Level:      "info",
		Filename:   logPath,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		AppEnv:     env,
	})

	// Load local environment variables (overloading existing variables like those exported from .env.local via Makefile)
	if err := godotenv.Overload(envPath); err != nil {
		logger.Log.Warn().Msg("⚠️ Unable to load env file")
	} else {
		logger.Log.Info().Msg("✅ Environment variables loaded successfully")
	}

	// Load config
	var configPath = path.Join(rootDir, "internal/services/identity/config")

	log.Printf("env 123: %s", env)

	if err := config.LoadConfig(configPath, env); err != nil {
		logger.Log.Fatal().Msg("❌ Unable to load config")
	}

	app, err := app.NewApplication()
	if err != nil {
		logger.Log.Fatal().Msg("❌ Unable to create application")
	}

	if err := app.Run(); err != nil {
		logger.Log.Fatal().Msg("❌ Unable to run application")
	}
}
