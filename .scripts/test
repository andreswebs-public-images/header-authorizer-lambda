#!/usr/bin/env bash

set -o nounset
set -o errexit
set -o pipefail

curl \
    --request POST \
    --header "Content-Type: application/json" \
    --data "@${EVENT}" \
    "http://localhost:9000/2015-03-31/functions/function/invocations"
