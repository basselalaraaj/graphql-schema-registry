package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/basselalaraaj/graphql-schema-registry/schema"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func main() {
	schema, err := graphql.NewSchema(schema.SchemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	h := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	})

	http.Handle("/graphql", h)
	fmt.Println("Now server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
