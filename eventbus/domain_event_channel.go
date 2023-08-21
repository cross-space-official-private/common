package eventbus

import (
	"context"
	"github.com/cross-space-official-private/common/logger"
	"time"
)

var domainEventChannel *DomainEventChannel

type (
	DomainEvent interface {
		GetContext() context.Context
		GetPayload() interface{}
		GetIdentifier() string
		GetTimestamp() time.Time
		WithContext(ctx context.Context) DomainEvent
	}

	DomainEventHandler interface {
		CanHandle(event DomainEvent) bool
		Handle(event DomainEvent) error
	}

	DomainEventChannel struct {
		channel  chan DomainEvent
		handlers []DomainEventHandler
	}
)

func initializeChannel() {
	domainEventChannel = &DomainEventChannel{
		channel:  make(chan DomainEvent),
		handlers: make([]DomainEventHandler, 0),
	}
}

func GetDomainEventChannel() *DomainEventChannel {
	return domainEventChannel
}

func RegisterHandlers(handlers ...DomainEventHandler) {
	if domainEventChannel == nil {
		return
	}

	for _, handler := range handlers {
		domainEventChannel.registerHandler(handler)
	}
}

func SendDomainEvent(ctx context.Context, event DomainEvent) {
	if domainEventChannel == nil {
		logger.GetLoggerEntry(ctx).Warn("The domain event channel is not initialized")
		return
	}

	domainEventChannel.sendDomainEvent(ctx, event)
}

func (s *DomainEventChannel) registerHandler(handler DomainEventHandler) {
	s.handlers = append(s.handlers, handler)
}

func (s *DomainEventChannel) sendDomainEvent(ctx context.Context, event DomainEvent) {
	s.channel <- event.WithContext(ctx)
}
