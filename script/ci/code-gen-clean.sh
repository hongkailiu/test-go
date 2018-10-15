#!/bin/bash

set -e

rm -rfv codegen/pkg/client
rm -fv codegen/pkg/apis/app.example.com/v1alpha1/zz_generated.deepcopy.go