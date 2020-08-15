#!/bin/bash

set -eu

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/.."

docker run --rm \
  --volume=$PWD:/root \
  --workdir=/root \
  golang \
  ./scripts/build.local.sh "$@"
