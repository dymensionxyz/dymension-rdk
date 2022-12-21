#!/usr/bin/env bash

set -eo pipefail

# Generate the `types` proto files
buf generate --path="./proto/sequencers" --template="buf.gen.yaml" --config="buf.yaml"