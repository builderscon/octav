#!/bin/bash

set -e

# tags is required to install optional dependencies
go get -t -v -tags debug0 ./...

if [ -z "$OCTAV_DB_NAME" ]; then
    OCTAV_DB_NAME=octav
fi

echo " + Creating database '$OCTAV_DB_NAME'"
mysql -u root -e "CREATE DATABASE $OCTAV_DB_NAME"
mysql -u root octav < sql/octav.sql

export OCTAV_TEST_DSN="root:@/$OCTAV_DB_NAME?parseTime=true"
export OCTAV_TRACE_DB=1

exec go test -v -tags debug0 