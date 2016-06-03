TOPDIR=$(CURDIR)/../../..
BUILDER_IMAGE=octav/adminweb-builder
REPOSITORY_HOST=asia.gcr.io
PROJECT_ID=$(shell gcloud config list 2>&1 | grep 'project = ' | sed -e 's/project = //')
CONTAINER=app
APPNAME=adminweb
ADMINWEB_DIR=$(TOPDIR)/adminweb
TAG:=$(shell date +"%Y%m%d.%H%M%S")
IMAGE_NAME_BASE:=$(REPOSITORY_HOST)/$(PROJECT_ID)/$(APPNAME)
IMAGE_NAME_LATEST:=$(IMAGE_NAME_BASE):latest
IMAGE_NAME:=$(IMAGE_NAME_BASE):$(TAG)

.PHONY: docker-build

$(CURDIR)/cpanfile:
	cp $(ADMINWEB_DIR)/cpanfile $(CURDIR)/cpanfile

build-local-lib:
	docker build -f Dockerfile.build -t octav/adminweb-builder .

$(CURDIR)/perl.tar.gz: $(CURDIR)/adminweb-builder.cid
	docker cp `cat adminweb-builder.cid`:/perl.tar.gz .

$(CURDIR)/adminweb-local-lib.tar.gz: $(CURDIR)/adminweb-builder.cid
	docker cp `cat adminweb-builder.cid`:/adminweb-local-lib.tar.gz .

$(CURDIR)/adminweb-builder.cid: build-local-lib
	@-rm adminweb-builder.cid
	docker run --cidfile adminweb-builder.cid $(BUILDER_IMAGE)

$(CURDIR)/adminweb.tar.gz: $(wildcard $(ADMINWEB_DIR)/**/*)
	tar -C $(ADMINWEB_DIR)/.. --exclude adminweb/local --exclude adminweb/tmp -cvzf $(CURDIR)/adminweb.tar.gz adminweb

$(CURDIR)/Octav.pm: $(TOPDIR)/p5/lib/WebService/Octav.pm
	cp  $^ .

docker-build:
	$(MAKE) -C $(TOPDIR) docker-env-ready
	@docker build -t octav/$(APPNAME) .

docker: $(CURDIR)/cpanfile $(CURDIR)/adminweb.tar.gz $(CURDIR)/adminweb-local-lib.tar.gz $(CURDIR)/perl.tar.gz $(CURDIR)/Octav.pm docker-build

publish:
	$(MAKE) -C $(TOPDIR) gke-publish IMAGE_NAME=$(IMAGE_NAME) APPNAME=$(APPNAME) IMAGE_NAME_LATEST=$(IMAGE_NAME_LATEST) \

deploy:
	$(MAKE) -C $(TOPDIR) gke-deploy IMAGE_NAME=$(IMAGE_NAME) APPNAME=$(APPNAME) CONTAINER=$(CONTAINER)

publish-deploy: publish deploy

clean-docker-images:
	$(MAKE) -C $(TOPDIR) clean-docker-images

clean: clean-docker-images

run:
	@docker run --rm -p 5000:5000 octav/$(APPNAME) 