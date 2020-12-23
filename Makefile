
.PHONY: all
all: test-all build

.PHONY: test-all
test-all: vet lint test

.PHONY: test
test:
	go test -v -parallel=4 ./...

.PHONY: lint
lint:
	docker run --rm \
		-v $(PWD):/app \
		-w /app \
		golangci/golangci-lint:v1.33.0 \
			golangci-lint --color never run

.PHONY: vet
vet:
	go vet ./...
