package servicebus

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/basselalaraaj/graphql-schema-registry/registry"
	"github.com/joho/godotenv"
)

var seconds int = 40

type client struct {
	topic *servicebus.Topic
}

// ServiceBus client
type ServiceBus struct {
}

var serviceBusClient *client

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	serviceBusClient = getClient()
}

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

func getClient() *client {
	client := &client{}

	serviceBusEnabled := os.Getenv("SERVICEBUS_ENABLED")
	if serviceBusEnabled == "" {
		return client
	}
	err := client.createClient()
	if err != nil {
		log.Fatal("not able to create a client")
	}

	return client
}

// SendNotification send message to the bus
func (s *ServiceBus) SendNotification(message *registry.SchemaRegistry) error {
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
