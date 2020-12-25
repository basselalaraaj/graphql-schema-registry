# graphql-schema-registry
[![Go Reference](https://pkg.go.dev/badge/github.com/basselalaraaj/graphql-schema-registry.svg)](https://pkg.go.dev/github.com/basselalaraaj/graphql-schema-registry)

## Getting started

Create a .env file for your configuration

```env
SERVICEBUS_CONNECTION_STRING=
SERVICEBUS_TOPIC_NAME=
```

Run the service through the command line

```bash
$ go run server.go
```

### Test

The `test-all` target run lint using (golangci-lint)[https://golangci-lint.run]
and unit tests:

```bash
$ make test-all
```

### Docker

```bash
$ export DOCKER_BUILDKIT=1
$ docker build -t graphql-schema-registry .
$ docker run -ti -p 8080:8080 graphql-schema-registry
```

## Development

Use air command line for live reloading for development
```bash
$ air
```