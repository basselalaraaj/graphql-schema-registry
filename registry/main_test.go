package registry_test

import (
	"os"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/basselalaraaj/graphql-schema-registry/registry"
	"github.com/go-redis/redis/v8"
)

var (
	client *redis.Client
)

func TestMain(m *testing.M) {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	err = mr.StartAddr("localhost:6379")
	if err != nil {
		panic(err)
	}

	client = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	code := m.Run()

	defer func() {
		mr.Close()
		os.Exit(code)
	}()
}

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

		ans := schema.Save()
		if ans != nil {
			t.Fail()
		}
	})
}

func TestGetServiceSchema(t *testing.T) {
	t.Run("Get service schema from redis", func(t *testing.T) {
		_, err := registry.GetServiceSchema("Cart")
		if err != nil {
			t.Fail()
		}
	})
	t.Run("Return error if service schema is not found in redis", func(t *testing.T) {
		_, err := registry.GetServiceSchema("Shop")
		if err == nil {
			t.Fail()
		}
	})
}

func TestGetAllServices(t *testing.T) {
	t.Run("Get all services schema from redis", func(t *testing.T) {
		_, err := registry.GetAllServices()
		if err != nil {
			t.Fail()
		}
	})
}
