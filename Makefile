# Makefile for Local RAG API

.PHONY: docs build run test clean

# Generate Swagger documentation
docs:
	$(HOME)/go/bin/swag init -g cmd/server/main.go -o docs
	@echo "Swagger documentation updated"

build: docs
	go build -o bin/server ./cmd/server

run: docs
	go run ./cmd/server

test:
	go test ./...

clean:
	rm -rf bin/
	rm -f docs/docs.go docs/swagger.json docs/swagger.yaml

deps:
	go mod tidy
	go mod download

fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# development workflow
dev: fmt docs run
