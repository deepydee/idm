.DEFAULT_GOAL := build
.PHONY: fmt vet build

include .env
export

MIGRATIONS_PATH=./migrations
DSN=$(DB_CONNECTION)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

test-integration:
	@export $(shell cat tests/.env | xargs); \
	DSN=$${DB_CONNECTION}://$${DB_USER}:$${DB_PASSWORD}@$${DB_HOST}:$${DB_PORT}/$${DB_NAME}?sslmode=disable; \
	migrate -path=$(MIGRATIONS_PATH) -database="$$DSN" up; \
	DSN="$$DSN" go test -v ./tests/...

test-inner:
	go test ./inner/...

fmt: test-inner
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o ./bin/idm ./cmd

up:
	docker compose -f ./docker/docker-compose.yml up -d --build

down:
	docker compose -f ./docker/docker-compose.yml down

migrate-up:
	migrate \
		-path=$(MIGRATIONS_PATH) \
		-database="$(DSN)" up

migrate-down:
	migrate \
		-path=$(MIGRATIONS_PATH) \
		-database="$(DSN)" down

migrate-version:
	migrate \
		-path=$(MIGRATIONS_PATH) \
		-database="$(DSN)" version

migrate-goto:
	@read -p "Enter target version: " version; \
	migrate \
	  -path=./migrations \
	  -database="postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" \
	  goto $$version

migrate-force:
	@read -p "Enter target version: " version; \
	migrate \
		-path=$(MIGRATIONS_PATH) \
		-database="$(DSN)" force $$version