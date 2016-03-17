#!/bin/bash

# We go through this hoopla to create the secret so that we don't have to
# commit extra files that otherwise may reveal sensitive information.
#
# The JSON file created this way can be fed into kubectl like so:
#
#    ./devtools/make_cloudsql_secret.sh | kubectl create -f -
#
# Kubernetes site shows that you can do this from the kubectl command line
# alone, but as of this writing at least kubectl that comes with the
# gcloud toolset doesn't, so... this workaround

if [ -z "$CLOUDSQL_DIR" ]; then
    CLOUDSQL_DIR=.gcloud/sql
fi

if [ -z "$CLOUDSQL_SECRET_NAME" ]; then
    CLOUDSQL_SECRET_NAME=cloudsql-dev
fi

which jo >/dev/null 2>&1
if [ $? -ne 0 ]; then
  echo "Missing 'jo' executable. Please install 'jo' to proceed"
  exit 1
fi

exec jo -p \
    kind=Secret \
    apiVersion=v1 \
    metadata=$(jo name=$CLOUDSQL_SECRET_NAME labels=$(jo name=cloudsql group=secrets)) \
    data[address]="$(base64 $CLOUDSQL_DIR/address)" \
    data[password]="$(base64 $CLOUDSQL_DIR/password)" \
    data[server-ca.pem]="$(base64 $CLOUDSQL_DIR/server-ca.pem)" \
    data[client-key.pem]="$(base64 $CLOUDSQL_DIR/client-key.pem)" \
    data[client-cert.pem]="$(base64 $CLOUDSQL_DIR/client-cert.pem)"

