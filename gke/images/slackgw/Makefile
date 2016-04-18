APPNAME=slackgw
TOPDIR=$(CURDIR)/../../..
REPOSITORY_HOST=asia.gcr.io
PROJECT_ID=$(shell gcloud config list 2>&1 | grep 'project = ' | sed -e 's/project = //')

TAG:=$(shell date +"%Y%m%d.%H%M%S")
IMAGE_NAME:=$(REPOSITORY_HOST)/$(PROJECT_ID)/$(APPNAME):$(TAG)

$(CURDIR)/slackgw-debug:
	GOOS=linux GOARCH=amd64 go build -tags debug0 -o slackgw-debug \
		$(shell go list -f '{{ .Dir }}' github.com/lestrrat/go-slackgw)/cmd/slackgw/slackgw.go

$(CURDIR)/slackgw:
	GOOS=linux GOARCH=amd64 go build -o slackgw-debug \
		$(shell go list -f '{{ .Dir }}' github.com/lestrrat/go-slackgw)/cmd/slackgw/slackgw.go

slackgw-debug: $(CURDIR)/slackgw-debug
slackgw: $(CURDIR)/slackgw

docker-debug: $(CURDIR)/slackgw-debug
	cp $(CURDIR)/slackgw-debug $(CURDIR)/slackgw
	$(MAKE) docker

docker: $(CURDIR)/slackgw
	docker build -t octav/slackgw .

publish:
	$(MAKE) -C $(TOPDIR) gke-publish IMAGE_NAME=$(IMAGE_NAME) APPNAME=$(APPNAME)

clean:
	@-rm slackgw
	@-rm slackgw-debug
