BUILD_DIR=./build
EXECUTABLE=gocsv

.DEFAULT_GOAL := build
.PHONY: clean build build-osx

build-osx:
	mkdir -p ${BUILD_DIR}/gocsv-mac-os-x
	go build -o ${BUILD_DIR}/gocsv-mac-os-x/${EXECUTABLE}
	cd ${BUILD_DIR} && zip -r gocsv-mac-os-x.zip ./gocsv-mac-os-x

build: clean build-osx

clean:
	rm -rf ${BUILD_DIR}
