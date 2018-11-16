#!/bin/bash

if [[ $(echo "${TRAVIS}" | awk '{print tolower($0)}') != "true" ]]; then
  echo "please run this on travis-ci"
  exit 1
fi

go get -u github.com/mattn/goveralls
mkdir -v build
go test -coverprofile build/coverage.out $(go list ./... | grep -v github.com/hongkailiu/test-go/pkg/flexy)
"${GOPATH}/bin/goveralls" -coverprofile=build/coverage.out -service travis-ci