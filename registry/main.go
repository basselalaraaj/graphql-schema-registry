package registry

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

// SchemaRegistry to add a schema to the registry
type SchemaRegistry struct {
	ServiceName string
	ServiceURL  string
	TypeDefs    string
}

var redisDB = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

var ctx = context.Background()

// ValidateSchema validates the graphql schema
func (s *SchemaRegistry) ValidateSchema() error {
	if s.TypeDefs == "" {
		return fmt.Errorf("typedefs should not be empty")
	}
	_, err := gqlparser.LoadSchema(&ast.Source{Name: s.ServiceName, Input: s.TypeDefs, BuiltIn: false})
	if err != nil {
		return err
	}
	return nil
}

// Save the schema in redis
func (s *SchemaRegistry) Save() error {
	value, _ := json.Marshal(s)

	err := redisDB.Set(ctx, s.ServiceName, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetServiceSchema get service schema from redis
func GetServiceSchema(service string) (*SchemaRegistry, error) {
	val2, err := redisDB.Get(ctx, service).Result()
	if err == redis.Nil {
		return &SchemaRegistry{}, err
	} else if err != nil {
		return &SchemaRegistry{}, err
	} else {
		var serviceSchema SchemaRegistry
		err := json.Unmarshal([]byte(val2), &serviceSchema)
		if err != nil {
			return &SchemaRegistry{}, err
		}
		return &serviceSchema, nil
	}
}

// GetAllServices returns all services names
func GetAllServices() (*[]string, error) {
	var cursor uint64
	keys, _, err := redisDB.Scan(ctx, cursor, "*", 100).Result()
	if err != nil {
		return &[]string{}, err
	}
	return &keys, nil
}
