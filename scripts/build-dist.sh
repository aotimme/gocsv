#!/usr/bin/env bash

DIST_DIR=$(pwd)/dist
EXECUTABLE=gocsv

GIT_HASH=$(git rev-parse HEAD)
VERSION=$(git describe --tags HEAD)
LD_FLAGS="-X github.com/aotimme/gocsv/cmd.VERSION=${VERSION} -X github.com/aotimme/gocsv/cmd.GIT_HASH=${GIT_HASH}"

rm -rf ${DIST_DIR}
mkdir ${DIST_DIR}

# Create an array of goos:goarch pairs
options=(
  "darwin:amd64"
  "darwin:arm64"
  "freebsd:amd64"
  "freebsd:arm64"
  "linux:386"
  "linux:amd64"
  "linux:arm"
  "linux:arm64"
  "linux:ppc64le"
  "linux:riscv64"
  "windows:amd64"
  "windows:arm64"
)

echo "Building into ${DIST_DIR}/..."
for option in "${options[@]}"; do
  IFS=':' read -r goos goarch <<< "$option"

  folder="${EXECUTABLE}-${goos}-${goarch}"
  echo "Building: ${folder}..."
  BIN_DIR=${DIST_DIR}/${folder}
  rm -rf ${BIN_DIR}
  mkdir -p ${BIN_DIR}
  binary=${EXECUTABLE}
  if [ "$goos" == "windows" ]; then
    binary="${EXECUTABLE}.exe"
  fi
  GOOS=${goos} GOARCH=${goarch} GO111MODULE=on go build -ldflags "${LD_FLAGS}" -o ${BIN_DIR}/${binary}
  cd ${DIST_DIR}
  zip -rq ${folder}.zip ${folder}
  rm -rf ${folder}
  cd - >/dev/null 2>&1
done