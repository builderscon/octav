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
import certifi
import feedparser
import functools
import json
import os
import re
import time
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

class MissingRequiredArgument(Exception):
    pass

def parse_rfc3339(s):
    return time.mktime(feedparser._parse_date(s))

EOM

my %classes;
for my $link (@{$schema->{links}}) {
    my $name = camelize_title($link->{title});

    my $path = $link->{href};
    my $link_schema = $link->{schema};
    my $required = ($link_schema && $link_schema->{required}) || [];
    my $props = $link_schema->{properties} || {};
    my $patternProperties = $link_schema->{patternProperties} || {};

    # This changes which object the method belongs to
    my %methods = (
        'Octav' => [],
        'Session' => [],
    );
    my $class = 'Octav';
    if (my $w = $link->{'hsup.wrapper'}) {
        if (ref $w ne 'ARRAY') {
            $link->{'hsup.wrapper'} = [$w];
        }
    }
    for my $w (@{$link->{'hsup.wrapper'}}) {
        if ($w eq 'httpWithClientSession') {
            $class = 'Session';
        }
    }
    push @{$methods{$class}}, $name;

    my $octav_method_out;
    my $octav_method_body;
    if (!open($octav_method_out, '>', \$octav_method_body)) {
        die $!;
    }

    my $session_method_out;
    my $session_method_body;
    if (!open($session_method_out, '>', \$session_method_body)) {
        die $!;
    }

    my @keys = keys %$props;
    my @args;

    # First, find out the required arguments
    my %required;
    my @proxy_args;
    foreach my $name (sort @$required) {
        $required{$name}++;
        push @args, $name; # required arguments are as-is
        push @proxy_args, $name;
    }
    foreach my $key (sort @keys) {
        next if $required{$key};
        push @args, "$key=None";
        push @proxy_args, "$key=$key";
    }
    push @args, 'extra_headers=None';
    push @proxy_args, q|extra_headers={'X-Octav-Session-ID': self.sid}|;
    if (keys %$patternProperties > 0) {
        push @args, '**args';
        push @proxy_args, '**args';
    }

    print_method_signature($octav_method_out, $name, @args);
    print_method_signature($session_method_out, $name, @args);

    # Session methods are just wrappers into Octav
    print $session_method_out <<EOM;
    self.renew()
    return self.client.$name(@{[join ", ", @proxy_args]})
EOM

    say $octav_method_out '    try:';
    say $octav_method_out '        payload = {}';
    say $octav_method_out '        hdrs = {}';
    foreach my $name (sort @$required) {
        say $octav_method_out "        if $name is None:";
        say $octav_method_out "            raise MissingRequiredArgument('property $name must be provided')";
        say $octav_method_out "        payload['" . $name . "'] = " . $name;
    }

    foreach my $key (sort @keys) {
        say $octav_method_out "        if $key is not None:";
        say $octav_method_out "            payload['" . $key . "'] = " . $key;
    }
    if (keys(%$patternProperties) > 0) {
        my @patternKeys;
        foreach my $key (sort keys %$patternProperties) {
            $key =~ s/'/\\'/g;
            push @patternKeys, $key;
        }
        say $octav_method_out "        patterns = [", (join(", ", map { "re.compile('$_')" } @patternKeys)), "]";
        say $octav_method_out "        for key in args:";
        say $octav_method_out "            for p in patterns:";
        say $octav_method_out "                if p.match(key):";
        say $octav_method_out "                    payload[key] = args[key]";
    }
    say $octav_method_out q|        uri = '%s|, $path, q|' % self.endpoint|;
    my $basic_auth = 0;
    for my $w (@{$link->{'hsup.wrapper'}}) {
        if ($w =~ /^httpWith(Optional)?BasicAuth/) {
            $basic_auth = 1;
            last;
        }
    }
    if ($basic_auth) {
        say $octav_method_out q|        hdrs = urllib3.util.make_headers(|;
        say $octav_method_out q|            basic_auth='%s:%s' % (self.key, self.secret),|;
        say $octav_method_out q|        )|;
    }
    if (lc($link->{method}) eq 'post') {
        say $octav_method_out q|        if self.debug:|;
        say $octav_method_out q|            print('POST %s' % uri)|;
        say $octav_method_out q|        hdrs['Content-Type']= 'application/json'|;
        say $octav_method_out q|        if extra_headers:|;
        say $octav_method_out q|            hdrs.update(extra_headers)|;
        say $octav_method_out q|        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))|;
    } else {
        say $octav_method_out q|        qs = urlencode(payload, True)|;
        say $octav_method_out q|        if self.debug:|;
        say $octav_method_out q|            print('GET %s?%s' % (uri, qs))|;
        say $octav_method_out q|        if extra_headers:|;
        say $octav_method_out q|            hdrs.update(extra_headers)|;
        say $octav_method_out q|        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)|;
    }

    say $octav_method_out '        if self.debug:';
    say $octav_method_out '            print(res)';
    say $octav_method_out '        self.res = res';
    say $octav_method_out '        if res.status != 200:';
    say $octav_method_out '            self.extract_error(res)';
    say $octav_method_out '            return None';
    if ($link->{targetSchema}) {
        say $octav_method_out '        return json.loads(res.data)';
    } else {
        say $octav_method_out '        return True';
    }
    say $octav_method_out '    except BaseException as e:';
    say $octav_method_out '        if self.debug:';
    say $octav_method_out '            print("error during http access: " + repr(e))';
    say $octav_method_out '        self.error = repr(e)';
    say $octav_method_out '        return None';
    say $octav_method_out '';

    push @{$classes{'Octav'}->{$name}}, $octav_method_body;
    push @{$classes{'Session'}->{$name}}, $session_method_body;
}

