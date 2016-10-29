BIN_DIR=./bin
DIST_DIR=./dist
EXECUTABLE=gocsv

.DEFAULT_GOAL := bin
.PHONY: cleanall dist bin

dist:
	rm -rf ${DIST_DIR}
	mkdir ${DIST_DIR}
	# Build for Mac OS X
	mkdir ${DIST_DIR}/gocsv-darwin-amd64
	env GOOS=darwin GOARCH=amd64 go build -o ${DIST_DIR}/gocsv-darwin-amd64/${EXECUTABLE}
	cd ${DIST_DIR} && zip -r gocsv-darwin-amd64.zip gocsv-darwin-amd64
	rm -r ${DIST_DIR}/gocsv-darwin-amd64
	# Build for Linux
	mkdir ${DIST_DIR}/gocsv-linux-amd64
	env GOOS=linux GOARCH=amd64 go build -o ${DIST_DIR}/gocsv-linux-amd64/${EXECUTABLE}
	cd ${DIST_DIR} && zip -r gocsv-linux-amd64.zip gocsv-linux-amd64
	rm -r ${DIST_DIR}/gocsv-linux-amd64
	# Build for Windows
	mkdir ${DIST_DIR}/gocsv-windows-amd64
	env GOOS=windows GOARCH=amd64 go build -o ${DIST_DIR}/gocsv-windows-amd64/${EXECUTABLE}
	cd ${DIST_DIR} && zip -r gocsv-windows-amd64.zip gocsv-windows-amd64
	rm -r ${DIST_DIR}/gocsv-windows-amd64

bin:
	rm -rf ${BIN_DIR}
	mkdir -p ${BIN_DIR}
	go build -o ${BIN_DIR}/${EXECUTABLE}

cleanall:
	rm -rf ${BIN_DIR}
	rm -rf ${DIST_DIR}
