#!/bin/bash

set -e

go env
go get -t -v ./...

echo " + Creating database '$OCTAV_DB_NAME'"
mysql -u root -e "CREATE DATABASE $OCTAV_DB_NAME"
mysql -u root octav < sql/octav.sql

export OCTAV_TEST_DSN="root:@/$OCTAV_DB_NAME?parseTime=true"
export OCTAV_TRACE_DB=1

exec go test -v -tags debug0 