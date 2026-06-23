package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/logger"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/config"
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

	if err := config.LoadConfig(configPath, env); err != nil {
		logger.Log.Fatal().Msg("❌ Unable to load config")
	}

	// Start worker
	worker := NewWorker()
	if worker == nil {
		logger.Log.Fatal().Msg("Failed to create worker")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := worker.StartWorker(ctx); err != nil && err != context.Canceled {
			logger.Log.Error().Err(err).Msg("Worker failed to start")
		}
	}()

	<-ctx.Done()
	logger.Log.Info().Msg("Received shutdown signal")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := worker.ShutdownWorker(shutdownCtx); err != nil {
		logger.Log.Error().Err(err).Msg("Shutdown failed")
	}

	wg.Wait()
	logger.Log.Info().Msg("Main process terminated")
}
