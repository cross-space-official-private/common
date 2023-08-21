package eventbus

import (
	"context"
	"github.com/cross-space-official-private/common/logger"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type (
	TestDomainEvent struct {
	}
	TestDomainEventHandler struct {
		Counter int
	}
)

func (t *TestDomainEvent) GetTimestamp() time.Time {
	//TODO implement me
	panic("implement me")
}

func (t *TestDomainEventHandler) CanHandle(event DomainEvent) bool {
	return true
}

func (t *TestDomainEventHandler) Handle(event DomainEvent) error {
	logger.GetLoggerEntry(event.GetContext()).Info("Handling domain event: ", event.GetIdentifier())
	t.Counter++

	panic("test panic")
	return nil
}

func (t *TestDomainEvent) GetContext() context.Context {
	return context.TODO()
}

func (t *TestDomainEvent) GetPayload() interface{} {
	return context.TODO()
}

func (t *TestDomainEvent) GetIdentifier() string {
	return "test"
}

func (t *TestDomainEvent) WithContext(ctx context.Context) DomainEvent {
	return t
}

func Test(t *testing.T) {
	Initialize()

	SendDomainEvent(context.Background(), &TestDomainEvent{})

	time.Sleep(1 * time.Second)

	handler := &TestDomainEventHandler{Counter: 0}
	RegisterHandlers(handler)
	SendDomainEvent(context.Background(), &TestDomainEvent{})
	time.Sleep(1 * time.Second)

	SendDomainEvent(context.Background(), &TestDomainEvent{})
	assert.Equal(t, handler.Counter, 1)

	time.Sleep(1 * time.Second)
	assert.Equal(t, handler.Counter, 2)
}
