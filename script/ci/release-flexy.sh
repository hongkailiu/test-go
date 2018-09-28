#!/bin/bash

set -e

readonly GO_VERSION="$(go version)"
echo "go version: ${GO_VERSION}"
if [[ ! "${GO_VERSION}" == *"go1.10"* ]]; then
  echo "skip release up to go version"
  exit 0
fi

echo "RELEASE VAR: ${RELEASE}"
if [[ "${RELEASE}" != "true" ]]; then
  echo "skip release up to env var value"
  exit 0
fi

readonly SOURCE_FOLDER="$(dirname "$(readlink -f "${0}")")"
readonly APP_FOLDER="$(dirname "$(dirname "${SOURCE_FOLDER}")")"
readonly BUILD_DIR="${APP_FOLDER}/build"
readonly RELEASE_DIR="${BUILD_DIR}/release"

readonly BUILD_FILE="${APP_FOLDER}/build/flexy.test"
readonly PKG_DIR_NAME="pkg-flexy"
readonly PKG_DIR="${APP_FOLDER}/build/${PKG_DIR_NAME}"

readonly VERSION=$(git describe --tags --always --dirty)
readonly GO_OS="$(uname -s)"
readonly GO_ARCH="$(uname -m)"

mkdir -p "${PKG_DIR}"
cp -fv "${BUILD_FILE}" "${PKG_DIR}/"
cp -rfv "${APP_FOLDER}/test_files" "${PKG_DIR}/"

readonly PKG_BASENAME="flexy-${VERSION}-${GO_OS}-${GO_ARCH}.tar.gz"
readonly PKG_FULLNAME="${BUILD_DIR}/${PKG_BASENAME}"

current_dir="$(pwd)"
cd "${BUILD_DIR}" || exit 1
tar -czf "${PKG_FULLNAME}" --transform "s/${PKG_DIR_NAME}/flexy/" "${PKG_DIR_NAME}"
cd "${current_dir}" || exit 1

if [[ ! -f "${PKG_FULLNAME}" ]]; then
  echo "pkg file does not exits: ${PKG_FULLNAME}"
  exit 1
fi

rm -rf "${RELEASE_DIR}"
mkdir -p "${RELEASE_DIR}"
readonly REPO_NAME="svt-release"
readonly GH_TOKEN=${GH_TOKEN}
readonly REPO_URL="https://${GH_TOKEN}:x-oauth-basic@github.com/cduser/${REPO_NAME}.git"

current_dir="$(pwd)"
cd "${RELEASE_DIR}" || exit 1

git clone "${REPO_URL}"
cd "${REPO_NAME}"
git checkout -b tempB
cp -f "${PKG_FULLNAME}" .
git add "${PKG_BASENAME}"
if [[ -n "${TRAVIS}" ]]; then
  echo "release by travis ci to branch: travis_${TRAVIS_BUILD_NUMBER}"
  git config user.email "cduser@@users.noreply.github.com"
  git config user.name "CD User"
  msg_body_line1="TRAVIS_BUILD_NUMBER: ${TRAVIS_BUILD_NUMBER}"
  msg_body_line2="TRAVIS_BUILD_ID: ${TRAVIS_BUILD_ID}"
  msg_body_line3="TRAVIS_JOB_NUMBER: ${TRAVIS_JOB_NUMBER}"
  git commit -m "travis: ${PKG_BASENAME}" -m "${msg_body_line1}" -m "${msg_body_line2}" -m "${msg_body_line3}"
  git push origin "HEAD:travis_${TRAVIS_BUILD_NUMBER}"
else
  echo "release by dev to branch: dev_${HOSTNAME}_${USERNAME}"
  git commit -m "dev: ${PKG_BASENAME}"
  git push origin "HEAD:dev_${HOSTNAME}_${USERNAME}"
fi
cd "${current_dir}" || exit 1