APPNAME=apiserver
OCTAV_DIR = $(CURDIR)/../../..
REPOSITORY_HOST=asia.gcr.io
PROJECT_ID=$(shell gcloud config list 2>&1 | grep 'project = ' | sed -e 's/project = //')

# Unfortunately this had to be replicated from top Makefile
SRC_FILES = $(wildcard $(OCTAV_DIR)/*.go $(OCTAV_DIR)/validator/*.go $(OCTAV_DIR)/client/*.go $(OCTAV_DIR)/db/*.go $(OCTAV_DIR)/model/*.go)

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

.PHONY: all octav always-rebuild

$(CURDIR)/octav: $(SRC_FILES)
	@echo "Building octav"
	cd $(OCTAV_DIR) ; make GOOS=linux GOARCH=amd64 octav; cd $(CURDIR)
	@cp $(OCTAV_DIR)/_bin/linux/amd64/octav .

$(CURDIR)/octav-debug: $(SRC_FILES)
	@echo "Building octav-debug"
	cd $(OCTAV_DIR) ; make GOOS=linux GOARCH=amd64 octav-debug; cd $(CURDIR)
	@cp $(OCTAV_DIR)/_bin/linux/amd64/octav-debug .

build-octav: $(CURDIR)/octav
build-octav-debug: $(CURDIR)/octav-debug

dm-running:
ifneq (,$(strip $(shell docker-machine env default 2>&1 | grep -o 'not running')))
	@echo "Starting docker-machine (default)"
	@docker-machine start default
endif



docker-debug: dm-running build-octav-debug
	$(DOCKER_MACHINE_ENV)
	cp octav-debug octav
	@docker build -t octav/$(APPNAME) .

docker: dm-running build-octav
	$(DOCKER_MACHINE_ENV)
	@docker build -t octav/$(APPNAME) .

publish:
	@echo "Publishing [ $(IMAGE_NAME) ]"
	docker tag octav/$(APPNAME) $(IMAGE_NAME)
	gcloud docker push $(IMAGE_NAME)

clean_images:
	@echo "Deleting old images"
	@-docker images -q --filter dangling=true | xargs docker rmi

clean_executables:
	@echo "Deleting generated executables"
	@-rm octav

clean: clean_images clean_executables

deploy:
	@echo "Deploying $(IMAGE_NAME) via rolling update"
	kubectl rolling-update --update-period=5s --image=$(IMAGE_NAME) $(APPNAME)
