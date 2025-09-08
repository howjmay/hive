#!/bin/bash

version=latest
if [ -n "$1" ]; then
    version="$1"
fi

sims=$(find simulators -name go.mod)
for d in $sims; do
    d="$(dirname $d)"
    echo "upgrading hivesim in $d"
    (
        set -e
        cd $d
        go get -d "github.com/iotaledger/wasp-hive@$version"
        go mod tidy -compat=1.21
    )
done
