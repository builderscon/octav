INTERNAL_BIN_DIR = _internal_bin
INTERNAL_CMDS    = gendb genmodel gentransport
BIN_DIR          = _bin
GCLOUD_CONFIG_DIR = .gcloud
GCLOUD_PROJECT_ID = builderscon-1248
CLOUDSQL_ADDRESS_FILE = $(GCLOUD_CONFIG_DIR)/sql/address
CLOUDSQL_SSL_FILES    = \
	$(GCLOUD_CONFIG_DIR)/sql/server-ca.pem \
	$(GCLOUD_CONFIG_DIR)/sql/client-key.pem \
	$(GCLOUD_CONFIG_DIR)/sql/client-cert.pem

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
	@echo "Building $(@F)..."
	@go build -o $@ $<

buildspec: $(INTERNAL_BIN_DIR)/hsup ../spec/v1/api.json
	@echo "Regenerating files based on ../spec/v1/api.json"
	@PATH=$(CURDIR)/_internal_bin:$(PATH) hsup -s ../spec/v1/api.json -d . -O

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

appserve:
	@echo "Starting appengine local server via 'goapp serve'"
	@cd appengine && goapp serve

appdeploy:
	@echo "Deploying to App Engine $(GCLOUD_PROJECT_ID) ($(GCLOUD_PROJECT_ID).appspot.com)"
	@cd appengine && goapp deploy -application $(GCLOUD_PROJECT_ID) app.yaml

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
	@mysql -uroot -h `cat $(GCLOUD_CONFIG_DIR)/sql/address` -p --ssl-ca=$(GCLOUD_CONFIG_DIR)/sql/server-ca.pem --ssl-cert=$(GCLOUD_CONFIG_DIR)/sql/client-cert.pem --ssl-key=$(GCLOUD_CONFIG_DIR)/sql/client-key.pem
