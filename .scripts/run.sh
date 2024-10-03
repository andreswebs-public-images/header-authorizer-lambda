#!/usr/bin/env bash

set -o nounset
set -o errexit
set -o pipefail

eval "$(aws-vault export "${AWS_PROFILE}")"

NAME="local-lambda"

docker run \
    --name "${NAME}" \
    --platform "${PLATFORM}" \
    --rm \
    --detach \
    --volume ~/.aws-lambda-rie:/aws-lambda \
    --publish 9000:8080 \
    --entrypoint /aws-lambda/aws-lambda-rie \
    --env AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID}" \
    --env AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY}" \
    --env AWS_SESSION_TOKEN="${AWS_SESSION_TOKEN}" \
    --env AWS_REGION="${AWS_REGION}" \
    --env AWS_DEFAULT_REGION="${AWS_DEFAULT_REGION}" \
    --env HEADER_KEY \
    --env HEADER_VALUE_PARAMETER \
    --env PRINCIPAL_ID \
    "${IMG}" /bootstrap
