BIN_DIR          = _bin
EXTRA_TAGS       =
INTERNAL_BIN_DIR = _internal_bin
INTERNAL_CMDS    = gendb genmodel gentransport
GCLOUD_CONFIG_DIR = .gcloud
CLOUDSQL_DIR=$(GCLOUD_CONFIG_DIR)/sql
CLOUDSQL_ADDRESS_FILE = $(CLOUDSQL_DIR)/address
CLOUDSQL_ADDRESS = $(shell cat $(CLOUDSQL_ADDRESS_FILE))
CLOUDSQL_SSL_FILES    = \
	$(CLOUDSQL_DIR)/server-ca.pem \
	$(CLOUDSQL_DIR)/client-key.pem \
	$(CLOUDSQL_DIR)/client-cert.pem
GOVERSION=$(shell go version)
GOOS=$(word 1,$(subst /, ,$(word 4, $(GOVERSION))))
GOARCH=$(word 2,$(subst /, ,$(word 4, $(GOVERSION))))

ifndef OCTAV_DB_NAME
	OCTAV_DB_NAME = octav
endif

SRC_FILES = $(wildcard *.go validator/*.go client/*.go db/*.go)

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

$(INTERNAL_BIN_DIR)/hsup:
	@echo "Building hsup..."
	@go get -d github.com/lestrrat/go-hsup
	@go build -o $@ github.com/lestrrat/go-hsup/cmd/hsup

build_internal_cmd: $(patsubst %,$(INTERNAL_BIN_DIR)/%, $(INTERNAL_CMDS))

generate: build_internal_cmd
	@echo "go generate ./..."
	@PATH=$(CURDIR)/_internal_bin:$(PATH) go generate ./...

test: generate buildspec
	$(GO_TEST_CMD) ./...

$(CURDIR)/_bin/octav: cmd/octav/octav.go $(SRC_FILES)
	@echo "Building $(@F) for $(GOOS)/$(GOARCH)..."
ifeq ($(DEBUG),1)
	@echo "DEBUG specified, enabling '-tags debug0'"
	$(eval $@_EXTRA_TAGS=-tags debug0)
endif
	go build $($@_EXTRA_TAGS) -o $@ $<

buildspec: $(INTERNAL_BIN_DIR)/hsup ../spec/v1/api.json
	@echo "Regenerating files based on ../spec/v1/api.json"
	@PATH=$(CURDIR)/_internal_bin:$(PATH) hsup -s ../spec/v1/api.json -d . -O

build: $(CURDIR)/_bin/octav

run: build
	./_bin/octav

clean:
	-rm _internal_bin/*
	-rm _bin/*

init_db:
	mysql -u root -e "DROP DATABASE IF EXISTS $(OCTAV_DB_NAME)"
	mysql -u root -e "CREATE DATABASE $(OCTAV_DB_NAME)"
	mysql -u root octav < sql/octav.sql

k8s-bootstrap:
	./devtools/k8s-bootstrap.sh

k8s-shutdown:
	./devtools/k8s-shutdown.sh

$(CLOUDSQL_ADDRESS_FILE):
ifeq (,$(wildcard $@))
	@echo "Missing '$@'"
	@echo "CloudSQL address needed. Please ask a builderscon admin to provide it for you"
	@exit 1
endif

$(CLOUDSQL_SSL_FILES):
ifeq (,$(wildcard $@))
	@echo "Missing '$@'"
	@echo "CloudSQL requires SSL certificates to connect. Please ask a builderscon admin to provide one for you"
	@exit 1
endif

cloudsql_files: $(CLOUDSQL_ADDRESS_FILE) $(CLOUDSQL_SSL_FILES)

cloudsql: cloudsql_files
	@echo "Connecting to mysql..."
	mysql -uroot -h $(CLOUDSQL_ADDRESS) -p --ssl-ca=$(GCLOUD_CONFIG_DIR)/sql/server-ca.pem --ssl-cert=$(GCLOUD_CONFIG_DIR)/sql/client-cert.pem --ssl-key=$(GCLOUD_CONFIG_DIR)/sql/client-key.pem

# This rule creates a secrets file so that it can be fed into kubectl.
# We go through this hoopla to create the secret so that we don't have to
# commit extra files that otherwise may reveal sensitive information.
#
# Kubernetes site shows that you can do this from the kubectl command line
# alone, but as of this writing at least kubectl that comes with the
# gcloud toolset doesn't, so... this workaround
cloudsql_secret: cloudsql_files
	./devtools/make_cloudsql_secret.sh | kubectl create -f -

