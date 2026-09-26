.PHONY: up down test

up:
	docker compose up -d

down:
	docker compose down

test:
	go test -race ./...
