BIN_DIR=bin
DIST_DIR=dist
SRC_DIR=gocsv
EXECUTABLE=gocsv

.DEFAULT_GOAL := bin
.PHONY: cleanall dist bin

dist:
	bash scripts/build-dist.sh

tag:
	bash scripts/update-latest-tag.sh

unit-test:
	cd $(SRC_DIR) && go test -cover

bin:
	rm -rf $(BIN_DIR)
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(EXECUTABLE) ./$(SRC_DIR)

cleanall:
	rm -rf $(BIN_DIR)
	rm -rf $(DIST_DIR)
