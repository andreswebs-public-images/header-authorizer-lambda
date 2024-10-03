#!/usr/bin/env bash

set -o nounset
set -o errexit
set -o pipefail

docker build --platform "${PLATFORM}" --load --tag "${IMG}" .
