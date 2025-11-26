all: help

help: ## Print this help message
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: hello
hello: ## Print "Hello, World!"
	echo "Hello, World!"

.PHONY: run
run: ## Run the main application with GitHub token
	GITHUB_TOKEN=$$(gh auth token) go run main.go

.PHONY: test
test: ## Run all tests with GitHub token
	GITHUB_TOKEN=$$(gh auth token) go test ./...
