#!/bin/bash

# We go through this hoopla to create the secret so that we don't have to
# commit extra files that otherwise may reveal sensitive information.
#
# The JSON file created this way can be fed into kubectl like so:
#
#    ./devtools/make_slack_secret.sh | kubectl create -f -
#
# Kubernetes site shows that you can do this from the kubectl command line
# alone, but as of this writing at least kubectl that comes with the
# gcloud toolset doesn't, so... this workaround

if [ -z "$SLACK_DIR" ]; then
    SLACK_DIR=.slack/
fi

if [ -z "$SLACK_SECRET_NAME" ]; then
    SLACK_SECRET_NAME=slack-dev
fi

which jo >/dev/null 2>&1
if [ $? -ne 0 ]; then
  echo "Missing 'jo' executable. Please install 'jo' to proceed"
  exit 1
fi

exec jo -p \
    kind=Secret \
    apiVersion=v1 \
    metadata=$(jo name=$SLACK_SECRET_NAME labels=$(jo name=slack group=secrets)) \
    data[token]="$(base64 $SLACK_DIR/token)" \
    data[gwtoken]="$(base64 $SLACK_DIR/gwtoken)"
