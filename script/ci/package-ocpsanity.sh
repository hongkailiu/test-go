#!/bin/bash

set -e

readonly SOURCE_FOLDER="$(dirname "$(readlink -f "${0}")")"
readonly APP_FOLDER="$(dirname "$(dirname "${SOURCE_FOLDER}")")"
readonly BUILD_DIR="${APP_FOLDER}/build"
readonly RELEASE_DIR="${BUILD_DIR}/release"

readonly BUILD_FILE="${APP_FOLDER}/build/ocpsanity"
readonly PKG_DIR_NAME="pkg-ocpsanity"
readonly PKG_DIR="${APP_FOLDER}/build/${PKG_DIR_NAME}"

readonly VERSION=$(git describe --tags --always --dirty)
readonly GO_OS="$(uname -s)"
readonly GO_ARCH="$(uname -m)"

rm -rfv "${PKG_DIR}"

mkdir -p "${PKG_DIR}/build"
cp -fv "${BUILD_FILE}" "${PKG_DIR}/build"

readonly PKG_BASENAME="ocpsanity-${VERSION}-${GO_OS}-${GO_ARCH}.tar.gz"
readonly PKG_FULLNAME="${BUILD_DIR}/${PKG_BASENAME}"

rm -fv "${PKG_FULLNAME}"

current_dir="$(pwd)"
cd "${BUILD_DIR}" || exit 1
tar -czf "${PKG_FULLNAME}" --transform "s/${PKG_DIR_NAME}/ocpsanity/" "${PKG_DIR_NAME}"
cd "${current_dir}" || exit 1

