NAME=jat
CURRENT_RELEASE=v0.2.0
GIT_COMMIT=$(shell git rev-parse --short HEAD)
BUILD_DATE=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
CC=go
FLAGS=-X github.com/filz0r/jat/internal/version.Version=$(CURRENT_RELEASE) -X github.com/filz0r/jat/internal/version.Commit=$(GIT_COMMIT) -X github.com/filz0r/jat/internal/version.BuildDate=$(BUILD_DATE)

COMPOSE_FILE=docker-compose.db.yml

all: $(NAME)

$(NAME): build_web
	@echo "Compiling new server binary"
	@echo "Version: $(CURRENT_RELEASE)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"
	@$(CC) build -ldflags "$(FLAGS)" -o $(NAME)
	@echo "Done!"

release:
	@if git rev-parse "$(CURRENT_RELEASE)" >/dev/null 2>&1; then \
              echo "Error: tag $(CURRENT_RELEASE) already exists"; \
              exit 1; \
	fi
	@echo "Creating release $(CURRENT_RELEASE)"
	@git tag -a "$(CURRENT_RELEASE)" -m "Release $(CURRENT_RELEASE)"
	@git push origin "$(CURRENT_RELEASE)"

clean_server:
	@rm -rf ./internal/api/webdist
	@echo "Removing old server binary"
	@rm -rf $(NAME)

build_server: $(NAME)

re: clean_server build_server

generate_api:
	@swag init -g main.go -o docs/api
	@rm -f docs/api/docs.go
	@cd web && npx swagger2openapi ../docs/api/swagger.json -o ../docs/api/openapi.json
	@cd web && npx openapi-typescript ../docs/api/openapi.json -o ./src/api/gen-spec.ts
	@cd web && npx prettier --write src/api/gen-spec.ts

build_web:
	@cd web && npm run build
	@rm -rf internal/api/webdist
	@cp -r web/dist internal/api/webdist

dev_setup_web:
	@cd web && npm install

dev_backend:
	@air

dev_frontend:
	@cd web && npm run dev

dev_db_down:
	@docker compose -f $(COMPOSE_FILE) down

dev_db_up:
	@docker compose -f $(COMPOSE_FILE) up -d

reset_db:
	@echo "dropping $(NAME) database and creating a new one with the same name"
	@docker exec jat sh -c 'psql -U postgres -d postgres -c "DROP DATABASE IF EXISTS $(NAME) WITH (FORCE);" && psql -U postgres -d postgres -c "CREATE DATABASE $(NAME);"'

connect_db:
	@docker exec -it $(NAME) psql -U postgres -d $(NAME)

demo:
	@echo "Building JAT"
	@make &> /dev/null
	@cp .env.sample .env
	@echo "launching database server using docker"
	@make dev_db_up &> /dev/null
	@echo "resetting database to blank state"
	@make reset_db &> /dev/null
	@echo "Seeding database with 'docs/assets/demo-db.sql'"
	@docker exec -i jat psql -U postgres -d jat < docs/assets/demo-db.sql &> /dev/null
	@echo "launching server on http://localhost:4200"
	@./jat server


.PHONY: all re test clean_server build_server generate_api dev_setup_web build_web dev_backend dev_frontend reset_db dev_db_up dev_db_down demo