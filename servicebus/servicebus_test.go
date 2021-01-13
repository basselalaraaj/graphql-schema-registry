package servicebus

import (
	"fmt"
	"os"
	"testing"

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

func TestInitialize(t *testing.T) {
	t.Run("Should throw an error that configuration 'SERVICEBUS_CONNECTION_STRING' is missing", func(t *testing.T) {
		os.Setenv("SERVICEBUS_ENABLED", "True")
		Initialize()
		if serviceBusClient.topic != nil {
			t.Fail()
		}
	})
	t.Run("Should throw an error that configuration 'SERVICEBUS_TOPIC_NAME' is missing", func(t *testing.T) {
		os.Setenv("SERVICEBUS_CONNECTION_STRING", connectionString)
		Initialize()
		if serviceBusClient.topic != nil {
			fmt.Println(serviceBusClient)
			t.Fail()
		}
	})
	t.Run("Should create a client", func(t *testing.T) {
		os.Setenv("SERVICEBUS_CONNECTION_STRING", connectionString)
		os.Setenv("SERVICEBUS_TOPIC_NAME", "abc")
		Initialize()
		if serviceBusClient.topic == nil {
			t.Fail()
		}
	})
}

func TestSendNotification(t *testing.T) {
	t.Run("Should send messages correctly", func(t *testing.T) {
		os.Setenv("SERVICEBUS_CONNECTION_STRING", connectionString)
		os.Setenv("SERVICEBUS_TOPIC_NAME", "abc")

		Initialize()
		message := registry.SchemaRegistry{
			ServiceName: "Cart",
			ServiceURL:  "http://cart-service",
			TypeDefs:    "type Query { placeHolder: String }",
		}
		if serviceBusClient.topic == nil {
			t.Fail()
		}

		serviceBus := ServiceBus{}

		err := serviceBus.SendNotification(&message)

		if err == nil {
			t.Fail()
		}
	})
}
