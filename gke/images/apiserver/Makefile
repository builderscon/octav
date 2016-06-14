APPNAME=apiserver
TOPDIR=$(CURDIR)/../../..
OCTAV_DIR = $(TOPDIR)/octav
REPOSITORY_HOST=asia.gcr.io
PROJECT_ID=$(shell gcloud config list 2>&1 | grep 'project = ' | sed -e 's/project = //')

# Unfortunately this had to be replicated from top Makefile
SRC_FILES = $(wildcard $(OCTAV_DIR)/*.go $(OCTAV_DIR)/validator/*.go $(OCTAV_DIR)/client/*.go $(OCTAV_DIR)/db/*.go $(OCTAV_DIR)/model/*.go)

TAG:=$(shell date +"%Y%m%d.%H%M%S")
IMAGE_NAME_BASE:=$(REPOSITORY_HOST)/$(PROJECT_ID)/$(APPNAME)
IMAGE_NAME_LATEST:=$(IMAGE_NAME_BASE):latest
IMAGE_NAME:=$(IMAGE_NAME_BASE):$(TAG)

.PHONY: all octav always-rebuild

$(CURDIR)/octav: $(SRC_FILES)
	@echo "Building octav"
	cd $(OCTAV_DIR) ; make GOOS=linux GOARCH=amd64 octav; cd $(CURDIR)
	@cp $(OCTAV_DIR)/_bin/linux/amd64/octav .

$(CURDIR)/octav-debug: $(SRC_FILES)
	@echo "Building octav-debug"
	cd $(OCTAV_DIR) ; make GOOS=linux GOARCH=amd64 octav-debug; cd $(CURDIR)
	@mv $(OCTAV_DIR)/_bin/linux/amd64/octav-debug .

build-octav: $(CURDIR)/octav
build-octav-debug: $(CURDIR)/octav-debug

docker-debug: build-octav-debug
	$(MAKE) -C $(TOPDIR) docker-env-ready
	cp octav-debug octav
	@docker build -t octav/$(APPNAME) .

docker: build-octav
	$(MAKE) -C $(TOPDIR) docker-env-ready
	@docker build -t octav/$(APPNAME) .

publish:
	$(MAKE) -C $(TOPDIR) gke-publish IMAGE_NAME=$(IMAGE_NAME) APPNAME=$(APPNAME) IMAGE_NAME_LATEST=$(IMAGE_NAME_LATEST)

clean-docker-images:
	$(MAKE) -C $(TOPDIR) clean-docker-images

clean_executables:
	@echo "Deleting generated executables"
	@-rm octav octav-debug

clean: clean-docker-images clean_executables

deploy:
	$(MAKE) -C $(TOPDIR) gke-deploy IMAGE_NAME=$(IMAGE_NAME) APPNAME=$(APPNAME)

publish-deploy: publish deploy
