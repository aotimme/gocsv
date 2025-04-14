BIN_DIR=bin
DIST_DIR=dist
EXECUTABLE=gocsv

.DEFAULT_GOAL := bin
.PHONY: cleanall dist bin

dist:
	bash scripts/build-dist.sh

test:
	go test -cover ./...

bin:
	bash scripts/build-bin.sh

cleanall:
	rm -rf $(BIN_DIR)
	rm -rf $(DIST_DIR)
