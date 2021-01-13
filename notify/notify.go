package notify

import (
	"github.com/basselalaraaj/graphql-schema-registry/registry"
	"github.com/basselalaraaj/graphql-schema-registry/servicebus"
)

type notifier interface {
	SendNotification(message *registry.SchemaRegistry) error
}

func sendMessage(message *registry.SchemaRegistry, n notifier) error {
	return n.SendNotification(message)
}

// SendNotification send notification
func SendNotification(message *registry.SchemaRegistry) error {
	serviceBus := servicebus.ServiceBus{}
	return sendMessage(message, &serviceBus)
}
