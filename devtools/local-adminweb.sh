#!/bin/bash

CURDIR=$(cd $(dirname $0); pwd -P)
TOPDIR=$(cd ${CURDIR}/..; pwd -P)
ADMINWEB_DIR=${TOPDIR}/adminweb

if [[ -z "$OCTAV_GOOGLE_MAPS_API_KEY" ]]; then
    OCTAV_GOOGLE_MAPS_API_KEY=${TOPDIR}/.googlemaps/apikey
fi

if [[ -z "$OCTAV_GITHUB_CLIENT_ID" ]]; then
    OCTAV_GITHUB_CLIENT_ID=${TOPDIR}/.github/id
fi

if [[ -z "$OCTAV_GITHUB_CLIENT_SECRET" ]]; then
    OCTAV_GITHUB_CLIENT_SECRET=${TOPDIR}/.github/secret
fi

if [[ -z "$OCTAV_REDIS" ]]; then
    OCTAV_REDIS=127.0.0.1:6379
fi

export OCTAV_GOOGLE_MAPS_API_KEY
export OCTAV_GITHUB_CLIENT_ID
export OCTAV_GITHUB_CLIENT_SECRET
export OCTAV_REDIS

for file in $OCTAV_GOOGLE_MAPS_API_KEY $OCTAV_GITHUB_CLIENT_ID $OCTAV_GITHUB_CLIENT_SECRET; do
    if [[ ! -e "$file" ]]; then
        echo "Required file '$file' not found"
    fi
done

echo
echo "Going to run plackup via carton..."
echo "Do NOT forget to change /etc/hosts to point admin.builderscon.io to 127.0.0.1"
echo

# export MOJO_MODE=production
# export MOJO_LOG_LEVEL=debug
cd $ADMINWEB_DIR
exec carton exec -- \
    plackup -I../p5/lib -Ilib -s Starlet