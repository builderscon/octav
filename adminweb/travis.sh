#!/bin/bash

set -e

perl --version
cpanm --version
cpanm --quiet --installdeps --notest .
perl Build.PL && ./Build && ./Build test

