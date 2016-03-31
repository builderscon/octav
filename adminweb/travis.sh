#!/bin/bash

set -e

echo "====> perl version"
perl --version

export PATH=local/bin:$PATH
mkdir -p local/bin
curl -L http://cpanmin.us > local/bin/cpanm
chmod +x local/bin/cpanm

echo "====> cpanm version"
cpanm --version
echo "====> Installing dependencies..."
cpanm --quiet --installdeps --notest -Llocal .
echo "====> Running tests..."
export PERL5OPT="-Mlib=local/lib/perl5 -Mlib=../p5/lib"
perl Build.PL && ./Build && ./Build test

