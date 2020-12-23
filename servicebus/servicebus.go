package servicebus

import (
	"context"
	"fmt"
	"os"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
)

type serviceBus struct {
	connectionString string
	client           *servicebus.Queue
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

	queueName := os.Getenv("SERVICEBUS_QUEUE_NAME")
	if queueName == "" {
		fmt.Println("FATAL: expected environment variable SERVICEBUS_QUEUE_NAME not set")
		return
	}

	client, err := ns.NewQueue(queueName)
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
func SendMessage() {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	if err := serviceBusConnection.client.Send(ctx, servicebus.NewMessageFromString("Hello World!!!")); err != nil {
		fmt.Println("FATAL: ", err)
	}
}
