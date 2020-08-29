#!/usr/bin/env bash

BUILD_DIR=$(pwd)/build
DIST_DIR=$(pwd)/dist
EXECUTABLE=gocsv

GIT_HASH=$(git rev-parse HEAD)
VERSION=$(cat VERSION)
LD_FLAGS="-X github.com/aotimme/gocsv/cmd.VERSION=${VERSION} -X github.com/aotimme/gocsv/cmd.GIT_HASH=${GIT_HASH}"

rm -rf ${BUILD_DIR}
mkdir ${BUILD_DIR}

rm -rf ${DIST_DIR}
mkdir ${DIST_DIR}

cd ${BUILD_DIR}
# Use `xgo` to truly handle cross-compiling, mainly due to gocsv's dependency
# on `go-sqlite3`, which is a cgo package.
# See: https://github.com/mattn/go-sqlite3#cross-compile
xgo -ldflags "${LD_FLAGS}" github.com/aotimme/gocsv

# Move files to `dist` and zip them.
for file in $(ls);
do
  folder=${file}
  binary=${EXECUTABLE}
  if [ ${file: -4} == ".exe" ];
  then
    binary="${EXECUTABLE}.exe"
    folder=$(echo ${folder} | sed -E 's/.exe$//g')
  fi
  mkdir ${DIST_DIR}/${folder}
  mv ${file} ${DIST_DIR}/${folder}/${binary}
  cd ${DIST_DIR}
  zip -rq ${folder}.zip ${folder}
  rm -r ${folder}
  cd ${BUILD_DIR}
done

rm -r ${BUILD_DIR}