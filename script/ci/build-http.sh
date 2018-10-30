#!/bin/bash

set -e

go build -o build/http ./cmd/
cp -rv http/static build/
cp -rv http/swagger build/