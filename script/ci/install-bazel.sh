#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

readonly BAZELISK_VERSION=v0.0.8

if [[ "${TRAVIS:-false}" == "true" ]]; then
  echo "installing bazel on travis-ci ..."
  echo "deb [arch=amd64] http://storage.googleapis.com/bazel-apt stable jdk1.8" | sudo tee /etc/apt/sources.list.d/bazel.list
  curl https://bazel.build/bazel-release.pub.gpg | sudo apt-key add -
  sudo apt-get update && sudo apt-get install bazel
  curl -OL https://github.com/bazelbuild/bazelisk/releases/download/${BAZELISK_VERSION}/bazelisk-linux-amd64
  sudo mv ./bazelisk-linux-amd64 /usr/bin/bazel
  sudo chmod +x /usr/bin/bazel
  exit 0
fi

if [[ "${CIRCLECI:-false}" == "true" ]]; then
  echo "using my own image now ... skipping installation ..."
  exit 0
  echo "installing openjdk on circle-ci ..."
  #sudo apt-get install -y software-properties-common
  #sudo add-apt-repository -y ppa:webupd8team/java
  #sudo apt-get update
  #sudo apt-get install --allow-unauthenticated oracle-java8-installer
  #sudo add-apt-repository -y ppa:openjdk-r/ppa
  #sudo apt-get update
  sudo apt-get install openjdk-8-jdk
  echo "installing patch on circle-ci ..."
  sudo apt-get install patch
  echo "installing bazel on circle-ci ..."
  curl -OL https://github.com/bazelbuild/bazelisk/releases/download/${BAZELISK_VERSION}/bazelisk-linux-amd64
  sudo mv ./bazelisk-linux-amd64 /usr/bin/bazel
  sudo chmod +x /usr/bin/bazel
  #echo "deb [arch=amd64] http://storage.googleapis.com/bazel-apt stable jdk1.8" | sudo tee /etc/apt/sources.list.d/bazel.list
  #curl https://bazel.build/bazel-release.pub.gpg | sudo apt-key add -
  #sudo apt-get update && sudo apt-get install bazel
fi

echo "not supported CI environment ... exit 1"
exit 1