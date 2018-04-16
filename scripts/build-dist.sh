#!/usr/bin/env bash

BIN_DIR=bin
DIST_DIR=dist
SRC_DIR=src
EXECUTABLE=gocsv

rm -rf ${DIST_DIR}
mkdir ${DIST_DIR}
for os in darwin windows linux; do
	for arch in amd64; do
		basename=gocsv-${os}-${arch}
		mkdir ${DIST_DIR}/${basename}
		env GOOS=${os} GOARCH=${arch} go build -o ${DIST_DIR}/${basename}/${EXECUTABLE} ./${SRC_DIR}
		cd ${DIST_DIR} && zip -rq ${basename}.zip ${basename}
    cd ~-
		rm -r ${DIST_DIR}/${basename}
	done
done
