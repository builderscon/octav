APPNAME=acmebot-k8s
TOPDIR=$(CURDIR)/../../..
REPOSITORY_HOST=asia.gcr.io
PROJECT_ID=$(shell gcloud config list 2>&1 | grep 'project = ' | sed -e 's/project = //')

TAG:=$(shell date +"%Y%m%d.%H%M%S")
IMAGE_NAME_BASE:=$(REPOSITORY_HOST)/$(PROJECT_ID)/$(APPNAME)
IMAGE_NAME_LATEST:=$(IMAGE_NAME_BASE):latest
IMAGE_NAME:=$(IMAGE_NAME_BASE):$(TAG)

$(CURDIR)/acmebot-k8s-debug:
	cd $(TOPDIR)/worker/acmebot-k8s && \
	GOOS=linux GOARCH=amd64 go build -tags debug0 -o $(CURDIR)/acmebot-k8s-debug .

$(CURDIR)/acmebot-k8s:
	cd $(TOPDIR)/worker/acmebot-k8s && \
	GOOS=linux GOARCH=amd64 go build -o $(CURDIR)/acmebot-k8s .

acmebot-k8s-debug: $(CURDIR)/acmebot-k8s-certupload-debug
acmebot-k8s: $(CURDIR)/acmebot-k8s-certupload

docker-debug: $(CURDIR)/acmebot-k8s-debug
	cp $(CURDIR)/acmebot-k8s-debug $(CURDIR)/acmebot-k8s-certupload
	$(MAKE) docker

docker: $(CURDIR)/acmebot-k8s
	docker build -t octav/acmebot-k8s .

publish:
	$(MAKE) -C $(TOPDIR) gke-publish \
                IMAGE_NAME_LATEST=$(IMAGE_NAME_LATEST) \
                IMAGE_NAME=$(IMAGE_NAME) \
                APPNAME=$(APPNAME)

clean:
	@-rm acmebot-k8s
	@-rm acmebot-k8s-debug
