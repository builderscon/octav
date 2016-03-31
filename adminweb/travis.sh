#!/bin/bash

set -e

echo "====> perl version"
perl --version
echo "====> cpanm version"
cpanm --version
echo "====> Installing dependencies..."
cpanm -v --installdeps --notest -Llocal .
echo "====> Running tests..."
export PERL5OPT=-Mlib=local/lib/perl5
perl Build.PL && ./Build && ./Build test

