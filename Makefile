.PHONY: help test test-release
.DEFAULT_GOAL:=help

test: ## Run tests
	go mod tidy
	go test -v ./...

test-release: test ## Test release procedure with goreleaser
	goreleaser --snapshot --skip-publish --rm-dist

clean: ## Clean up repo
	rm -rf vendor dist

help: ## Show help
	@echo "Makefile help:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
