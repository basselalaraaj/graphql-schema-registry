
.PHONY: all
all: test-lint build

.PHONY: build
build:
	go build -o graphql-schema-registry .

.PHONY: test-lint
test-lint: lint test

.PHONY: test
test:
	go test -v -parallel=4 ./...

.PHONY: lint
lint:
	@golangci-lint run
