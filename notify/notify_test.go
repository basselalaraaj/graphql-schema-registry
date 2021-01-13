package notify

import (
	"testing"

	"github.com/basselalaraaj/graphql-schema-registry/registry"
)

func TestSendNotification(t *testing.T) {
	t.Run("Test send notification", func(t *testing.T) {
		message := registry.SchemaRegistry{
			ServiceName: "Cart",
			ServiceURL:  "http://cart-service",
			TypeDefs:    "type Query { placeHolder: String }",
		}
		err := SendNotification(&message)
		if err != nil {
			t.Fail()
		}
	})
}
