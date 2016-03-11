INTERNAL_BIN_DIR = _internal_bin
INTERNAL_CMDS    = gendb genmodel gentransport
BIN_DIR          = _bin

ifndef OCTAV_DB_NAME
	OCTAV_DB_NAME = octav
endif

SRC_FILES = *.go validator/*.go client/*.go db/*.go

GO_TEST_CMD      = OCTAV_TEST_DSN="root:@/$(OCTAV_DB_NAME)?parseTime=true" go test -v $(GO_TAGS_OPT)
ifdef OCTAV_TEST_DSN
	GO_TEST_CMD = go test -v $(GO_TAGS_OPT)
endif

.PHONY: install_deps build_internal_cmd generate build run all clean test init_db
all: test build

install_deps:
	go get -d -v -t $(GO_TAGS_OPT) ./...

$(INTERNAL_BIN_DIR)/%: internal/cmd/%/*.go
	@echo "Building $(@F)..."
	@go build -o $@ $^

build_internal_cmd: $(patsubst %,$(INTERNAL_BIN_DIR)/%, $(INTERNAL_CMDS))

generate: build_internal_cmd
	@echo "go generate ./..."
	@PATH=$(CURDIR)/_internal_bin:$(PATH) go generate ./...

test: generate
	$(GO_TEST_CMD) ./...

$(CURDIR)/_bin/octav: cmd/octav/octav.go $(SRC_FILES)
	@echo "Building $(@F)..."
	@go build -o $@ $<

build: $(CURDIR)/_bin/octav

run: build
	./_bin/octav

clean:
	rm _internal_bin/*
	rm _bin/*

init_db:
	mysql -u root -e "DROP DATABASE IF EXISTS $(OCTAV_DB_NAME)"
	mysql -u root -e "CREATE DATABASE $(OCTAV_DB_NAME)"
	mysql -u root octav < sql/octav.sql
