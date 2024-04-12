#!/usr/bin/env bash
set -o nounset
cat <<EOT
{
  "type": "REQUEST",
  "methodArn": "arn:aws:execute-api:us-east-1:123456789012:abcdef123/test/GET/request",
  "headers": {
    "X-Telegram-Bot-Api-Secret-Token": "${TELEGRAM_BOT_SECRET}"
  }
}
EOT
