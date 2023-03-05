#!/bin/bash
set -euf -o pipefail

SOURCE_DIR=$(readlink -f "$(dirname "${BASH_SOURCE[0]}")")
CONTAINER_NAME="htar-dev"

if [[ $( docker container inspect -f "{{.State.Running}}" "${CONTAINER_NAME}" 2> /dev/null ) == "true" ]]; then
  if [[ -z "$*" ]]; then
    set -- "/bin/sh"
  fi
  docker exec -it "${CONTAINER_NAME}" "$@"
  exit $?
fi

docker build --tag "localhost/htar:builder" --file "${SOURCE_DIR}/ci/builder.Dockerfile" "${SOURCE_DIR}/ci/"

docker run -it --rm \
  --user "$(id -u "${USER}"):$(id -g "${USER}")" \
  -e "GOCACHE=/tmp/go-cache" \
  -v "golang-cache:/go/pkg" \
  -v "${SOURCE_DIR}:/go/src/htar" \
  -w "/go/src/htar" \
  --name "${CONTAINER_NAME}" \
  localhost/htar:builder "$@"
