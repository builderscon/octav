#!/bin/bash
set -e

# This script creates a new GKE environment for octav. It should not
# be used to "upgrade" an existing cluster, as the basic way this script
# handles collisions and pre-existing resources is to delete it and
# then to re-create it. You should only use it create a new environment,
# or to anihilate an existing environment.
#
# See also: k8s-shutdown.sh

if [ -z "$MY_CONFIG_NAME" ]; then
    MY_CONFIG_NAME=octav-dev
fi

# Name of the Container Cluster
if [ -z "$CLUSTER_NAME" ]; then
    CLUSTER_NAME="${MY_CONFIG_NAME}-cluster-1"
fi

if [ -z "$CLUSTER_NODE_SIZE" ]; then
    CLUSTER_NODE_SIZE=3
fi

if [ -z "$CLUSTER_MACHINE_TYPE" ]; then
    CLUSTER_MACHINE_TYPE=g1-small
fi

SECRETS="cloudsql logging"
SERVICES=apiserver
REPLICATION_CONTROLLERS=apiserver adminweb

# gcloud needs to know which configuration set we're dealing with
echo "* Activating $MY_CONFIG_NAME configuration..."
OLDCONFIG=$(gcloud config configurations list | grep True | awk '{print $1}')
gcloud config configurations activate $MY_CONFIG_NAME > /dev/null
# restore old configuration
trap "gcloud config configurations activate $OLDCONFIG > /dev/null" EXIT

echo "* Checking if we have cluster '$CLUSTER_NAME'..."
set +e
RETVAL=$(gcloud container clusters describe $CLUSTER_NAME > /dev/null 2>&1; echo $?)
set -e

if [ $RETVAL == 0 ]; then
    echo "--> Cluster '$CLUSTER_NAME' already exists."
else
    echo "--> Cluster '$CLUSTER_NAME' not found. Creating cluster..."
    gcloud container clusters create $CLUSTER_NAME \
        --num-nodes=$CLUSTER_NODE_SIZE \
        --machine-type=$CLUSTER_MACHINE_TYPE
    CREATED_CLUSTER=1
fi

INSTANCE_GROUP=$(basename $(gcloud container clusters describe $CLUSTER_NAME | grep gke | awk '{print $2}'))
INSTANCE_TAG=$(echo $INSTANCE_GROUP | sed -e s/group/node/)
cat <<EOM
--> Cluster information:
  Group: $INSTANCE_GROUP
  Tag:   $INSTANCE_TAG
EOM


echo "* Creating replication controller(s)..."
for name in $REPLICATION_CONTROLLERS; do
    echo "Processing replication controller $name ..."

    set +e
    RETVAL=$(kubectl get rc $name 2>&1 > /dev/null ; echo $?)
    set -e
    if [ "$RETVAL" == "0" ]; then 
        echo " - Deleting old replication controller $name ..."
        kubectl delete rc -l name=$name
    fi
    echo " + Creating replication controller $name ..."
    kubectl create -f gke/rc/$name.yaml
done

echo "* Creating service(s)..."
for name in $SERVICES; do
    echo "Processing service $name ..."

    set +e
    RETVAL=$(kubectl get service $name >/dev/null 2>&1; echo $?)
    set -e
    if [ "$RETVAL" == "0" ]; then
        echo " - Deleting old service $name..."
        kubectl delete service $name
    fi

    echo " + Creating service $name..."
    kubectl create -f gke/service/$name.yaml
done

