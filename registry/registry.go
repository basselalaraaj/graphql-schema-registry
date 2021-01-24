package registry

import (
	"fmt"

	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

// SchemaRegistry to add a schema to the registry
type SchemaRegistry struct {
	ServiceName string
	ServiceURL  string
	TypeDefs    string
}

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
	err := MongoDb.saveSchema(s)
	if err != nil {
		return err
	}

	err = s.setSchema()
	if err != nil {
		return err
	}
	return nil
}

// GetServiceSchema get service schema from redis
func GetServiceSchema(service string) (*SchemaRegistry, error) {
	result, err := getSchema(service)

	if err != nil {
		return &SchemaRegistry{}, err
	}

	return result, nil
}

// GetAllServices returns all services names
func GetAllServices() (*[]string, error) {
	serviceSchemas, err := scanSchemas()
	if err != nil {
		return &[]string{}, err
	}

	if len(*serviceSchemas) == 0 {
		var results []string

		if err := MongoDb.getServiceSchemas(&results); err != nil {
			return nil, err
		}

		return &results, nil
	}
	return serviceSchemas, nil
}
