APPNAME=h2o
BUILDER_IMAGE=octav/h2o-static-builder
TOPDIR=$(CURDIR)/../../..
REPOSITORY_HOST=asia.gcr.io
PROJECT_ID=$(shell gcloud config list 2>&1 | grep 'project = ' | sed -e 's/project = //')

TAG:=$(shell date +"%Y%m%d.%H%M%S")
IMAGE_NAME:=$(REPOSITORY_HOST)/$(PROJECT_ID)/$(APPNAME):$(TAG)


build-h2o-binary:
	docker build -f Dockerfile.build -t octav/h2o-static-builder .

h2o:
	$(MAKE) build-h2o-binary
	docker run --cidfile h2o-static-builder.cid $(BUILDER_IMAGE)
	docker cp `cat h2o-static-builder.cid`:/h2o/h2o .
	@-docker rm `cat h2o-static-builder.cid`
	@-rm h2o-static-builder.cid

docker: h2o
	docker build -t octav/h2o .

publish:
	$(MAKE) -C $(TOPDIR) gke-publish IMAGE_NAME=$(IMAGE_NAME) APPNAME=$(APPNAME)

clean:
	@-rm h2o