print $tmpout <<'EOM';
class Session(object):
  def __init__(self, client, f, sid, expires):
    if not f:
      raise MissingRequiredArgument('f is required')
    if not client:
      raise MissingRequiredArgument('client is required')
    if not sid:
      raise MissingRequiredArgument('sid is required')
    if not expires:
      raise MissingRequiredArgument('expires is required')

    self.update_func = f
    self.sid = sid
    self.expires = expires
    self.client = client

  def last_error(self):
    return self.client.last_error()

  # renews the octav session. returns false if there was no need to
  # renew, true if the session was renewed. None is returned if
  # there was an error
  def renew(self):
    if self.expires > time.mktime(time.gmtime()):
        return False
    s = self.update_func()
    if s is None:
        return None
    self.sid = s.get('sid')
    self.expires = parse_rfc3339(s.get('expires'))
    return True

EOM

for my $name (sort keys %{$classes{'Session'}}) {
    my $methods = $classes{'Session'}->{$name};
    for my $method (@$methods) {
        say $tmpout $method;
    }
}

print $tmpout <<'EOM';

class Octav(object):
  def __init__(self, endpoint, key, secret, debug=False):
    if not endpoint:
      raise MissingRequiredArgument('endpoint is required')
    if not key:
      raise MissingRequiredArgument('key is required')
    if not secret:
      raise MissingRequiredArgument('secret is required')
    self.debug = debug
    self.endpoint = endpoint
    self.error = None
    self.http = PoolManager(cert_reqs='CERT_REQUIRED', ca_certs=certifi.where())
    self.key = key
    self.secret = secret

  def new_session(self, access_token, auth_via):
    if not access_token:
      raise MissingRequiredArgument('access_token is required')
    if not auth_via:
      raise MissingRequiredArgument('auth_via is required')

    f = functools.partial(self.create_client_session, access_token, auth_via)
    s = f()
    if s is None:
        return None
    return Session(self, f, s.get('sid'), parse_rfc3339(s.get('expires')))

  def extract_error(self, r):
    try:
      js = json.loads(r.data)
      if 'error' in js:
        self.error = js['error']
      elif 'message' in js:
        self.error = js['message']
    except BaseException:
      self.error = r.status

  def last_error(self):
    return self.error

  def last_response(self):
    return self.res
EOM

for my $name (sort keys %{$classes{'Octav'}}) {
    my $methods = $classes{'Octav'}->{$name};
    for my $method (@$methods) {
        say $tmpout $method;
    }
}


print $buf;

sub camelize_title {
    my $s = shift;
    $s =~ s/[\W+]([\w])/_$1/g;
    return lc($s);
}

sub print_method_signature {
    my ($fh, $name, @args) = @_;
    print $fh "  def $name (self";
    if (@args) {
        print $fh ", ", join(", ", @args);
    }
    print $fh "):\n";
}
