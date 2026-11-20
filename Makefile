-include .env
export

APP_NAME := lorebound
BIN_DIR  := bin
DB_URL   ?= postgres://$(DB_USER):$(DB_PASSWORD)@postgres:5432/$(DB_NAME)?sslmode=disable

.DEFAULT_GOAL := help

.PHONY: help build run test lint fmt vet tidy \
	compose-build compose-up compose-debug compose-stop compose-reset \
	compose-build-no-cache compose-logs compose-ps \
	migrate-up migrate-down migrate-reset

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-24s %s\n", $$1, $$2}'

## Go

build: ## Build the binary into bin/
	go build -o $(BIN_DIR)/$(APP_NAME) ./cmd

run: ## Run the app locally (needs compose services up)
	go run ./cmd

test: ## Run tests with the race detector
	go test -race -count=1 ./...

lint: ## Run golangci-lint
	golangci-lint run

fmt: ## Format code
	gofmt -w .

vet: ## Run go vet
	go vet ./...

tidy: ## Tidy go.mod and go.sum
	go mod tidy

## Docker Compose

compose-build: ## Build Docker images
	docker compose build

compose-build-no-cache: ## Build images without cache
	docker compose build --no-cache

compose-up: ## Build and start all containers in the background
	docker compose up -d --build

compose-debug: compose-build ## Start in the foreground with logs
	docker compose up

compose-stop: ## Stop and remove containers (keeps data volumes)
	docker compose down

compose-reset: ## Stop containers AND delete all data volumes
	docker compose down -v --remove-orphans

compose-logs: ## Follow app logs
	docker compose logs -f app

compose-ps: ## Show container status and health
	docker compose ps

## Migrations (run through the compose migrate service)

migrate-up: ## Apply all migrations
	docker compose run --rm migrate -path=/migrations -database="$(DB_URL)" up

migrate-down: ## Roll back the last migration
	docker compose run --rm migrate -path=/migrations -database="$(DB_URL)" down 1

migrate-reset: ## Roll back all migrations
	docker compose run --rm migrate -path=/migrations -database="$(DB_URL)" down -all