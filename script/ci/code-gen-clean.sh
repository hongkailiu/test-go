#!/bin/bash

set -e

rm -rfv pkg/codegen/pkg/client
rm -fv pkg/codegen/pkg/apis/app.example.com/v1alpha1/zz_generated.deepcopy.go