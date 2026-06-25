package domain

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/events"
	domainevents "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases/users/events/domain-events"
)

type UserCreatedDomainEvent struct {
	topic events.DomainEventTopic
	data  any
}

func (e UserCreatedDomainEvent) GetTopic() events.DomainEventTopic {
	return e.topic
}

func (e UserCreatedDomainEvent) GetPayload() []byte {
	bytes, _ := json.Marshal(e.data)
	return bytes
}

type User struct {
	UUID         uuid.UUID
	domainEvents []events.DomainEvent
}

func NewUser(uuid uuid.UUID) *User {

	user := &User{
		UUID: uuid,
	}

	domainEvent := UserCreatedDomainEvent{
		topic: domainevents.UserCreatedTopic,
		data:  user,
	}

	user.domainEvents = append(user.domainEvents, domainEvent)

	return user

}

func (u *User) PublishDomainEvent(eventBus *events.EventBus) {

	if eventBus == nil {
		panic("nil eventbus")
	}

	if len(u.domainEvents) == 0 {
		panic("No available domain events to publish")
	}

	for _, domainEvent := range u.domainEvents {
		eventBus.Publish(domainEvent.GetTopic(), domainEvent.GetPayload())
	}

}
