#!/bin/bash

set -e

go build -o build/http ./cmd/http/
cp -rv pkg/http/static build/
cp -rv pkg/http/swagger build/