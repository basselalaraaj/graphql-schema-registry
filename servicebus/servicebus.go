package servicebus

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/basselalaraaj/graphql-schema-registry/registry"
)

var seconds int = 40

type client struct {
	topic *servicebus.Topic
}

var serviceBusClient *client

func (s *client) createClient() error {
	connectionString := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	if connectionString == "" {
		return fmt.Errorf("FATAL: expected environment variable SERVICEBUS_CONNECTION_STRING not set")
	}

	// Create a client to communicate with a Service Bus Namespace.
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connectionString))
	if err != nil {
		return err
	}

	topicName := os.Getenv("SERVICEBUS_TOPIC_NAME")
	if topicName == "" {
		return fmt.Errorf("FATAL: expected environment variable SERVICEBUS_TOPIC_NAME not set")
	}

	client, err := ns.NewTopic(topicName)
	if err != nil {
		return err
	}
	s.topic = client

	return nil
}

// Initialize the service bus
func Initialize() {
	serviceBusEnabled := os.Getenv("SERVICEBUS_ENABLED")
	if serviceBusEnabled == "" {
		return
	}
	serviceBusClient = &client{}
	err := serviceBusClient.createClient()
	if err != nil {
		fmt.Println("not able to create a client")
		return
	}
}

// SendMessage to send messages on the service bus
func SendMessage(message *registry.SchemaRegistry) error {
	serviceBusEnabled := os.Getenv("SERVICEBUS_ENABLED")
	if serviceBusEnabled == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	defer cancel()

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	if err := serviceBusClient.topic.Send(ctx, servicebus.NewMessageFromString(string(jsonMessage))); err != nil {
		return err
	}

	return nil
}
