BIN_DIR=bin
DIST_DIR=dist
CMD_DIR=cmd
CSV_DIR=csv
EXECUTABLE=gocsv

.DEFAULT_GOAL := bin
.PHONY: cleanall dist bin

dist:
	bash scripts/build-dist.sh

tag:
	bash scripts/update-latest-tag.sh

test:
	cd $(CMD_DIR) && GO111MODULE=on go test -cover
	cd $(CSV_DIR) && GO111MODULE=on go test -cover

bin:
	bash scripts/build-bin.sh

cleanall:
	rm -rf $(BIN_DIR)
	rm -rf $(DIST_DIR)
