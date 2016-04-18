APPNAME=acmebot
TOPDIR=$(CURDIR)/../../..
REPOSITORY_HOST=asia.gcr.io
PROJECT_ID=$(shell gcloud config list 2>&1 | grep 'project = ' | sed -e 's/project = //')

TAG:=$(shell date +"%Y%m%d.%H%M%S")
IMAGE_NAME:=$(REPOSITORY_HOST)/$(PROJECT_ID)/$(APPNAME):$(TAG)

$(CURDIR)/acmebot-debug:
	cd $(TOPDIR)/worker/acmebot && \
	GOOS=linux GOARCH=amd64 go build -tags debug0 -o $(CURDIR)/acmebot-debug .

$(CURDIR)/acmebot:
	cd $(TOPDIR)/worker/acmebot && \
	GOOS=linux GOARCH=amd64 go build -o $(CURDIR)/acmebot .

acmebot-debug: $(CURDIR)/acmebot-debug
acmebot: $(CURDIR)/acmebot

docker-debug: $(CURDIR)/acmebot-debug
	cp $(CURDIR)/acmebot-debug $(CURDIR)/acmebot
	$(MAKE) docker

docker: $(CURDIR)/acmebot
	docker build -t octav/acmebot .

publish:
	$(MAKE) -C $(TOPDIR) gke-publish IMAGE_NAME=$(IMAGE_NAME) APPNAME=$(APPNAME)

clean:
	@-rm acmebot
	@-rm acmebot-debug
