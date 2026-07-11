# ==============================================================================
# gobooking - Development Makefile
# ==============================================================================

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.DEFAULT_GOAL := help

COMPOSE := docker compose -f docker-compose.yml

APP_CONTAINER := gobooking_app

DOMAINS := api.gobooking.local mail.gobooking.local db.gobooking.local

GREEN  := \033[0;32m
YELLOW := \033[1;33m
BLUE   := \033[0;34m
RED    := \033[0;31m
RESET  := \033[0m

CERT_FILE := certs/gobooking.local.pem
CERT_KEY  := certs/gobooking.local-key.pem

.PHONY: help run build \
        fmt vet test check \
        tidy update \
        clean doctor \
        hosts certs up down restart logs ps bash \
        migrate-up migrate-down sqlc

help: ## Show available commands
	@echo ""
	@echo "$(BLUE)gobooking Development Commands$(RESET)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_-]+:.*##/ {printf "  \033[32m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""

# ==============================================================================
# Development
# ==============================================================================

run: ## Run the server
	cd app && go run ./cmd/api

build: ## Build the local binary
	@mkdir -p app/bin
	cd app && go build -o bin/booking ./cmd/api
	@echo "$(GREEN)✓ Binary generated in app/bin/booking$(RESET)"

# ==============================================================================
# Quality
# ==============================================================================

fmt: ## Format the source code
	cd app && go fmt ./...

vet: ## Run go vet
	cd app && go vet ./...

test: ## Run unit tests
	cd app && go test ./...

check: fmt vet test ## Run all quality checks

# ==============================================================================
# Dependencies
# ==============================================================================

tidy: ## Clean up go.mod / go.sum
	cd app && go mod tidy

update: ## Update dependencies
	cd app && go get -u ./...
	cd app && go mod tidy

# ==============================================================================
# Database
# ==============================================================================
migrations: ## Create migrations | make migrations t="table_name"
	@if [ $(t) ]; then \
  		echo "$(GREEN)Migrations building ... $(RESET)"; \
		cd app && migrate create -ext sql -dir db/migrations -seq ${t}; \
		echo "$(GREEN)Migrations built $(RESET)"; \
	else \
		echo "$(RED)(t) param is required (make migrations t='table_name') $(RESET)"; \
	fi

migrate-up: ## Apply migrations
	@echo "$(GREEN)Database migrations up ... $(RESET)";
	migrate -path app/db/migrations -database "$$DATABASE_URL" up
	@echo "$(GREEN)Database migrations finished! $(RESET)";

migrate-down: ## Roll back the last migration
	@echo "$(GREEN)Rollback last database migration ... $(RESET)";
	migrate -path app/db/migrations -database "$$DATABASE_URL" down 1;
	@echo "$(GREEN)Rollback done! $(RESET)";

sqlc: ## Regenerate Go code from SQL queries
	cd app && sqlc generate

# ==============================================================================
# Docker
# ==============================================================================

hosts: ## Add local domains to /etc/hosts (requires sudo)
	@echo "$(YELLOW)Updating /etc/hosts...$(RESET)"
	@for domain in $(DOMAINS); do \
		if grep -qE "^127\.0\.0\.1[[:space:]]+$$domain$$" /etc/hosts; then \
			echo "$(GREEN)$$domain already present$(RESET)"; \
		else \
			echo "127.0.0.1 $$domain" | sudo tee -a /etc/hosts > /dev/null; \
			echo "$(GREEN)$$domain added$(RESET)"; \
		fi; \
	done

certs: hosts ## Generate local TLS certificates if missing (requires mkcert)
	@if [ -f $(CERT_FILE) ] && [ -f $(CERT_KEY) ]; then \
		echo "$(GREEN)Certificates already present$(RESET)"; \
	else \
		echo "$(YELLOW)Generating certificates...$(RESET)"; \
		mkcert -install; \
		mkcert -cert-file $(CERT_FILE) -key-file $(CERT_KEY) gobooking.local "*.gobooking.local"; \
		echo "$(GREEN)Certificates generated in certs/$(RESET)"; \
	fi

up: certs ## Build and start the containers
	@echo "$(YELLOW)Starting containers...$(RESET)"
	$(COMPOSE) up -d --build
	@echo "$(GREEN)Containers started$(RESET)"
	@echo "$(BLUE)Traefik dashboard: http://localhost:8080$(RESET)"
	@echo "$(BLUE)API URL: https://api.gobooking.local$(RESET)"
	@echo "$(YELLOW)  -> root path returns 404, call a route declared in the router (e.g. /health)$(RESET)"
	@echo "$(BLUE)Mailpit URL: https://mail.gobooking.local$(RESET)"
	@echo "$(BLUE)Adminer URL: https://db.gobooking.local$(RESET)"

down: ## Stop the containers
	@echo "$(YELLOW)Stopping containers...$(RESET)"
	$(COMPOSE) down
	@echo "$(GREEN)Containers stopped$(RESET)"

restart: down up ## Restart the containers

logs: ## Show container logs
	@echo "$(YELLOW)Showing logs...$(RESET)"
	$(COMPOSE) logs -f

ps: ## List containers
	@echo "$(YELLOW)Listing containers...$(RESET)"
	$(COMPOSE) ps

bash: ## Access the app container
	@echo "$(YELLOW)Accessing the app container...$(RESET)"
	docker exec -it $(APP_CONTAINER) sh

# ==============================================================================
# Utilities
# ==============================================================================

clean: ## Remove generated files
	rm -rf app/bin

doctor: ## Show the development environment
	@echo ""
	@echo "$(BLUE)Environment$(RESET)"
	@echo ""
	@go version
	@git --version
