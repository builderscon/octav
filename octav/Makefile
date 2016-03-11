INTERNAL_BIN_DIR = _internal_bin
INTERNAL_CMDS    = gendb genmodel gentransport
BIN_DIR          = _bin

.PHONY: install_deps build_internal_cmd generate build run all clean
all: install_deps build_internal_cmd generate build

install_deps:
	go get -d -v ./...

build_internal_cmd:
	@for cmd in $(INTERNAL_CMDS); do\
		go build -o $(INTERNAL_BIN_DIR)/"$$cmd" internal/cmd/"$$cmd"/*.go;\
	done

generate:
	PATH=$(CURDIR)/_internal_bin:$(PATH) go generate ./...

build:
	go build -o $(CURDIR)/_bin/octav cmd/octav/*.go

run:
	./_bin/octav

clean:
	rm _internal_bin/*
	rm _bin/*
