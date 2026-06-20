package main

import (
	"context"
	"encoding/json"
	"net/http"
	"path"

	apiresponse "github.com/lxbachit03/ygz-microservices-golang/internal/pkg/api-response"
	errorcode "github.com/lxbachit03/ygz-microservices-golang/internal/pkg/error_code"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/logger"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/queue/rabbitmq"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/config"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/external/mail"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/utils"
	"github.com/rs/zerolog"
)

type Worker struct {
	rabbitMQ    rabbitmq.MessageQueueService
	mailService mail.MailService
	logger      *zerolog.Logger
}

func NewWorker() *Worker {
	rootDir := utils.GetWorkingDir()
	workerLogPath := path.Join(rootDir, "logs/identity/worker.log")

	workerLogger := logger.NewLogger(logger.LoggerConfig{
		Level:      "info",
		Filename:   workerLogPath,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		AppEnv:     config.AppConfig.Server.AppEnv,
	})

	// init rabbitmq
	messageQueueService := rabbitmq.NewRabbitMQService(
		rabbitmq.MessageQueueConfig{
			Host:     config.AppConfig.MessageQueue.RabbitMQ.Host,
			Port:     config.AppConfig.MessageQueue.RabbitMQ.Port,
			User:     config.AppConfig.MessageQueue.RabbitMQ.User,
			Password: config.AppConfig.MessageQueue.RabbitMQ.Password,
		}, logger.Log)

	// init mail service
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

	return &Worker{
		rabbitMQ:    messageQueueService,
		mailService: mailService,
		logger:      workerLogger,
	}
}

func (w *Worker) StartWorker(ctx context.Context) error {
	const emailQueueName = "reset_password_queue"

	handler := func(body []byte) error {

		w.logger.Debug().Msgf("Received message: %s", string(body))

		var email mail.Email
		if err := json.Unmarshal(body, &email); err != nil {
			w.logger.Error().Err(err).Msg("Failed to unmarshal message")
			return err
		}

		if err := w.mailService.SendMail(ctx, &email); err != nil {
			return &apiresponse.ApiError{
				StatusCode:   http.StatusInternalServerError,
				ErrorCode:    errorcode.ErrInternalServer,
				Message:      "Failed to send email",
				ErrorDetails: err,
			}
		}

		w.logger.Info().Msgf("Email sent successfully to %v", email.To)

		return nil
	}

	if err := w.rabbitMQ.Consume(ctx, emailQueueName, handler); err != nil {
		w.logger.Error().Err(err).Msg("Failed to satrt consumer")
		return err
	}

	w.logger.Info().Msgf("Worker started, consuming from queue: %s", emailQueueName)

	<-ctx.Done()
	w.logger.Info().Msgf("Worker stopped consuming due to context cancellation")

	return ctx.Err()
}

func (w *Worker) ShutdownWorker(ctx context.Context) error {
	w.logger.Info().Msg("Shutting down worker ... ")

	if err := w.rabbitMQ.Close(); err != nil {
		w.logger.Error().Err(err).Msg("Failed to close RabbitMQ")
		return err
	}

	w.logger.Info().Msg("RabbitMQ connection closed successfully")

	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			w.logger.Warn().Msg("Shutdown timeout exceeded")
			return ctx.Err()
		}
	default:
	}

	w.logger.Info().Msg("Worker shutdown completed")

	return nil
}
