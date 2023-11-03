RUN ?= .*
PKG ?= ./...
.PHONY: test
test: ## Run tests in local environment
	golangci-lint run --timeout=5m $(PKG)
	go test -cover -benchmem -bench=. -run=$(RUN) $(PKG)

.PHONY: license-check
license-check:
	licensed cache
	licensed status

.PHONY: docker-license-check
docker-license-check:
	@docker run --workdir /app --entrypoint make -v $(shell pwd):/app public.ecr.aws/kanopy/licensed-go license-check

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
