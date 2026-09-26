.PHONY: up down lint test

up:
	docker compose up -d

down:
	docker compose down

lint:
	golangci-lint run

test:
	go test -race ./...
