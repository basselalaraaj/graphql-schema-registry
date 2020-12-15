package registry

import (
	"context"
	"encoding/json"

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

var rdb = redis.NewClient(&redis.Options{
	Addr: ":6379",
})

var ctx = context.Background()

// ValidateSchema validates the graphql schema
func (s *SchemaRegistry) ValidateSchema() error {
	_, err := gqlparser.LoadSchema(&ast.Source{Name: s.ServiceName, Input: s.TypeDefs, BuiltIn: false})
	if err != nil {
		return err
	}
	return nil
}

// Save the schemaa in redis
func (s *SchemaRegistry) Save() error {
	value, _ := json.Marshal(s)

	err := rdb.Set(ctx, s.ServiceName, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
