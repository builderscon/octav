#!perl
use strict;
use feature 'say';
use JSON;
use File::Basename;
use File::Path;
use File::Spec;

my $sfile = shift @ARGV;
open (my $fh, '<', $sfile) or die "Failed to open $sfile: $!";

my $sbody = do { local $/; <$fh> };

my $schema = JSON::decode_json($sbody);

my $buf = '';
open(my $tmpout, '>', \$buf);

say $tmpout <<EOM;
"""OCTAV Client Library"""
"""DO NOT EDIT: This file was generated from $sfile on @{[scalar localtime]}"""
EOM

say $tmpout <<'EOM';
import json
import os
import requests

class Octav(object):
  def __init__(self, endpoint=None, key=None, secret=None, debug=False):
    if not endpoint:
      raise "endpoint is required"
    if not key:
      raise "key is required"
    if not secret:
      raise "secret is required"
    self.debug = debug
    self.endpoint = endpoint
    self.error = None
    self.key = key
    self.secret = secret
    self.session = requests.Session()
    self.session.mount('http://', requests.adapters.HTTPAdapter(max_retries=0))

  def extract_error(self, r):
    try:
      js = r.json()
      self.error = js["message"]
    except:
      self.error = r.status_code

  def last_error(self):
    return self.error
EOM

for my $link (@{$schema->{links}}) {
    my $name = camelize_title($link->{title});

    my $path = $link->{href};
    my $link_schema = $link->{schema};
    print $tmpout "  def $name (self";
    my $props = $link_schema->{properties} || {};
    if (my @keys = keys %$props) {
        print $tmpout ", ";
        print $tmpout join (", ", map { "$_=None" } @keys);
    }
    print $tmpout "):\n";

    my %required;
    say $tmpout '    payload = {}';
    if (my $link_schema = $link->{schema}) {
        my $required = $link_schema->{required};
        if ($required && scalar(@{$required}) > 0) {
            foreach my $name (sort @$required) {
                $required{$name}++;
                say $tmpout '    if ' . $name . ' is None:';
                say $tmpout '            raise "property \"" + required + "\" must be provided"';
                say $tmpout '    payload[' . $name . '] = ' . $name;
            }
        }
    }
    if (my @keys = keys %$props) {
        foreach my $key (sort @keys) {
            next if $required{$key};
            say $tmpout "    if not $key is None:";
            say $tmpout "        payload['" . $key . "'] = " . $key;
        }
    }

    say $tmpout '    uri = self.endpoint + "' . $path . '"';
    if (lc($link->{method}) eq 'post') {
        say $tmpout '    if self.debug:';
        say $tmpout '        print("POST " + uri)';
        say $tmpout '    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)';
    } else {
        say $tmpout '    if self.debug:';
        say $tmpout '        print("GET " + uri)';
        say $tmpout '    res = self.session.get(uri, params=payload)';
    }
    say $tmpout '    if res.status_code != 200:';
    say $tmpout '        self.extract_error(res)';
    say $tmpout '        return None';
    if ($link->{targetSchema}) {
        say $tmpout '    return res.json()'
    } else {
        say $tmpout '    return True';
    }
    say $tmpout '';
}

print $buf;

sub camelize_title {
    my $s = shift;
    $s =~ s/[\W+]([\w])/_$1/g;
    return lc($s);
}