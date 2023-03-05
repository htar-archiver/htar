#!/bin/bash
set -euf -o pipefail

SOURCE_DIR=$(readlink -f "$(dirname "${BASH_SOURCE[0]}")")

"${SOURCE_DIR}/go-docker.sh" go build \
  -ldflags "-linkmode external -extldflags -static" \
  -o ./cmd/htar/htar \
  htar/cmd/htar
