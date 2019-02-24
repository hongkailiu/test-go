#!/bin/bash

if [[ $(echo "${CI}" | awk '{print tolower($0)}') != "true" ]]; then
  echo "please run this on ci system, like travis ci or circle ci"
  exit 1
fi

ci_service="unknown"

if [[ $(echo "${TRAVIS}" | awk '{print tolower($0)}') == "true" ]]; then
  ci_service="travis-ci"
fi

if [[ $(echo "${CIRCLECI}" | awk '{print tolower($0)}') == "true" ]]; then
  ci_service="circle-ci"
fi

if [[ ${ci_service} == "unknown" ]]; then
  echo "only travis ci and circle ci are supported for now"
  exit 1
fi

go get -u github.com/mattn/goveralls

"${GOPATH}/bin/goveralls" -coverprofile=build/coverage.out -service "${ci_service}"
