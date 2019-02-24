#!/bin/bash

if [[ $(echo "${CI}" | awk '{print tolower($0)}') != "true" ]]; then
  echo "please run this on ci system, like travis ci or circle ci"
  exit 1
fi

go get -u github.com/mattn/goveralls

if [[ $(echo "${TRAVIS}" | awk '{print tolower($0)}') == "true" ]]; then
  "${GOPATH}/bin/goveralls" -coverprofile=build/coverage.out -service=travis-ci
fi

if [[ $(echo "${CIRCLECI}" | awk '{print tolower($0)}') == "true" ]]; then
  ##https://github.com/lemurheavy/coveralls-public/issues/632
  "${GOPATH}/bin/goveralls" -coverprofile=build/coverage.out -service=circle-ci -repotoken="${COVERALLS_TOKEN}"
fi


