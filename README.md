# graphql-schema-registry

## Getting started

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

Use air command line for live reloading
```bash
$ air
```