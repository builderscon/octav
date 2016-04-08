TOPDIR=$(CURDIR)/../../..
REPOSITORY_HOST=asia.gcr.io
PROJECT_ID=$(shell gcloud config list 2>&1 | grep 'project = ' | sed -e 's/project = //')
CONTAINER=app
APPNAME=adminweb
ADMINWEB_DIR=$(TOPDIR)/adminweb
TAG:=$(shell date +"%Y%m%d.%H%M%S")
IMAGE_NAME:=$(REPOSITORY_HOST)/$(PROJECT_ID)/$(APPNAME):$(TAG)

.PHONY: docker-build

$(CURDIR)/cpanfile:
	cp $(ADMINWEB_DIR)/cpanfile $(CURDIR)/cpanfile

$(CURDIR)/adminweb.tar.gz: $(wildcard $(ADMINWEB_DIR)/**/*)
	tar -C $(ADMINWEB_DIR)/.. --exclude adminweb/local --exclude adminweb/tmp -cvzf $(CURDIR)/adminweb.tar.gz adminweb

$(CURDIR)/Octav.pm: $(TOPDIR)/p5/lib/WebService/Octav.pm
	cp  $^ .

docker-build:
	$(MAKE) -C $(TOPDIR) docker-env-ready
	@docker build -t octav/$(APPNAME) .

docker: $(CURDIR)/cpanfile $(CURDIR)/adminweb.tar.gz $(CURDIR)/Octav.pm docker-build

publish:
	$(MAKE) -C $(TOPDIR) gke-publish IMAGE_NAME=$(IMAGE_NAME) APPNAME=$(APPNAME)

deploy:
	$(MAKE) -C $(TOPDIR) gke-deploy IMAGE_NAME=$(IMAGE_NAME) APPNAME=$(APPNAME) CONTAINER=$(CONTAINER)

publish-deploy: publish deploy

clean-docker-images:
	$(MAKE) -C $(TOPDIR) clean-docker-images

clean: clean-docker-images

run:
	@docker run --rm -p 5000:5000 octav/$(APPNAME) 