package eventbus

import (
	"fmt"
	"github.com/cross-space-official-private/common/logger"
	"github.com/cross-space-official-private/common/utils"
)

func Initialize() {
	initializeChannel()
	go listenEvent()
}

func listenEvent() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
		go listenEvent()
	}()

	for {
		select {
		case event := <-domainEventChannel.channel:
			logger.GetLoggerEntry(event.GetContext()).Info("Handling domain event: ", event.GetIdentifier())
			for _, handler := range domainEventChannel.handlers {
				if handler.CanHandle(event) {
					err := handler.Handle(event)
					if !utils.IsNil(err) {
						logger.GetLoggerEntry(event.GetContext()).Error(err)
					}
				}
			}
		}
	}
}
