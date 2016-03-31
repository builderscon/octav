#!/bin/bash

set -e

curl -L http://install.perlbrew.pl | bash
source ~/perl5/perlbrew/etc/bashrc

perlbrew use 5.22
perl --version
cpanm --version
cpanm --quiet --installdeps --notest .
perl Build.PL && ./Build && ./Build test

