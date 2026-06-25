package domainevents

import (
	"encoding/json"
	"log"

	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/events"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/logger"
)

const (
	UserCreatedTopic events.DomainEventTopic = "user.created"
)

type UserCreatedDomainEventHandler struct {
	eventBus *events.EventBus
}

func NewUserCreatedDomainEventHandler(eventBus *events.EventBus) UserCreatedDomainEventHandler {
	handler := UserCreatedDomainEventHandler{
		eventBus: eventBus,
	}

	eventBus.Subcribe(UserCreatedTopic, func(data []byte) {
		handler.Handle(data)
	})

	return handler
}

func (h UserCreatedDomainEventHandler) Handle(data []byte) {
	// data is already a JSON encoded byte slice from Publish
	var dummy interface{}
	err := json.Unmarshal(data, &dummy)
	if err != nil {
		logger.Log.Error().Err(err).Msg("❌ User created domain event marshal failed")
		return
	}

	log.Print("✅ User created domain event handled successfully")
}
