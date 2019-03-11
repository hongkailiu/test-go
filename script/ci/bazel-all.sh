#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

if [[ "${TRAVIS:-false}" == "true" ]]; then
  echo "bazel on travis-ci ..."
  echo "bazel build ..."
  bazel build --jvmopt='-Xmx:2048m' --jvmopt='-Xms:2048m' //cmd/...
  echo "bazel test ..."
  ### ignore those package: seems bazel needs test file and target file are in the same pkg
  ### however, it is not the case for those 2 pkgs
  bazel test -- //... -//pkg/flexy/... -//pkg/ocptf/...
  exit 0
fi

if [[ "${CIRCLECI:-false}" == "true" ]]; then
  echo "bazel on circle-ci ..."
  echo "bazel build is skipped on circle-ci because of resource issue"
  #echo "bazel build ..."
  #bazel build --jobs=1 --jvmopt='-Xmx:2048m' --jvmopt='-Xms:2048m' //cmd/...
  echo "bazel test ..."
  bazel test -- //... -//pkg/flexy/... -//pkg/ocptf/...
  exit 0
fi

echo "not supported CI environment ... exit 1"
exit 1