APPNAME=acmebot-k8s-certupload
TOPDIR=$(CURDIR)/../../..
REPOSITORY_HOST=asia.gcr.io
PROJECT_ID=$(shell gcloud config list 2>&1 | grep 'project = ' | sed -e 's/project = //')

TAG:=$(shell date +"%Y%m%d.%H%M%S")
IMAGE_NAME:=$(REPOSITORY_HOST)/$(PROJECT_ID)/$(APPNAME):$(TAG)

$(CURDIR)/acmebot-k8s-certupload-debug:
	cd $(TOPDIR)/worker/acmebot-k8s-certupload && \
	GOOS=linux GOARCH=amd64 go build -tags debug0 -o $(CURDIR)/acmebot-k8s-certupload-debug .

$(CURDIR)/acmebot-k8s-certupload:
	cd $(TOPDIR)/worker/acmebot-k8s-certupload && \
	GOOS=linux GOARCH=amd64 go build -o $(CURDIR)/acmebot-k8s-certupload .

acmebot-k8s-certupload-debug: $(CURDIR)/acmebot-k8s-certupload-debug
acmebot-k8s-certupload: $(CURDIR)/acmebot-k8s-certupload

docker-debug: $(CURDIR)/acmebot-k8s-certupload-debug
	cp $(CURDIR)/acmebot-k8s-certupload-debug $(CURDIR)/acmebot-k8s-certupload
	$(MAKE) docker

docker: $(CURDIR)/acmebot-k8s-certupload
	docker build -t octav/acmebot-k8s-certupload .

publish:
	$(MAKE) -C $(TOPDIR) gke-publish IMAGE_NAME=$(IMAGE_NAME) APPNAME=$(APPNAME)

clean:
	@-rm acmebot-k8s-certupload
	@-rm acmebot-k8s-certupload-debug
