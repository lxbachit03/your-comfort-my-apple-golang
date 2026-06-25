package events

import "sync"

type Event struct {
	Payload interface{}
}

type DomainEvent interface {
	GetTopic() DomainEventTopic
	GetPayload() []byte
}

type DomainEventTopic string

type EventHandler func(Event)

type EventBus struct {
	handlers map[string][]func([]byte)
}

func NewDomainEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]func([]byte)),
	}
}

func (eb *EventBus) Subcribe(topic DomainEventTopic, handler func([]byte)) {
	eb.handlers[string(topic)] = append(eb.handlers[string(topic)], handler)
}

func (eb *EventBus) Publish(topic DomainEventTopic, data []byte) {
	handlers, exist := eb.handlers[string(topic)]
	if !exist {
		return
	}

	var wg sync.WaitGroup

	for _, h := range handlers {
		wg.Add(1)
		go func(handler func([]byte)) {
			defer wg.Done()
			handler(data)
		}(h)
	}

	wg.Wait()
}
