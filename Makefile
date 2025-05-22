.DEFAULT_GOAL := build
.PHONY: fmt vet build

export

MIGRATIONS_PATH=./migrations
DSN=postgres://postgres:postgres@127.0.0.1:54323/idm?sslmode=disable

test:
	go test ./...

fmt: test
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o ./bin/idm ./cmd

up:
	docker compose -f ./docker/docker-compose.yml up -d

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