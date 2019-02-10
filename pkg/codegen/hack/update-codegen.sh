#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail


SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/../../..

go get -d k8s.io/code-generator/cmd/client-gen
git -C ${GOPATH}/src/k8s.io/code-generator/ checkout kubernetes-1.12.1

CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${SCRIPT_ROOT}/; ls -d -1 ${GOPATH}/src/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

echo "SCRIPT_ROOT: ${SCRIPT_ROOT}"
echo "CODEGEN_PKG: ${CODEGEN_PKG}"

${CODEGEN_PKG}/generate-groups.sh all github.com/hongkailiu/test-go/pkg/codegen/pkg/client github.com/hongkailiu/test-go/pkg/codegen/pkg/apis "app:v1alpha1"
#  --go-header-file ${SCRIPT_ROOT}/hack/custom-boilerplate.go.txt
