#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

CODEGEN_PKG=${GOPATH}/src/k8s.io/code-generator

echo "CODEGEN_PKG: ${CODEGEN_PKG}"

${CODEGEN_PKG}/generate-groups.sh all github.com/hongkailiu/test-go/codegen/pkg/client github.com/hongkailiu/test-go/codegen/pkg/apis "app:v1alpha1"
#  --go-header-file ${SCRIPT_ROOT}/hack/custom-boilerplate.go.txt

