package graph

import (
	"context"
	"fmt"

	"github.com/basselalaraaj/graphql-schema-registry/graph/generated"
	"github.com/basselalaraaj/graphql-schema-registry/graph/model"
	"github.com/basselalaraaj/graphql-schema-registry/notify"
	"github.com/basselalaraaj/graphql-schema-registry/registry"
	"github.com/basselalaraaj/graphql-schema-registry/servicebus"
)

// Resolver It serves as dependency injection for your app, add any dependencies you require here.
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
	schemaRegistry := &registry.SchemaRegistry{
		ServiceName: schemaInput.ServiceName,
		ServiceURL:  schemaInput.ServiceURL,
		TypeDefs:    schemaInput.TypeDefs,
	}

	err := schemaRegistry.ValidateSchema()
	if err != nil {
		return false, err
	}

	err = schemaRegistry.Save()
	if err != nil {
		return false, err
	}

	go func() {
		serviceBus := servicebus.ServiceBus{}
		err := notify.SendNotification(schemaRegistry, &serviceBus)
		if err != nil {
			fmt.Println("send message to service bus failed")
		}
	}()

	return true, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) GetSchema(ctx context.Context, services []string) ([]*model.Schema, error) {
	servicesSchema := []*model.Schema{}
	for _, service := range services {
		schema, err := registry.GetServiceSchema(service)
		if err != nil {
			return nil, fmt.Errorf("not able to get schema for the service %v", service)
		}
		newSchema := model.Schema(*schema)
		servicesSchema = append(servicesSchema, &newSchema)
	}
	return servicesSchema, nil
}

func (r *queryResolver) GetAllSchemas(ctx context.Context) ([]*model.Schema, error) {
	servicesSchema := []*model.Schema{}
	services, err := registry.GetAllServices()
	if err != nil {
		return nil, fmt.Errorf("not able to get schema for the services")
	}
	for _, service := range *services {
		schema, err := registry.GetServiceSchema(service)
		if err != nil {
			fmt.Printf("not able to get schema for the service %v \n", service)
			continue
		}
		newSchema := model.Schema(*schema)
		servicesSchema = append(servicesSchema, &newSchema)
	}
	return servicesSchema, nil
}
