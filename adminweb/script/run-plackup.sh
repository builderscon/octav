#!/bin/sh

exec plackup -a /adminweb/app.psgi \
    -s Starlet \
    --max-workers=10 \
    --max-reqs-per-child=500 \
    --min-reqs-per-child=350 
    --spawn-internal=0.5

    