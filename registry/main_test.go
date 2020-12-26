package registry_test

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/basselalaraaj/graphql-schema-registry/registry"
	"github.com/elliotchance/redismock/v8"
	"github.com/go-redis/redis/v8"
)

func TestSchemaValidation(t *testing.T) {
	t.Run("Should be successful if schema is valid", func(t *testing.T) {
		schema := registry.SchemaRegistry{
			ServiceName: "Cart",
			ServiceURL:  "http://cart-service",
			TypeDefs:    "type Query { placeHolder: String }",
		}

		ans := schema.ValidateSchema()
		if ans != nil {
			t.Fail()
		}
	})

	t.Run("Should return an error if schema is invalid", func(t *testing.T) {
		schema := registry.SchemaRegistry{
			ServiceName: "Cart",
			ServiceURL:  "http://cart-service",
			TypeDefs:    "type Query { placeHolder: random }",
		}

		ans := schema.ValidateSchema()
		if ans == nil {
			t.Fail()
		}
	})
}

func TestSchemaSave(t *testing.T) {
	t.Run("Should save the schema correctly to redis", func(t *testing.T) {
		schema := registry.SchemaRegistry{
			ServiceName: "Cart",
			ServiceURL:  "http://cart-service",
			TypeDefs:    "type Query { placeHolder: String }",
		}

		mr, err := miniredis.Run()
		if err != nil {
			panic(err)
		}
		defer mr.Close()

		err = mr.StartAddr("localhost:6379")
		if err != nil {
			panic(err)
		}

		client := redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		})

		mock := redismock.NewNiceMock(client)
		mock.On("Set").Return(redis.NewStatusResult("", nil))

		ans := schema.Save()
		if ans != nil {
			t.Fail()
		}
		mr.Close()
	})
}
