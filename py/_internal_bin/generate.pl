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
import urllib3

if os.getenv('SERVER_SOFTWARE', '').startswith('Google App Engine/') or os.getenv('SERVER_SOFTWARE', '').startswith('Development/'):
    from urllib3.contrib.appengine import AppEngineManager as PoolManager
else:
    from urllib3 import PoolManager

import sys
if sys.version[0] == "3":
    from urllib.parse import urlencode
else:
    from urllib import urlencode

class Octav(object):
  def __init__(self, endpoint, key, secret, debug=False):
    if not endpoint:
      raise "endpoint is required"
    if not key:
      raise "key is required"
    if not secret:
      raise "secret is required"
    self.debug = debug
    self.endpoint = endpoint
    self.error = None
    self.http = PoolManager()
    self.key = key
    self.secret = secret

  def extract_error(self, r):
    try:
      js = r.json()
      self.error = js["message"]
    except:
      self.error = r.status

  def last_error(self):
    return self.error
EOM

for my $link (@{$schema->{links}}) {
    my $name = camelize_title($link->{title});

    my $path = $link->{href};
    my $link_schema = $link->{schema};
    my $required = $link_schema && ($link_schema->{required} || []);
    my $props = $link_schema->{properties} || {};
    my @keys = keys %$props;
    my @args;

    # First, find out the required arguments
    my %required;
    foreach my $name (sort @$required) {
        $required{$name}++;
        push @args, $name; # required arguments are as-is
    }
    foreach my $key (sort @keys) {
        next if $required{$key};
        push @args, "$key=None";
    }
    print $tmpout "  def $name (self";
    if (@args) {
        print $tmpout ", ", join(", ", @args);
    }
    print $tmpout "):\n";

    say $tmpout '    try:';
    say $tmpout '        payload = {}';
    say $tmpout '        hdrs = {}';
    foreach my $name (sort @$required) {
        say $tmpout "        if $name is None:";
        say $tmpout "            raise 'property $name must be provided'";
        say $tmpout "        payload['" . $name . "'] = " . $name;
    }

    foreach my $key (sort @keys) {
        say $tmpout "        if $key is not None:";
        say $tmpout "            payload['" . $key . "'] = " . $key;
    }
    say $tmpout q|        uri = '%s|, $path, q|' % self.endpoint|;
    if ($link->{'hsup.wrapper'} eq 'httpWithBasicAuth') {
        say $tmpout q|        hdrs = urllib3.util.make_headers(|;
        say $tmpout q|            basic_auth='%s:%s' % (self.key, self.secret),|;
        say $tmpout q|        )|;
    }
    if (lc($link->{method}) eq 'post') {
        say $tmpout q|        if self.debug:|;
        say $tmpout q|            print('POST %s' % uri)|;
        say $tmpout q|        hdrs['Content-Type']= 'application/json'|;
        say $tmpout q|        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))|;
    } else {
        say $tmpout q|        qs = urlencode(payload)|;
        say $tmpout q|        if self.debug:|;
        say $tmpout q|            print('GET %s?%s' % (uri, qs))|;
        say $tmpout q|        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)|;
    }

    say $tmpout '        if self.debug:';
    say $tmpout '            print(res)';
    say $tmpout '        if res.status != 200:';
    say $tmpout '            self.extract_error(res)';
    say $tmpout '            return None';
    if ($link->{targetSchema}) {
        say $tmpout '        return json.loads(res.data)';
    } else {
        say $tmpout '        return True';
    }
    say $tmpout '    except BaseException, e:';
    say $tmpout '        if self.debug:';
    say $tmpout '            print("error during http access: " + repr(e))';
    say $tmpout '        self.error = repr(e)';
    say $tmpout '        return None';
    say $tmpout '';
}

print $buf;

sub camelize_title {
    my $s = shift;
    $s =~ s/[\W+]([\w])/_$1/g;
    return lc($s);
}