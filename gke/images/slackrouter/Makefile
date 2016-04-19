APPNAME=slackrouter
TOPDIR=$(CURDIR)/../../..
REPOSITORY_HOST=asia.gcr.io
PROJECT_ID=$(shell gcloud config list 2>&1 | grep 'project = ' | sed -e 's/project = //')

TAG:=$(shell date +"%Y%m%d.%H%M%S")
IMAGE_NAME:=$(REPOSITORY_HOST)/$(PROJECT_ID)/$(APPNAME):$(TAG)

$(CURDIR)/slackrouter-debug:
	cd $(TOPDIR)/worker/slackrouter && \
	GOOS=linux GOARCH=amd64 go build -tags debug0 -o $(CURDIR)/slackrouter-debug .

$(CURDIR)/slackrouter:
	cd $(TOPDIR)/worker/slackrouter && \
	GOOS=linux GOARCH=amd64 go build -o $(CURDIR)/slackrouter .

slackrouter-debug: $(CURDIR)/slackrouter-debug
slackrouter: $(CURDIR)/slackrouter

docker-debug: $(CURDIR)/slackrouter-debug
	cp $(CURDIR)/slackrouter-debug $(CURDIR)/slackrouter
	$(MAKE) docker

docker: $(CURDIR)/slackrouter
	docker build -t octav/slackrouter .

publish:
	$(MAKE) -C $(TOPDIR) gke-publish IMAGE_NAME=$(IMAGE_NAME) APPNAME=$(APPNAME)

clean:
	@-rm slackrouter
	@-rm slackrouter-debug
