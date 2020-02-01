LDFLAGS := -ldflags="-s -w"

.PHONY: help test lint mod-check

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test: ## Run tests
	go test -race -v ./... -count=1

lint: ## Run GolangCI-Lint
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run -v

mod-check: ## Run check on go mod tidy
	go mod tidy && git --no-pager diff --exit-code -- go.mod go.sum

tfvar: ## Build tfvar binary
	go build $(LDFLAGS) -o tfvar

install: ## Install tfvar into $GOBIN
	go install $(LDFLAGS)

clean: ## Cleanup
	rm tfvar
