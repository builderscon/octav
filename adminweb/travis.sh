#!/bin/bash

set -e

echo "====> perl version"
perl --version
echo "====> cpanm version"
cpanm --version
echo "====> Installing dependencies..."
cpanm -v --installdeps --notest .
echo "====> Running tests..."
perl Build.PL && ./Build && ./Build test

