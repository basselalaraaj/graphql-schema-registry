package registry

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
)

var (
	db   *redis.Client
	mock redismock.ClientMock
)

func TestMain(m *testing.M) {
	db, mock = redismock.NewClientMock()
	redisDB = db

	code := m.Run()

	defer func() {
		os.Exit(code)
	}()
}

func TestSchemaValidation(t *testing.T) {
	t.Run("Should be successful if schema is valid", func(t *testing.T) {
		schema := SchemaRegistry{
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
		schema := SchemaRegistry{
			ServiceName: "Cart",
			ServiceURL:  "http://cart-service",
			TypeDefs:    "type Query { placeHolder: random }",
		}

		ans := schema.ValidateSchema()
		if ans == nil {
			t.Fail()
		}
	})

	t.Run("Should return an error if schema is empty", func(t *testing.T) {
		schema := SchemaRegistry{
			ServiceName: "Cart",
			ServiceURL:  "http://cart-service",
			TypeDefs:    "",
		}

		ans := schema.ValidateSchema()
		if ans == nil {
			t.Fail()
		}
	})
}

func TestSchemaSave(t *testing.T) {
	t.Run("Should save the schema correctly to redis", func(t *testing.T) {
		schema := SchemaRegistry{
			ServiceName: "Cart",
			ServiceURL:  "http://cart-service",
			TypeDefs:    "type Query { placeHolder: String }",
		}
		value, _ := json.Marshal(schema)

		mock.ExpectSet("Cart", value, 0).SetVal("{}")
		ans := schema.Save()
		if ans != nil {
			t.Fail()
		}
	})

	t.Run("Should throw error while saving the schema", func(t *testing.T) {
		schema := SchemaRegistry{
			ServiceName: "Cart",
			ServiceURL:  "http://cart-service",
			TypeDefs:    "type Query { placeHolder: String }",
		}
		value, _ := json.Marshal(schema)

		mock.ExpectSet("Cart", value, 0).SetErr(errors.New("fail"))
		ans := schema.Save()
		if ans == nil {
			t.Fail()
		}
	})
}

func TestGetServiceSchema(t *testing.T) {
	t.Run("Get service schema from redis", func(t *testing.T) {
		mock.ExpectGet("Cart").SetVal("{}")
		_, err := GetServiceSchema("Cart")

		if err != nil {
			t.Fail()
		}
	})

	t.Run("Return error if service schema is not found in redis", func(t *testing.T) {
		mock.ExpectGet("Cart").SetErr(errors.New("fail"))
		_, err := GetServiceSchema("Cart")

		if err == nil {
			t.Fail()
		}
	})

	t.Run("Return redis nil if service schema is not found in redis", func(t *testing.T) {
		mock.ExpectGet("Cart").RedisNil()
		_, err := GetServiceSchema("Cart")

		if err == nil {
			t.Fail()
		}
	})
}

func TestGetAllServices(t *testing.T) {
	t.Run("Get all services schema from redis", func(t *testing.T) {
		var cursor uint64
		mock.ExpectScan(cursor, "*", 100).SetVal([]string{"{}"}, cursor)
		_, err := GetAllServices()

		if err != nil {
			t.Fail()
		}
	})

	t.Run("Get all services schema from redis", func(t *testing.T) {
		var cursor uint64
		mock.ExpectScan(cursor, "*", 100).SetErr(errors.New("fail"))
		_, err := GetAllServices()

		if err == nil {
			t.Fail()
		}
	})
}
