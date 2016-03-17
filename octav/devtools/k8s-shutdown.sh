#!/bin/bash

set -e

if [ -z "$MY_CONFIG_NAME" ]; then
    MY_CONFIG_NAME=octav-dev
fi

# Name of the Container Cluster
if [ -z "$CLUSTER_NAME" ]; then
    CLUSTER_NAME="${MY_CONFIG_NAME}-cluster-1"
fi

# gcloud needs to know which configuration set we're dealing with
echo "* Activating $MY_CONFIG_NAME configuration..."
OLDCONFIG=$(gcloud config configurations list | grep True | awk '{print $1}')
gcloud config configurations activate $MY_CONFIG_NAME > /dev/null
# restore old configuration
trap "gcloud config configurations activate $OLDCONFIG > /dev/null" EXIT

INSTANCE_GROUP=$(basename $(gcloud container clusters describe $CLUSTER_NAME | grep gke | awk '{print $2}'))

for service in $(gcloud compute backend-services list | grep $INSTANCE_GROUP | awk '{print $1}'); do
    gcloud compute backend-services remove-backend --instance-group=$INSTANCE_GROUP $service
done

yes | gcloud container clusters delete $CLUSTER_NAME