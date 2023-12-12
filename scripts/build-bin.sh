#!/usr/bin/env bash

BIN_DIR=$(pwd)/bin
EXECUTABLE=gocsv

GIT_HASH=$(git rev-parse HEAD)
VERSION=$(git describe --tags HEAD)
LD_FLAGS="-X github.com/aotimme/gocsv/cmd.VERSION=${VERSION} -X github.com/aotimme/gocsv/cmd.GIT_HASH=${GIT_HASH}"

rm -rf ${BIN_DIR}
mkdir -p ${BIN_DIR}
go build -ldflags "${LD_FLAGS}" -o ${BIN_DIR}/${EXECUTABLE}