#!/usr/bin/env bash

set -euo pipefail

command="cincinnati"
git_commit="$( git describe --tags --always --dirty )"
build_date="$( date -u '+%Y%m%d' )"
version="v${build_date}-${git_commit}"

eval $(go env | grep -e "GOHOSTOS" -e "GOHOSTARCH")
GOOS=${GOOS:-${GOHOSTOS}}

OUT_DIR=${OUT_DIR-_out}
mkdir -p ${OUT_DIR}

set -x

CGO_ENABLED=0 GOOS="${GOOS}" go build -ldflags "-X 'github.com/hongkailiu/test-go/pkg/version.Name=${command}' -X 'github.com/hongkailiu/test-go/pkg/version.Version=${version}'" -a -installsuffix cgo -o "${OUT_DIR}/${command}" "./cmd/${command}"
