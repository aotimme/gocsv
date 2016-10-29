BUILD_DIR=./build
EXECUTABLE=gocsv

.DEFAULT_GOAL := build
.PHONY: clean build

build: clean
	mkdir ${BUILD_DIR}
	go build -o ${BUILD_DIR}/${EXECUTABLE}

clean:
	rm -rf ${BUILD_DIR}
