#!/usr/bin/make -f

# ---------------------------------------------------------------------------- #
#                                 Make targets                                 #
# ---------------------------------------------------------------------------- #

.PHONY: clean
clean: ## Clean temporary files
	go clean

.PHONY: clean-cache
clean-cache: ## Clean go build cache
	go clean -cache

vet: ## Run go vet
	go vet ./cmd/rollappd

lint: ## Run linter
	golangci-lint run

# ---------------------------------------------------------------------------- #
#                                    testing                                   #
# ---------------------------------------------------------------------------- #
.PHONY: test
test: ## Run go test
	go test ./... 

.PHONY: test_evm
test_evm: ## Run go test
	go test ./... --tags=evm

# ---------------------------------------------------------------------------- #
#                                   Protobuf                                   #
# ---------------------------------------------------------------------------- #

DOCKER := $(shell which docker)
protoVer=0.14.0
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-format:
	@$(protoImage) find ./proto -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint --error-format=json

proto-gen:
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./proto/protocgen.sh

proto-check-bc-breaking:
	@rm -rf ./.proto-export
	@$(protoImage) buf breaking --against .git#branch=$$(git merge-base HEAD origin/main)

proto-export:
	@rm -rf ./.proto-export && cd proto && buf export --config ./buf.yaml --output ../.proto-export

proto-export-deps:
	@rm -rf ./.proto-export-deps
	@cd proto && buf export --config ./buf.yaml --output ../.proto-export-deps --exclude-imports && buf export buf.build/cosmos/cosmos-sdk:v0.50.0 --output ../.proto-export-deps

PROTO_DIRS=$(shell find .proto-export-deps -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)

.PHONY: proto-format proto-lint proto-check-bc-breaking proto-export proto-export-deps


# ---------------------------------------------------------------------------- #
#                                   Misc                                       #
# ---------------------------------------------------------------------------- #

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## generates help for all targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
