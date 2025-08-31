SHELL := /bin/bash

DB_USER := $(shell sed -n 's/^DB_USER=//p' .env)
DB_PASSWORD := $(shell sed -n 's/^DB_PASSWORD=//p' .env)
DB_HOST := $(shell sed -n 's/^DB_HOST=//p' .env)
DB_PORT := $(shell sed -n 's/^DB_PORT=//p' .env)

install_tools:
	go install go.uber.org/mock/mockgen@latest

migration:
	docker run -v ./db/migrations/$(folder):/migrations migrate/migrate create -ext sql -dir /migrations -seq $(name)

migrate:
	docker run -v ./db/migrations/$(folder):/migrations --network host migrate/migrate -path=/migrations/ -database postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(folder)?sslmode=disable up $(version)

migrate-force:
	docker run -v ./db/migrations/$(folder):/migrations --network host migrate/migrate -path=/migrations/ -database postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(folder)?sslmode=disable force $(version)

generate:
	docker run --rm -v .:/src -w /src sqlc/sqlc generate
	mockgen -source=internal/repository/sequence.go -destination=internal/repository/mocks/sequence.go -package=mocks
	mockgen -source=internal/repository/step.go -destination=internal/repository/mocks/step.go -package=mocks

test:
	go test -v -cover ./internal/...

integ-test:
	go test -v ./integ-tests/...

run:
	go run cmd/api/main.go

build:
	@-mkdir build
	@go build -o build/sequence-technical-test cmd/api/main.go

run-bin:
	./build/sequence-technical-test

clear:
	-rm -r build