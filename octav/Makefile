EXTRA_TAGS       =
OCTAV_EXECUTABLE = octav
INTERNAL_BIN_DIR = $(CURDIR)/_internal_bin
INTERNAL_CMDS    = gendb genmodel gentransport genoctavctl
GOVERSION=$(shell go version)
GOOS=$(word 1,$(subst /, ,$(word 4, $(GOVERSION))))
GOARCH=$(word 2,$(subst /, ,$(word 4, $(GOVERSION))))
HSUP_DIR=$(shell go list -f '{{ .Dir }}' github.com/lestrrat/go-hsup)
BIN_DIR=$(CURDIR)/_bin/$(GOOS)/$(GOARCH)
ifndef OCTAV_DB_NAME
	OCTAV_DB_NAME = octav
endif

SRC_FILES = $(wildcard *.go validator/*.go client/*.go db/*.go model/*.go)

.PHONY: installdeps build_internal_cmd generate build run all clean test initdb gke-apiserver $(INTERNAL_BIN_DIR)/glide

all: test build

installdeps:
	@PATH=$(INTERNAL_BIN_DIR):$(PATH) glide install

$(INTERNAL_BIN_DIR):
	@mkdir -p $(INTERNAL_BIN_DIR)

$(INTERNAL_BIN_DIR)/%: $(INTERNAL_BIN_DIR) internal/cmd/%/*.go
	@echo "Building $(@F)..."
	go build -o $@ $^

$(INTERNAL_BIN_DIR)/hsup: $(INTERNAL_BIN_DIR)
	@echo "Building hsup..."
	@go get -x -d github.com/lestrrat/go-hsup
	@cd $(HSUP_DIR) && PATH=$(INTERNAL_BIN_DIR):$(PATH) glide install && go build -o $@ cmd/hsup/*.go

glide: $(INTERNAL_BIN_DIR)/glide

$(INTERNAL_BIN_DIR)/glide: $(INTERNAL_BIN_DIR)/glide-$(GOOS)-$(GOARCH)

$(INTERNAL_BIN_DIR)/glide-$(GOOS)-$(GOARCH):
	@mkdir -p $(INTERNAL_BIN_DIR)
	@echo "Installing glide for $(GOOS)/$(GOARCH)..."
	@wget -O - https://github.com/Masterminds/glide/releases/download/0.10.2/glide-0.10.2-$(GOOS)-$(GOARCH).tar.gz | tar xvz
	@mv $(GOOS)-$(GOARCH)/glide $(INTERNAL_BIN_DIR)/glide-$(GOOS)-$(GOARCH)
	@rm -rf $(GOOS)-$(GOARCH)
	@ln -sf glide-$(GOOS)-$(GOARCH) $(INTERNAL_BIN_DIR)/glide 

build_internal_cmd: $(patsubst %,$(INTERNAL_BIN_DIR)/%, $(INTERNAL_CMDS))

generate: build_internal_cmd
	@echo "go generate ."
	@PATH=$(CURDIR)/_internal_bin:$(PATH) go generate .

test:
	@OCTAV_TEST_DSN=$(OCTAV_TEST_DSN) go test -v $(GO_TAGS_OPT) $(shell PATH=$(INTERNAL_BIN_DIR):$(PATH) glide novendor)

$(BIN_DIR)/octav: cmd/octav/octav.go $(SRC_FILES)
	@echo "Building $(@F) for $(GOOS)/$(GOARCH)..."
	go build -o $@ $<

octav: $(BIN_DIR)/octav

$(BIN_DIR)/octav-debug: cmd/octav/octav.go $(SRC_FILES)
	@echo "Building $(@F) for $(GOOS)/$(GOARCH)..."
	go build -tags debug0 -o $@ $<

octav-debug: $(BIN_DIR)/octav-debug

buildspec: $(INTERNAL_BIN_DIR)/hsup ../spec/v1/api.json
	@echo "Regenerating files based on ../spec/v1/api.json"
	@PATH=$(CURDIR)/_internal_bin:$(PATH) hsup -s ../spec/v1/api.json -d . -O

build: $(BIN_DIR)/octav

run: build
	$(BIN_DIR)/octav

clean:
	-rm _internal_bin/*
	-rm -rf _bin/*

initdb:
	$(MAKE) -C ../ initdb
