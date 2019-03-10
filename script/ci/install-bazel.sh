#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

if [[ "${TRAVIS:-false}" == "true" ]]; then
  echo "installing bazel on travis-ci ..."
  echo "deb [arch=amd64] http://storage.googleapis.com/bazel-apt stable jdk1.8" | sudo tee /etc/apt/sources.list.d/bazel.list
  curl https://bazel.build/bazel-release.pub.gpg | sudo apt-key add -
  sudo apt-get update && sudo apt-get install bazel
  exit 0
fi

if [[ "${CIRCLECI:-false}" == "true" ]]; then
  echo "installing openjdk on circle-ci ..."
  sudo apt-get install -y software-properties-common
  sudo add-apt-repository -y ppa:webupd8team/java
  sudo apt-get update
  sudo apt-get install --allow-unauthenticated oracle-java8-installer
  echo "installing bazel on circle-ci ..."
  echo "deb [arch=amd64] http://storage.googleapis.com/bazel-apt stable jdk1.8" | sudo tee /etc/apt/sources.list.d/bazel.list
  curl https://bazel.build/bazel-release.pub.gpg | sudo apt-key add -
  sudo apt-get update && sudo apt-get install bazel
  exit 0
fi

echo "not supported CI environment ... exit 1"
exit 1