# graphql-schema-registry
[![Go Reference](https://pkg.go.dev/badge/github.com/basselalaraaj/graphql-schema-registry.svg)](https://pkg.go.dev/github.com/basselalaraaj/graphql-schema-registry)

GraphQL Schema registry is used for discovery of services and keeps the schema of the gateway in sync with the service schemas it consumes, by automatically pushing schema changes.

When using a GraphQL gateway together with multiple GraphQL services, you will run into the issue that the gateway requires static paths to connect to multiple GraphQL services.
Using introspection the gateway reads the schemas of the GraphQL services and creates its own schema. But if introspection of one of the GraphQL services fails for example due to timeout, the gateway will fail generating the schema, resulting in a failing gateway. By validating the schemas and caching them in a schema registry, you only depend on one call to fetch all the required schemas to build the gateway schema. Lowering the risk of having an unavailable gateway.

The gateway reflects the schemas of multiple GraphQL services. When a GraphQL service releases a new schema, the gateway should be notified with the schema changes to avoid new queries that are not available due to an outdated gateway. By pushing the GraphQL schema to the registry, after storing the schema, the registry pushes the new GraphQL schema to the gateway. This way all new operation in the underlying GraphQL services, are made available through the gateway.

![registry](https://user-images.githubusercontent.com/5745279/103581587-98e75080-4edc-11eb-86c9-9d60329a2dc6.jpg)


## Getting started

Run the service through the command line

```bash
$ go run server.go
```

### Test

```bash
$ make test
```

### Lint

```bash
$ make lint
```

### Docker

```bash
$ export DOCKER_BUILDKIT=1
$ docker build -t graphql-schema-registry .
$ docker run -ti -p 8080:8080 graphql-schema-registry
```

### Tooling

For easily pushing and retrieving schemas from this graphql schema registry you can use our javascript package [graphql-schema-registry-tooling](https://github.com/basselalaraaj/graphql-schema-registry-tooling).

### Azure Service Bus events

Send a message to Azure service bus when a new schema is pushed to the registry.

Create a .env file for your configuration

```env
SERVICEBUS_ENABLED=True
SERVICEBUS_CONNECTION_STRING=
SERVICEBUS_TOPIC_NAME=
```
