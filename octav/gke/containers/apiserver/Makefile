APPNAME=apiserver
OCTAV_DIR = $(CURDIR)/../../..
REPOSITORY_HOST=asia.gcr.io
PROJECT_ID=$(shell gcloud config list 2>&1 | grep 'project = ' | sed -e 's/project = //')

TAG:=$(shell date +"%Y%m%d.%H%M%S")
IMAGE_NAME:=$(REPOSITORY_HOST)/$(PROJECT_ID)/$(APPNAME):$(TAG)

DOCKER_MACHINE_ENV := @echo "Docker environment set"
DOCKER_MACHINE_EV1 := $(DOCKER_TLS_VERIFY)
DOCKER_MACHINE_EV2 := $(DOCKER_HOST)
DOCKER_MACHINE_EV3 := $(DOCKER_CERT_PATH)
DOCKER_MACHINE_EV4 := $(DOCKER_MACHINE_NAME)
ifeq ($(strip $(DOCKER_MACHINE_EV1)),)
	DOCKER_MACHINE_ENV = @echo "\n Docker environment missing, run: eval \"\$$(docker-machine env default)\"\n" && exit 1
endif
ifeq ($(strip $(DOCKER_MACHINE_EV2)),)
	DOCKER_MACHINE_ENV = @echo "\n Docker environment missing, run: eval \"\$$(docker-machine env default)\"\n" && exit 1
endif
ifeq ($(strip $(DOCKER_MACHINE_EV3)),)
	DOCKER_MACHINE_ENV = @echo "\n Docker environment missing, run: eval \"\$$(docker-machine env default)\"\n" && exit 1
endif
ifeq ($(strip $(DOCKER_MACHINE_EV4)),)
	DOCKER_MACHINE_ENV = @echo "\n Docker environment missing, run: eval \"\$$(docker-machine env default)\"\n" && exit 1
endif

.PHONY: all octav

octav: 
	@cd $(OCTAV_DIR) ; make build GOOS=linux GOARCH=amd64; cd $(CURDIR)
	@cp $(OCTAV_DIR)/_bin/octav .

dm-running:
ifneq (,$(strip $(shell docker-machine env default 2>&1 | grep -o 'not running')))
	@echo "Starting docker-machine (default)"
	@docker-machine start default
endif

docker: dm-running 
	$(DOCKER_MACHINE_ENV)
	@docker build -t octav/$(APPNAME) .

publish:
	@echo "Publishing [ $(IMAGE_NAME) ]"
	docker tag octav/$(APPNAME) $(IMAGE_NAME)
	gcloud docker push $(IMAGE_NAME)

clean:
	@echo "Deleting old images"
	@-docker images -q --filter dangling=true | xargs docker rmi

deploy:
	@echo "Deploying $(IMAGE_NAME) via rolling update"
	@kubectl rolling-update --update-period=5s --image=$(IMAGE_NAME) $(APPNAME)
