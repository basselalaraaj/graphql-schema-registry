// GraphQL Schema registry is used for discovery of services and keeps the schema of the gateway in sync
// with the service schemas it consumes, by automatically pushing schema changes.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/basselalaraaj/graphql-schema-registry/graph"
	"github.com/basselalaraaj/graphql-schema-registry/graph/generated"
	"github.com/basselalaraaj/graphql-schema-registry/registry"
	"github.com/basselalaraaj/graphql-schema-registry/servicebus"
	"github.com/joho/godotenv"
)

const defaultPort = "8080"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	servicebus.Initialize()
	registry.InitializeDatabase()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	srv.Use(extension.Introspection{})
	http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	http.Handle("/graphql", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
