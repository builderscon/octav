APPNAME=slackbot
TOPDIR=$(CURDIR)/../../..
SLACKBOT_DIR = $(TOPDIR)/slackbot
REPOSITORY_HOST=asia.gcr.io
PROJECT_ID=$(shell gcloud config list 2>&1 | grep 'project = ' | sed -e 's/project = //')

# Unfortunately this had to be replicated from top Makefile
SRC_FILES = $(wildcard $(SLACKBOT_DIR)/*.go $(SLACKBOT_DIR)/validator/*.go $(SLACKBOT_DIR)/client/*.go $(SLACKBOT_DIR)/db/*.go $(SLACKBOT_DIR)/model/*.go)

TAG:=$(shell date +"%Y%m%d.%H%M%S")
IMAGE_NAME:=$(REPOSITORY_HOST)/$(PROJECT_ID)/$(APPNAME):$(TAG)

.PHONY: all slackbot always-rebuild

$(CURDIR)/slackbot: $(SRC_FILES)
	@echo "Building slackbot"
	cd $(SLACKBOT_DIR) ; make GOOS=linux GOARCH=amd64 slackbot; cd $(CURDIR)
	@cp $(SLACKBOT_DIR)/_bin/linux/amd64/slackbot .

$(CURDIR)/slackbot-debug: $(SRC_FILES)
	@echo "Building slackbot-debug"
	cd $(SLACKBOT_DIR) ; make GOOS=linux GOARCH=amd64 slackbot-debug; cd $(CURDIR)
	@mv $(SLACKBOT_DIR)/_bin/linux/amd64/slackbot-debug .

build-slackbot: $(CURDIR)/slackbot
build-slackbot-debug: $(CURDIR)/slackbot-debug

dm-running:
ifneq (,$(strip $(shell docker-machine env default 2>&1 | grep -o 'not running')))
	@echo "Starting docker-machine (default)"
	@docker-machine start default
endif

docker-debug: dm-running build-slackbot-debug
	$(MAKE) -C $(TOPDIR) docker-env-ready
	cp slackbot-debug slackbot
	@docker build -t octav/$(APPNAME) .

docker: dm-running build-slackbot
	$(MAKE) -C $(TOPDIR) docker-env-ready
	@docker build -t octav/$(APPNAME) .

publish:
	$(MAKE) -C $(TOPDIR) gke-publish IMAGE_NAME=$(IMAGE_NAME) APPNAME=$(APPNAME)

clean-docker-images:
	$(MAKE) -C $(TOPDIR) clean-docker-images

clean_executables:
	@echo "Deleting generated executables"
	@-rm slackbot

clean: clean-docker-images clean_executables

deploy:
	$(MAKE) -C $(TOPDIR) gke-deploy IMAGE_NAME=$(IMAGE_NAME) APPNAME=$(APPNAME)

publish-deploy: publish deploy
