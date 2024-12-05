# Load environment variables
include .envrc

# Define paths and database address
MIGRATIONS_PATH = ./cmd/migrate/migrations
# DB_ADDR = postgres://admin:adminpassword@localhost:5432/app?sslmode=disable
DOCKER_MIGRATE = docker run --rm -v "$(CURDIR)/cmd/migrate/migrations:/migrations" --network="host" migrate/migrate:latest

.PHONY: test
test:
	@go test -v ./...

.PHONY: migrate-create
migrate-create:
	@$(DOCKER_MIGRATE) create -seq -ext sql -dir /migrations $(name)

.PHONY: migrate-up
migrate-up:
	@$(DOCKER_MIGRATE) -path=/migrations -database "$(DB_ADDR)" up

.PHONY: migrate-down
migrate-down:
	@$(DOCKER_MIGRATE) -path=/migrations -database "$(DB_ADDR)" down -all

.PHONY: seed
seed: 
	@go run cmd/migrate/seed/main.go
	
.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt


