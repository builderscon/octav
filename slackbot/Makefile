INTERNAL_BIN_DIR = $(CURDIR)/_internal_bin
GOVERSION=$(shell go version)
GOOS=$(word 1,$(subst /, ,$(word 4, $(GOVERSION))))
GOARCH=$(word 2,$(subst /, ,$(word 4, $(GOVERSION))))
BIN_DIR=$(CURDIR)/_bin/$(GOOS)/$(GOARCH)

installdeps:
	@PATH=$(INTERNAL_BIN_DIR):$(PATH) glide up

$(INTERNAL_BIN_DIR)/hsup:
	@mkdir -p $(INTERNAL_BIN_DIR)
	@echo "Building hsup..."
	@go get -d github.com/lestrrat/go-hsup
	@go build -o $@ github.com/lestrrat/go-hsup/cmd/hsup

buildspec: $(INTERNAL_BIN_DIR)/hsup ../spec/v1/slackbot.json
	@echo "Regenerating files based on ../spec/v1/slackbot.json"
	@PATH=$(CURDIR)/_internal_bin:$(PATH) hsup -s ../spec/v1/slackbot.json -d . -O

$(BIN_DIR)/slackbot-debug: cmd/slackbot/slackbot.go $(SRC_FILES)
	@echo "Building $(@F) for $(GOOS)/$(GOARCH)..."
	go build -tags debug0 -o $@ $<

glide: $(INTERNAL_BIN_DIR)/glide

$(INTERNAL_BIN_DIR)/glide: $(INTERNAL_BIN_DIR)/glide-$(GOOS)-$(GOARCH)

$(INTERNAL_BIN_DIR)/glide-$(GOOS)-$(GOARCH):
	@mkdir -p $(INTERNAL_BIN_DIR)
	@echo "Installing glide for $(GOOS)/$(GOARCH)..."
	@wget -O - https://github.com/Masterminds/glide/releases/download/0.10.0/glide-0.10.0-$(GOOS)-$(GOARCH).tar.gz | tar xvz
	@mv $(GOOS)-$(GOARCH)/glide $(INTERNAL_BIN_DIR)/glide-$(GOOS)-$(GOARCH)
	@rm -rf $(GOOS)-$(GOARCH)
	@ln -sf glide-$(GOOS)-$(GOARCH) $(INTERNAL_BIN_DIR)/glide 

slackbot-debug: $(BIN_DIR)/slackbot-debug

test:
	@go test -v $(GO_TAGS_OPT) $(shell PATH=$(INTERNAL_BIN_DIR):$(PATH) glide novendor)

