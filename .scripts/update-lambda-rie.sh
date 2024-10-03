#!/usr/bin/env bash

set -o nounset
set -o errexit
set -o pipefail


LAMBDA_RIE_URL="https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie"

[ "${ARCH}" = "arm64" ] && {
  LAMBDA_RIE_URL="${LAMBDA_RIE_URL}-arm64"
}

echo "${LAMBDA_RIE_URL}"

mkdir -p ~/.aws-lambda-rie
pushd ~/.aws-lambda-rie
curl -Lo aws-lambda-rie "${LAMBDA_RIE_URL}"
chmod +x aws-lambda-rie
popd
