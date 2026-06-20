package rabbitmq

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/rabbitmq/amqp091-go"
)

type MessageQueueConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

type rabbitMQService struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	logger  *zerolog.Logger
}

type MessageQueueService interface {
	Publish() error
	Comsume() error
	Close() error
}

func NewRabbitMQService(config MessageQueueConfig, logger *zerolog.Logger) MessageQueueService {
	amqpURL := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		config.User,
		config.Password,
		config.Host,
		config.Port,
	)

	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		// logger.Error().Err(err).Msg("failed to connect to RabbitMQ")
		panic(err.Error())
		// return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		panic(err.Error())
		// logger.Error().Err(err).Msg("failed to open channel")
		// return nil
	}

	logger.Info().Msg("🍺 Connected RabbitMQ")

	return &rabbitMQService{
		conn:    conn,
		channel: ch,
		logger:  logger,
	}
}

func (r *rabbitMQService) Publish() error {
	return nil
}

func (r *rabbitMQService) Comsume() error {
	return nil
}

func (r *rabbitMQService) Close() error {
	return nil
}
