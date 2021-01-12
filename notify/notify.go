package notify

import (
	"github.com/basselalaraaj/graphql-schema-registry/registry"
)

type notifier interface {
	SendNotification(message *registry.SchemaRegistry) error
}

// SendNotification send notification
func SendNotification(message *registry.SchemaRegistry, n notifier) error {
	return n.SendNotification(message)
}
