#!/usr/bin/env bash

set -o nounset
set -o errexit
set -o pipefail

MSG=$(aws sqs receive-message --queue-url "${QUEUE_URL}")
MSG_RECEIPT=$([ -n "${MSG}"  ] && echo -E "${MSG}" | jq -r '.Messages[] | .ReceiptHandle')
aws sqs delete-message --queue-url "${QUEUE_URL}" --receipt-handle "${MSG_RECEIPT}"
echo -E "${MSG}" | jq -r '.Messages[] | .Body'
