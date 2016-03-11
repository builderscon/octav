INTERNAL_BIN_DIR = _internal_bin
INTERNAL_CMDS    = gendb genmodel gentransport
BIN_DIR          = _bin

ifndef OCTAV_DB_NAME
	OCTAV_DB_NAME = octav
endif

GO_TEST_CMD      = OCTAV_TEST_DSN="root:@/$(OCTAV_DB_NAME)?parseTime=true" go test -v $(GO_TAGS_OPT)
ifdef OCTAV_TEST_DSN
	GO_TEST_CMD = go test -v $(GO_TAGS_OPT)
endif

.PHONY: install_deps build_internal_cmd generate build run all clean test init_db
all: test build

install_deps:
	go get -d -v -t $(GO_TAGS_OPT) ./...

build_internal_cmd:
	@for cmd in $(INTERNAL_CMDS); do\
		go build -o $(INTERNAL_BIN_DIR)/"$$cmd" internal/cmd/"$$cmd"/*.go;\
	done

generate: build_internal_cmd
	@echo "go generate ./..."
	@PATH=$(CURDIR)/_internal_bin:$(PATH) go generate ./...

build: generate install_deps
	go build -o $(CURDIR)/_bin/octav cmd/octav/*.go

test: generate
	$(GO_TEST_CMD) ./...

run: build
	./_bin/octav

clean:
	rm _internal_bin/*
	rm _bin/*

init_db:
	mysql -u root -e "DROP DATABASE IF EXISTS $(OCTAV_DB_NAME)"
	mysql -u root -e "CREATE DATABASE $(OCTAV_DB_NAME)"
	mysql -u root octav < sql/octav.sql
