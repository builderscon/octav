GCLOUD_CONFIG_DIR = .gcloud
CONTAINER_DIR=$(CURDIR)/gke/containers
CLOUDSQL_DIR=$(GCLOUD_CONFIG_DIR)/sql
CLOUDSQL_ADDRESS_FILE = $(CLOUDSQL_DIR)/address
CLOUDSQL_ADDRESS = $(shell cat $(CLOUDSQL_ADDRESS_FILE))
CLOUDSQL_SSL_FILES    = \
	$(CLOUDSQL_DIR)/server-ca.pem \
	$(CLOUDSQL_DIR)/client-key.pem \
	$(CLOUDSQL_DIR)/client-cert.pem

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


.PHONY: gke-apiserver gke-adminweb k8-bootstrap k8-shutdown cloudsql initdb

k8s-bootstrap:
	./devtools/k8s-bootstrap.sh

k8s-shutdown:
	./devtools/k8s-shutdown.sh

$(CLOUDSQL_ADDRESS_FILE):
ifeq (,$(wildcard $@))
	@echo "Missing '$@'"
	@echo "CloudSQL address needed. Please ask a builderscon admin to provide it for you"
	@exit 1
endif

$(CLOUDSQL_SSL_FILES):
ifeq (,$(wildcard $@))
	@echo "Missing '$@'"
	@echo "CloudSQL requires SSL certificates to connect. Please ask a builderscon admin to provide one for you"
	@exit 1
endif

cloudsql_files: $(CLOUDSQL_ADDRESS_FILE) $(CLOUDSQL_SSL_FILES)

cloudsql: cloudsql_files
	@echo "Connecting to mysql..."
	mysql -uroot -h $(CLOUDSQL_ADDRESS) -p --ssl-ca=$(GCLOUD_CONFIG_DIR)/sql/server-ca.pem --ssl-cert=$(GCLOUD_CONFIG_DIR)/sql/client-cert.pem --ssl-key=$(GCLOUD_CONFIG_DIR)/sql/client-key.pem

# This rule creates a secrets file so that it can be fed into kubectl.
# We go through this hoopla to create the secret so that we don't have to
# commit extra files that otherwise may reveal sensitive information.
#
# Kubernetes site shows that you can do this from the kubectl command line
# alone, but as of this writing at least kubectl that comes with the
# gcloud toolset doesn't, so... this workaround
cloudsql_secret: cloudsql_files
	./devtools/make_cloudsql_secret.sh | kubectl create -f -

initdb:
	@echo "Initializing database..."
	@echo "  --> dropping old database $(OCTAV_DB_NAME) (if it exists)"
	@mysql -u root -e "DROP DATABASE IF EXISTS $(OCTAV_DB_NAME)"
	@echo "  --> creating new database $(OCTAV_DB_NAME)"
	@mysql -u root -e "CREATE DATABASE $(OCTAV_DB_NAME)"
	@echo "  --> running DDL..."
	@mysql -u root octav < sql/octav.sql

docker-env-ready:
	$(DOCKER_MACHINE_ENV)

gke-apiserver:
	@$(MAKE) -C $(CONTAINER_DIR)/apiserver clean
	@$(MAKE) -C $(CONTAINER_DIR)/apiserver docker DEBUG=1
	@$(MAKE) -C $(CONTAINER_DIR)/apiserver publish-deploy

gke-adminweb:
	@$(MAKE) -C $(CONTAINER_DIR)/adminweb clean
	@$(MAKE) -C $(CONTAINER_DIR)/adminweb docker DEBUG=1
	@$(MAKE) -C $(CONTAINER_DIR)/adminweb publish-deploy

gke-publish:
	@echo "Publishing [ $(IMAGE_NAME) ]"
	@docker tag octav/$(APPNAME) $(IMAGE_NAME)
	@echo " --> pushing $(IMAGE_NAME)..."
	@gcloud docker push $(IMAGE_NAME)

gke-deploy:
	@echo "Deploying $(IMAGE_NAME) via rolling update"
	kubectl rolling-update --update-period=5s --image=$(IMAGE_NAME) $(APPNAME)

