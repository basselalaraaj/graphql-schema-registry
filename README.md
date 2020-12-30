# graphql-schema-registry
[![Go Reference](https://pkg.go.dev/badge/github.com/basselalaraaj/graphql-schema-registry.svg)](https://pkg.go.dev/github.com/basselalaraaj/graphql-schema-registry)

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