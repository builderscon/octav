APPNAME=deploybot-k8s
TOPDIR=$(CURDIR)/../../..
REPOSITORY_HOST=asia.gcr.io
PROJECT_ID=$(shell gcloud config list 2>&1 | grep 'project = ' | sed -e 's/project = //')

TAG:=$(shell date +"%Y%m%d.%H%M%S")
IMAGE_NAME_BASE:=$(REPOSITORY_HOST)/$(PROJECT_ID)/$(APPNAME)
IMAGE_NAME_LATEST:=$(IMAGE_NAME_BASE):latest
IMAGE_NAME:=$(IMAGE_NAME_BASE):$(TAG)

$(CURDIR)/deploybot-k8s-debug:
	cd $(TOPDIR)/worker/deploybot-k8s && \
	GOOS=linux GOARCH=amd64 go build -tags debug0 -o $(CURDIR)/deploybot-k8s-debug .

$(CURDIR)/deploybot-k8s:
	cd $(TOPDIR)/worker/deploybot-k8s && \
	GOOS=linux GOARCH=amd64 go build -o $(CURDIR)/deploybot-k8s .

deploybot-k8s-debug: $(CURDIR)/deploybot-k8s-debug
deploybot-k8s: $(CURDIR)/deploybot-k8s

docker-debug: $(CURDIR)/deploybot-k8s-debug
	cp $(CURDIR)/deploybot-k8s-debug $(CURDIR)/deploybot-k8s
	$(MAKE) docker

docker: $(CURDIR)/deploybot-k8s
	docker build -t octav/deploybot-k8s .

publish:
	$(MAKE) -C $(TOPDIR) gke-publish \
		IMAGE_NAME_LATEST=$(IMAGE_NAME_LATEST) \
		IMAGE_NAME=$(IMAGE_NAME) \
		APPNAME=$(APPNAME)

clean:
	@-rm deploybot-k8s
	@-rm deploybot-k8s-debug
