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

containerProtoVer=v0.2
containerProtoImage=tendermintdev/sdk-proto-gen:$(containerProtoVer)
containerProtoGen=cosmos-sdk-proto-gen-$(containerProtoVer)
containerProtoFmt=cosmos-sdk-proto-fmt-$(containerProtoVer)

proto-gen: ## Generates protobuf files
	@echo "Generating Protobuf files"
	@which protoc >/dev/null || (echo "protoc not found. Installing..." && \
        (uname | grep -q Darwin && brew install protobuf || sudo apt install -y protobuf-compiler))
	@sh ./scripts/protocgen.sh

proto-format: ## Formats protobuf files
	@echo "Formatting Protobuf files"
	@which protoc >/dev/null || (echo "protoc not found. Installing..." && \
        (uname | grep -q Darwin && brew install protobuf || sudo apt install -y protobuf-compiler))
	@find ./ -not -path "./third_party/*" -name "*.proto" -exec clang-format -i {} \;


# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## generates help for all targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
