EXTRA_TAGS       =
OCTAV_EXECUTABLE = octav
INTERNAL_BIN_DIR = $(CURDIR)/_internal_bin
INTERNAL_CMDS    = gendb genmodel gentransport genoctavctl
GOVERSION=$(shell go version)
GOOS=$(word 1,$(subst /, ,$(word 4, $(GOVERSION))))
GOARCH=$(word 2,$(subst /, ,$(word 4, $(GOVERSION))))
BIN_DIR=$(CURDIR)/_bin/$(GOOS)/$(GOARCH)
ifndef OCTAV_DB_NAME
	OCTAV_DB_NAME = octav
endif

SRC_FILES = $(wildcard *.go validator/*.go client/*.go db/*.go model/*.go)

.PHONY: installdeps build_internal_cmd generate build run all clean test initdb gke-apiserver $(INTERNAL_BIN_DIR)/glide hsupsrc bindatasrc assets

all: test build

installdeps:
	@PATH=$(INTERNAL_BIN_DIR):$(PATH) glide install

hsupsrc:
	@go get -d github.com/lestrrat/go-hsup

bindatasrc:
	@go get -d github.com/lestrrat/go-bindata

$(INTERNAL_BIN_DIR):
	@mkdir -p $(INTERNAL_BIN_DIR)

$(INTERNAL_BIN_DIR)/%: internal/cmd/%/*.go
	@$(MAKE) $(INTERNAL_BIN_DIR)
	@echo "Building $(@F)..."
	go build -o $@ $^

$(INTERNAL_BIN_DIR)/hsup:
	@$(MAKE) $(INTERNAL_BIN_DIR)
	@echo "Building hsup..."
	@cd $(shell go list -f '{{ .Dir }}' github.com/lestrrat/go-hsup) && PATH=$(INTERNAL_BIN_DIR):$(PATH) glide install && go build -o $@ cmd/hsup/*.go

$(INTERNAL_BIN_DIR)/go-bindata:
	@$(MAKE) $(INTERNAL_BIN_DIR)
	@echo "Building bindata..."
	@cd $(shell go list -f '{{ .Dir }}' github.com/lestrrat/go-bindata) && PATH=$(INTERNAL_BIN_DIR):$(PATH) go build -o $@ ./go-bindata

glide: $(INTERNAL_BIN_DIR)/glide

$(INTERNAL_BIN_DIR)/glide: $(INTERNAL_BIN_DIR)/glide-$(GOOS)-$(GOARCH)

$(INTERNAL_BIN_DIR)/glide-$(GOOS)-$(GOARCH):
	@mkdir -p $(INTERNAL_BIN_DIR)
	@echo "Installing glide for $(GOOS)/$(GOARCH)..."
	@wget -O - https://github.com/Masterminds/glide/releases/download/v0.12.0/glide-v0.12.0-$(GOOS)-$(GOARCH).tar.gz | tar xvz
	@mv $(GOOS)-$(GOARCH)/glide $(INTERNAL_BIN_DIR)/glide-$(GOOS)-$(GOARCH)
	@rm -rf $(GOOS)-$(GOARCH)
	@ln -sf glide-$(GOOS)-$(GOARCH) $(INTERNAL_BIN_DIR)/glide 

build_internal_cmd: $(patsubst %,$(INTERNAL_BIN_DIR)/%, $(INTERNAL_CMDS))

generate: build_internal_cmd $(INTERNAL_BIN_DIR)/go-bindata
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

assets: $(INTERNAL_BIN_DIR)/go-bindata
	@echo "Regenerating assets"
	@cd assets; PATH=$(CURDIR)/_internal_bin:$(PATH) go-bindata -o=assets.go -pkg=assets -ignore=assets.go -ignore=\\.swp templates/...

