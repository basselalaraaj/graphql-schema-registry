package servicebus_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/basselalaraaj/graphql-schema-registry/registry"
	"github.com/basselalaraaj/graphql-schema-registry/servicebus"
)

var (
	connectionString string = "Endpoint=sb://localhost.servicebus.windows.net/;SharedAccessKeyName=123;SharedAccessKey=123"
)

func TestMain(m *testing.M) {
	code := m.Run()

	defer func() {
		os.Exit(code)
	}()
}

func TestInitialize(t *testing.T) {
	t.Run("Should throw an error that configuration 'SERVICEBUS_CONNECTION_STRING' is missing", func(t *testing.T) {
		servicebus.Initialize()
		if servicebus.ServiceBusClient.Topic != nil {
			t.Fail()
		}
	})
	t.Run("Should throw an error that configuration 'SERVICEBUS_TOPIC_NAME' is missing", func(t *testing.T) {
		os.Setenv("SERVICEBUS_CONNECTION_STRING", connectionString)
		servicebus.Initialize()
		if servicebus.ServiceBusClient.Topic != nil {
			fmt.Println(servicebus.ServiceBusClient)
			t.Fail()
		}
	})
	t.Run("Should create a client", func(t *testing.T) {
		os.Setenv("SERVICEBUS_CONNECTION_STRING", connectionString)
		os.Setenv("SERVICEBUS_TOPIC_NAME", "abc")
		servicebus.Initialize()
		if servicebus.ServiceBusClient.Topic == nil {
			t.Fail()
		}
	})
}

func TestSendMessage(t *testing.T) {
	t.Run("Should send messages correctly", func(t *testing.T) {
		os.Setenv("SERVICEBUS_CONNECTION_STRING", connectionString)
		os.Setenv("SERVICEBUS_TOPIC_NAME", "abc")

		servicebus.Initialize()
		message := registry.SchemaRegistry{
			ServiceName: "Cart",
			ServiceURL:  "http://cart-service",
			TypeDefs:    "type Query { placeHolder: String }",
		}
		if servicebus.ServiceBusClient.Topic == nil {
			t.Fail()
		}

		err := servicebus.SendMessage(&message)

		if err == nil {
			t.Fail()
		}
	})
}
