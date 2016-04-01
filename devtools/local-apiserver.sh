#!/bin/bash

set -e 
CURDIR=$(cd $(dirname $0); pwd -P)
TOPDIR=$(cd $CURDIR/..; pwd -P)
if [[ -z "$GCLOUD_CONFIG_DIR" ]]; then
    GCLOUD_CONFIG_DIR=$TOPDIR/.gcloud
fi

export OCTAV_MYSQL_DBNAME=octav
export OCTAV_MYSQL_USERNAME=root
export OCTAV_MYSQL_ADDRESS_FILE=$GCLOUD_CONFIG_DIR/sql/address
export OCTAV_MYSQL_PASSWORD_FILE=$GCLOUD_CONFIG_DIR/sql/password
export OCTAV_MYSQL_CA_CERT_FILE=$GCLOUD_CONFIG_DIR/sql/server-ca.pem
export OCTAV_MYSQL_CLIENT_CERT_FILE=$GCLOUD_CONFIG_DIR/sql/client-cert.pem
export OCTAV_MYSQL_CLIENT_KEY_FILE=$GCLOUD_CONFIG_DIR/sql/client-key.pem
export OCTAV_TRACE_DB=1

GOVERSION=$(go version)
IFS=' ' read -ra GOVCOMPS <<< "$GOVERSION"
IFS='/' read -ra GOVCOMPS <<< "${GOVCOMPS[3]}"
GOOS=${GOVCOMPS[0]}
GOARCH=${GOVCOMPS[1]}

OCTAV_BIN=octav
if [[ "$DEBUG" ]]; then
    OCTAV_BIN=${OCTAV_BIN}-debug
fi
make -C $TOPDIR/octav $OCTAV_BIN

exec octav/_bin/$GOOS/$GOARCH/$OCTAV_BIN
