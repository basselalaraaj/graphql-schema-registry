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

type serviceBus struct {
	connectionString string
	client           *servicebus.Topic
}

var serviceBusConnection *serviceBus

func (s *serviceBus) startConnection() {
	connectionString := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	if connectionString == "" {
		fmt.Println("FATAL: expected environment variable SERVICEBUS_CONNECTION_STRING not set")
		return
	}
	s.connectionString = connectionString
}

func (s *serviceBus) createClient() {
	// Create a client to communicate with a Service Bus Namespace.
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(s.connectionString))
	if err != nil {
		fmt.Println(err)
		return
	}

	topicName := os.Getenv("SERVICEBUS_TOPIC_NAME")
	if topicName == "" {
		fmt.Println("FATAL: expected environment variable SERVICEBUS_TOPIC_NAME not set")
		return
	}

	client, err := ns.NewTopic(topicName)
	if err != nil {
		return
	}
	s.client = client
}

// Initialize the service bus
func Initialize() {
	serviceBusConnection = &serviceBus{}

	serviceBusConnection.startConnection()
	serviceBusConnection.createClient()
}

// SendMessage to send messages on the service bus
func SendMessage(message *registry.SchemaRegistry) {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := serviceBusConnection.client.Send(ctx, servicebus.NewMessageFromString(string(jsonMessage))); err != nil {
		fmt.Println("FATAL: ", err)
	}
}
