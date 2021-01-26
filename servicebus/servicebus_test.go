package servicebus

import (
	"context"
	"os"
	"testing"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/basselalaraaj/graphql-schema-registry/registry"
)

var (
	connectionString string = "Endpoint=sb://localhost.windows.net/;SharedAccessKeyName=123;SharedAccessKey=123"
)

func TestMain(m *testing.M) {
	code := m.Run()

	defer func() {
		os.Exit(code)
	}()
}

func TestCreateClient(t *testing.T) {
	t.Run("Should throw an error that configuration 'SERVICEBUS_CONNECTION_STRING' is missing", func(t *testing.T) {
		os.Setenv("SERVICEBUS_ENABLED", "True")

		_ = ServiceBusClient.CreateClient()

		if ServiceBusClient.topic != nil {
			t.Fail()
		}
	})
	t.Run("Should throw an error that configuration 'SERVICEBUS_TOPIC_NAME' is missing", func(t *testing.T) {
		os.Setenv("SERVICEBUS_CONNECTION_STRING", connectionString)
		_ = ServiceBusClient.CreateClient()

		if ServiceBusClient.topic != nil {
			t.Fail()
		}
	})
	t.Run("Should create a client", func(t *testing.T) {
		os.Setenv("SERVICEBUS_CONNECTION_STRING", connectionString)
		os.Setenv("SERVICEBUS_TOPIC_NAME", "abc")

		_ = ServiceBusClient.CreateClient()

		if ServiceBusClient.topic == nil {
			t.Fail()
		}
	})
}

type TopicMock struct {
}

func (t *TopicMock) Send(ctx context.Context, event *servicebus.Message, opts ...servicebus.SendOption) error {
	return nil
}

func TestSendNotification(t *testing.T) {
	t.Run("Should send messages correctly", func(t *testing.T) {
		os.Setenv("SERVICEBUS_ENABLED", "True")
		os.Setenv("SERVICEBUS_CONNECTION_STRING", connectionString)
		os.Setenv("SERVICEBUS_TOPIC_NAME", "abc")

		message := registry.SchemaRegistry{
			ServiceName: "Cart",
			ServiceURL:  "http://cart-service",
			TypeDefs:    "type Query { placeHolder: String }",
		}
		if ServiceBusClient.topic == nil {
			t.Fail()
		}

		ServiceBusClient = ServiceBus{
			topic: &TopicMock{},
		}

		err := ServiceBusClient.SendNotification(&message)

		if err != nil {
			t.Fail()
		}
	})
}
