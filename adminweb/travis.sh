#!/bin/bash

set -e

curl -L http://install.perlbrew.pl | bash
source ~/perl5/perlbrew/etc/bashrc

perlbrew install perl-5.22.1
perlbrew use perl-5.22.1
perl --version
cpanm --version
cpanm --quiet --installdeps --notest .
perl Build.PL && ./Build && ./Build test

