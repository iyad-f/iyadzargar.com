# SPDX-FileCopyrightText: 2026 Iyad
# SPDX-License-Identifier: Apache-2.0

.DEFAULT_GOAL := help

.PHONY: help setup dev generate css test lint fmt tidy

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "} {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'

setup: ## Install git hooks and download dependencies
	pre-commit install
	go mod download

dev: ## Run the dev server with live reload (templ, CSS, and Go via air)
	go tool air

generate: ## Generate templ components
	go tool templ generate

css: ## Build the Tailwind CSS
	tailwindcss -i web/static/css/input.css -o web/static/css/styles.css --minify

test: ## Run tests with the race detector
	go test -race ./...

lint: ## Run golangci-lint
	golangci-lint run

fmt: ## Format templ and Go files
	go tool templ fmt .
	golangci-lint fmt

tidy: ## Tidy module dependencies
	go mod tidy
