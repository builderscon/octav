#!/bin/bash

set -e

if [ -z "$OCTAV_DB_NAME" ]; then
    OCTAV_DB_NAME=octav
fi

if [ "$TRAVIS" == "true" ]; then
    echo " + Detected running under Travis CI"
    echo " + Creating database '$OCTAV_DB_NAME'"
    mysql -u root -e "CREATE DATABASE $OCTAV_DB_NAME"
    mysql -u root octav < sql/octav.sql
    go get -t -v -tags debug0 ./...
fi

export OCTAV_TEST_DSN="root:@/$OCTAV_DB_NAME?parseTime=true"
export OCTAV_TRACE_DB=1
export OCTAV_DEBUG_FILE=debug.out

exec go test -v -tags debug0 