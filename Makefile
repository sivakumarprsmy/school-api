.PHONY: run build build-mcp build-agent test lint docker-up docker-down tidy

run:
	go run .

build:
	CGO_ENABLED=0 go build -o bin/school-api .

build-mcp:
	CGO_ENABLED=0 go build -o bin/school-mcp ./cmd/mcp

build-agent:
	CGO_ENABLED=0 go build -o bin/enrollment-agent ./cmd/agent

test:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -func=coverage.out

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

docker-up:
	docker compose up --build

docker-down:
	docker compose down -v
