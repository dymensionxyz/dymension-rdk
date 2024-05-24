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

###############################################################################
###                                Proto                                    ###
###############################################################################

# ------
# NOTE: Link to the tendermintdev/sdk-proto-gen docker images:
#       https://hub.docker.com/r/tendermintdev/sdk-proto-gen/tags
#
protoVer=v0.7
protoImageName=tendermintdev/sdk-proto-gen:$(protoVer)
containerProtoGen=cosmos-sdk-proto-gen-$(protoVer)
containerProtoFmt=cosmos-sdk-proto-fmt-$(protoVer)
# ------
# NOTE: cosmos/proto-builder image is needed because clang-format is not installed
#       on the tendermintdev/sdk-proto-gen docker image.
#		Link to the cosmos/proto-builder docker images:
#       https://github.com/cosmos/cosmos-sdk/pkgs/container/proto-builder
#
protoCosmosVer=0.11.2
protoCosmosName=ghcr.io/cosmos/proto-builder:$(protoCosmosVer)
protoCosmosImage=$(DOCKER) run --network host --rm -v $(CURDIR):/workspace --workdir /workspace $(protoCosmosName)

proto-gen:
	@echo "Generating Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGen}$$"; then docker start -a $(containerProtoGen); else docker run --name $(containerProtoGen) -v $(CURDIR):/workspace --workdir /workspace $(protoImageName) \
		sh ./scripts/protocgen.sh; fi
	@go mod tidy

proto-swagger-gen:
	@echo "Downloading Protobuf dependencies"
	@make proto-download-deps
	@echo "Generating Protobuf Swagger"
	$(protoCosmosImage) sh ./scripts/protoc-swagger-gen.sh

SWAGGER_DIR=./swagger-proto
THIRD_PARTY_DIR=$(SWAGGER_DIR)/third_party

proto-download-deps:
	mkdir -p "$(THIRD_PARTY_DIR)/cosmos_tmp" && \
	cd "$(THIRD_PARTY_DIR)/cosmos_tmp" && \
	git init && \
	git remote add origin "https://github.com/cosmos/cosmos-sdk.git" && \
	git config core.sparseCheckout true && \
	printf "proto\nthird_party\n" > .git/info/sparse-checkout && \
	git fetch --depth=1 origin "$(DEPS_COSMOS_SDK_VERSION)" && \
	git checkout FETCH_HEAD && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/cosmos_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/ethermint_tmp" && \
	cd "$(THIRD_PARTY_DIR)/ethermint_tmp" && \
	git init && \
	git remote add origin "https://github.com/dymensionxyz/ethermint.git" && \
	git config core.sparseCheckout true && \
	printf "proto\nthird_party\n" > .git/info/sparse-checkout && \
	git fetch --depth=1 origin "$(DEPS_ETHERMINT_VERSION)" && \
	git checkout FETCH_HEAD && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/ethermint_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/osmosis_tmp" && \
	cd "$(THIRD_PARTY_DIR)/osmosis_tmp" && \
	git init && \
	git remote add origin "https://github.com/dymensionxyz/osmosis.git" && \
	git config core.sparseCheckout true && \
	printf "proto\nthird_party\n" > .git/info/sparse-checkout && \
	git fetch --depth=1 origin "$(DEPS_OSMOSIS_VERSION)" && \
	git checkout FETCH_HEAD && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/osmosis_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/ibc_tmp" && \
	cd "$(THIRD_PARTY_DIR)/ibc_tmp" && \
	git init && \
	git remote add origin "https://github.com/cosmos/ibc-go.git" && \
	git config core.sparseCheckout true && \
	printf "proto\n" > .git/info/sparse-checkout && \
	git fetch --depth=1 origin "$(DEPS_IBC_GO_VERSION)" && \
	git checkout FETCH_HEAD && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/ibc_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/cosmos_proto_tmp" && \
	cd "$(THIRD_PARTY_DIR)/cosmos_proto_tmp" && \
	git init && \
	git remote add origin "https://github.com/cosmos/cosmos-proto.git" && \
	git config core.sparseCheckout true && \
	printf "proto\n" > .git/info/sparse-checkout && \
	git fetch --depth=1 origin "$(DEPS_COSMOS_PROTO_VERSION)" && \
	git checkout FETCH_HEAD && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/cosmos_proto_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/gogoproto" && \
	curl -SSL https://raw.githubusercontent.com/cosmos/gogoproto/$(DEPS_COSMOS_GOGOPROTO_VERSION)/gogoproto/gogo.proto > "$(THIRD_PARTY_DIR)/gogoproto/gogo.proto"

	mkdir -p "$(THIRD_PARTY_DIR)/google/api" && \
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto > "$(THIRD_PARTY_DIR)/google/api/annotations.proto"
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto > "$(THIRD_PARTY_DIR)/google/api/http.proto"

	mkdir -p "$(THIRD_PARTY_DIR)/confio/ics23" && \
	curl -sSL https://raw.githubusercontent.com/confio/ics23/$(DEPS_CONFIO_ICS23_VERSION)/proofs.proto > "$(THIRD_PARTY_DIR)/proofs.proto"


# ---------------------------------------------------------------------------- #
#                                   Misc                                       #
# ---------------------------------------------------------------------------- #

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## generates help for all targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
