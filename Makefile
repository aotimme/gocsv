BIN_DIR=bin
DIST_DIR=dist
SRC_DIR=src
EXECUTABLE=gocsv

.DEFAULT_GOAL := bin
.PHONY: cleanall dist bin

dist:
	bash scripts/build-dist.sh

tag:
	bash scripts/update-latest-tag.sh

test:
	cd $(SRC_DIR) && go test -cover

bin:
	bash scripts/build-bin.sh

cleanall:
	rm -rf $(BIN_DIR)
	rm -rf $(DIST_DIR)
