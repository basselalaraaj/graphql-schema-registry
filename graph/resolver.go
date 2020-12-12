package graph

import (
	"context"
	"fmt"

	"github.com/basselalaraaj/graphql-schema-registry/graph/generated"
	"github.com/basselalaraaj/graphql-schema-registry/graph/model"
)

// This file will not be regenerated automatically.
//
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

func (r *mutationResolver) PushSchema(ctx context.Context, schemaInput model.SchemaInput) (bool, error) {
	fmt.Println(schemaInput)
	return true, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) PlaceHolder(ctx context.Context) (*string, error) {
	text := "Hello world"
	return &text, nil
}
