#!/usr/bin/env bash

set -eo pipefail

proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
	for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
		if grep "option go_package" $file &>/dev/null; then
			buf generate --template ./proto/buf.gen.gogo.yaml $file
		fi
	done
done

# move proto files to the right places
cp -r .gen/github.com/dymensionxyz/dymension-rdk/* ./proto
rm -rf .gen/

# move generated .go files from proto to their respective directories
for file in $(find ./proto -name '*.go'); do
    dir=$(dirname "${file#./proto/}")
    mkdir -p "./${dir}"
    mv "$file" "./${dir}"
done

rm -rf proto/x
