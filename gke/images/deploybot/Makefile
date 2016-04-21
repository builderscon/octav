APPNAME=deploybot
TOPDIR=$(CURDIR)/../../..
REPOSITORY_HOST=asia.gcr.io
PROJECT_ID=$(shell gcloud config list 2>&1 | grep 'project = ' | sed -e 's/project = //')

TAG:=$(shell date +"%Y%m%d.%H%M%S")
IMAGE_NAME_BASE:=$(REPOSITORY_HOST)/$(PROJECT_ID)/$(APPNAME)
IMAGE_NAME_LATEST:=$(IMAGE_NAME_BASE):latest
IMAGE_NAME:=$(IMAGE_NAME_BASE):$(TAG)

$(CURDIR)/deploybot-debug:
	cd $(TOPDIR)/worker/deploybot && \
	GOOS=linux GOARCH=amd64 go build -tags debug0 -o $(CURDIR)/deploybot-debug .

$(CURDIR)/deploybot:
	cd $(TOPDIR)/worker/deploybot && \
	GOOS=linux GOARCH=amd64 go build -o $(CURDIR)/deploybot .

deploybot-debug: $(CURDIR)/deploybot-debug
deploybot: $(CURDIR)/deploybot

docker-debug: $(CURDIR)/deploybot-debug
	cp $(CURDIR)/deploybot-debug $(CURDIR)/deploybot
	$(MAKE) docker

docker: $(CURDIR)/deploybot
	docker build -t octav/deploybot .

publish:
	$(MAKE) -C $(TOPDIR) gke-publish \
		IMAGE_NAME_LATEST=$(IMAGE_NAME_LATEST) \
		IMAGE_NAME=$(IMAGE_NAME) \
		APPNAME=$(APPNAME)

clean:
	@-rm deploybot
	@-rm deploybot-debug
