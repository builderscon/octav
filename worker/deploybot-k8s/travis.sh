#!/bin/bash

set -e

export GO_TAGS_OPT="-tags debug0"

if [ "$TRAVIS" == "true" ]; then
    echo " + Detected running under Travis CI"
    make glide
    make installdeps
fi

export OCTAV_DEBUG_FILE=/tmp/debug.out

exec make test
