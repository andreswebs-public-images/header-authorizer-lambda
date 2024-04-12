#!/usr/bin/env bash

set -o nounset
set -o errexit
set -o pipefail

MSG="${MSG:-Hello from ApiGateway!}"

curl \
    -v \
    --location \
    --header 'Content-Type: application/json' \
    --header "${HEADER_KEY}: ${TELEGRAM_BOT_SECRET}" \
    --data-raw "{ \"TestMessage\": \"${MSG}\" }" \
    --request POST \
    "${INVOKE_URL}/"
