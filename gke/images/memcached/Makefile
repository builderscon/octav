APPNAME=memcached
BUILDER_IMAGE=octav/memcached-static-builder
TOPDIR=$(CURDIR)/../../..
REPOSITORY_HOST=asia.gcr.io
PROJECT_ID=$(shell gcloud config list 2>&1 | grep 'project = ' | sed -e 's/project = //')

TAG:=$(shell date +"%Y%m%d.%H%M%S")
IMAGE_NAME:=$(REPOSITORY_HOST)/$(PROJECT_ID)/$(APPNAME):$(TAG)


build-memcached-binary:
	docker build -f Dockerfile.build -t octav/memcached-static-builder .

memcached:
	$(MAKE) build-memcached-binary
	docker run --cidfile memcached-static-builder.cid $(BUILDER_IMAGE)
	docker cp `cat memcached-static-builder.cid`:/memcached/memcached .
	@-docker rm `cat memcached-static-builder.cid`
	@-rm memcached-static-builder.cid

docker: memcached
	docker build -t octav/memcached .

publish:
	$(MAKE) -C $(TOPDIR) gke-publish IMAGE_NAME=$(IMAGE_NAME) APPNAME=$(APPNAME)

clean:
	@-rm memcached
