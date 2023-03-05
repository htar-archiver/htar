#!/bin/bash
set -euo pipefail

SOURCE_DIR=$(readlink -f "$(dirname "${BASH_SOURCE[0]}")")
CONTAINER_NAME="htar-dev"

if [[ $( docker container inspect -f "{{.State.Running}}" "${CONTAINER_NAME}" 2> /dev/null ) == "true" ]]; then
  docker exec -it "${CONTAINER_NAME}" "$@"
  exit $?
fi

docker run -it --rm \
  --user "$(id -u "${USER}"):$(id -g "${USER}")" \
  -e "GOCACHE=/tmp/go-cache" \
  -v "golang-cache:/go/pkg" \
  -v "${SOURCE_DIR}:/go/src/htar" \
  -w "/go/src/htar" \
  --name "${CONTAINER_NAME}" \
  golang:1.16-alpine "$@"
