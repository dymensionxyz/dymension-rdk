#!/bin/sh

set -eo pipefail

# get protoc executions
# go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos 2>/dev/null

proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    if grep "option go_package" $file &> /dev/null ; then
      buf generate --template ./proto/buf.gen.gogo.yaml $file
    fi
  done
done

# move proto files to the right places
cp -r github.com/dymensionxyz/dymension-rdk/* ./
rm -rf github.com