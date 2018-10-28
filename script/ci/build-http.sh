#!/bin/bash


go build -o build/http ./cmd/
cp -rv http/static build/