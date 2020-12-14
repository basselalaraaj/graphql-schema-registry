package graph

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/basselalaraaj/graphql-schema-registry/graph/generated"
	"github.com/basselalaraaj/graphql-schema-registry/graph/model"
	"github.com/go-redis/redis/v8"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

var ctx = context.Background()

// It serves as dependency injection for your app, add any dependencies you require here.
type Resolver struct{}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

type schemaRegistry struct {
	ServiceName string
	ServiceURL  string
	TypeDefs    string
}

func (r *mutationResolver) PushSchema(ctx context.Context, schemaInput model.SchemaInput) (bool, error) {
	if _, err := gqlparser.LoadSchema(&ast.Source{Name: schemaInput.ServiceName, Input: schemaInput.TypeDefs, BuiltIn: false}); err != nil {
		fmt.Println(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})

	key := schemaInput.ServiceName

	schemaRegistry := &schemaRegistry{
		ServiceName: schemaInput.ServiceName,
		ServiceURL:  schemaInput.ServiceURL,
		TypeDefs:    schemaInput.TypeDefs,
	}

	value, _ := json.Marshal(schemaRegistry)

	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		fmt.Println(err)
	}

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(val)
	return true, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) PlaceHolder(ctx context.Context) (*string, error) {
	text := "Hello world"
	return &text, nil
}
