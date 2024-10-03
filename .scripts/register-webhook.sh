#!/usr/bin/env bash

set -o nounset
set -o errexit
set -o pipefail

WEBHOOK_URL="${INVOKE_URL}"

curl --request POST \
     --url "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/setWebhook?url=${WEBHOOK_URL}&secret_token=${TELEGRAM_BOT_SECRET}"
