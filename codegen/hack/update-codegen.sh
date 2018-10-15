#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/../..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${SCRIPT_ROOT}/; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

echo "SCRIPT_ROOT: ${SCRIPT_ROOT}"
echo "CODEGEN_PKG: ${CODEGEN_PKG}"

${CODEGEN_PKG}/generate-groups.sh all github.com/hongkailiu/test-go/codegen/pkg/client github.com/hongkailiu/test-go/codegen/pkg/apis "app.example.com:v1alpha1"
#  --go-header-file ${SCRIPT_ROOT}/hack/custom-boilerplate.go.txt

