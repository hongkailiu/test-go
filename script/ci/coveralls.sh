#!/bin/bash

if [[ $(echo "${TRAVIS}" | awk '{print tolower($0)}') != "true" ]]; then
  echo "please run this on travis-ci"
  exit 1
fi

go get -u github.com/mattn/goveralls

if [[ ! -d "build" ]]; then
  mkdir -v build
fi

#go test -coverprofile build/coverage.out $(go list ./... | grep -v github.com/hongkailiu/test-go/pkg/flexy)
go test -v -coverprofile build/coverage.out ./...
"${GOPATH}/bin/goveralls" -coverprofile=build/coverage.out -service travis-ci