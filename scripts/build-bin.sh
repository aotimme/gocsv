#!/usr/bin/env bash

BIN_DIR=bin
SRC_DIR=src
EXECUTABLE=gocsv

GIT_HASH=$(git rev-parse HEAD)
VERSION=$(cat VERSION)
LD_FLAGS="-X main.VERSION=${VERSION} -X main.GIT_HASH=${GIT_HASH}"

rm -rf ${BIN_DIR}
mkdir -p ${BIN_DIR}
go build -ldflags "${LD_FLAGS}" -o ${BIN_DIR}/${EXECUTABLE} ./${SRC_DIR}