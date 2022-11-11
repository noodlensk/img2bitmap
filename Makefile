setup-lint: ## Set up linter
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.50.1
setup: setup-lint ## Setup tooling
fmt: ## gofmt and goimports all go files
	find . -name '*.go' | while read -r file; do gofumpt -w "$$file"; goimports -w "$$file"; done
lint: ## Lint
	golangci-lint run
lint-fix: ## Autofix lint errors
	golangci-lint run --fix
test: ## Run tests
	go test -count=1 -p=8 -parallel=8 -race ./...
dep: ## Get all dependencies
	go mod download && go mod tidy
install : ## Install app
	go install .
build: ## Build all projects
	go build
# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
