INTERNAL_BIN_DIR = _internal_bin
INTERNAL_CMDS    = gendb genmodel gentransport
BIN_DIR          = _bin
OCTAV_DB_NAME    = octav

.PHONY: install_deps build_internal_cmd generate build run all clean test init_db
all: test build

install_deps:
	go get -d -v -t ./...

build_internal_cmd: install_deps
	@for cmd in $(INTERNAL_CMDS); do\
		go build -o $(INTERNAL_BIN_DIR)/"$$cmd" internal/cmd/"$$cmd"/*.go;\
	done

generate: build_internal_cmd
	PATH=$(CURDIR)/_internal_bin:$(PATH) go generate ./...

build: generate
	go build -o $(CURDIR)/_bin/octav cmd/octav/*.go

test: generate
	OCTAV_TEST_DSN="root:@/$(OCTAV_DB_NAME)?parseTime=true" go test ./...

run: build
	./_bin/octav

clean:
	rm _internal_bin/*
	rm _bin/*

init_db:
	mysql -u root -e "DROP DATABASE IF EXISTS $(OCTAV_DB_NAME)"
	mysql -u root -e "CREATE DATABASE $(OCTAV_DB_NAME)"
	mysql -u root octav < sql/octav.sql
