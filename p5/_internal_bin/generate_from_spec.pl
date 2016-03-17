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

say $tmpout <<'EOM';
package WebService::Octav;
use strict;
use JSON;
use LWP::UserAgent;
use URI;

sub new {
    my ($class, %args) = @_;
    if (! $args{endpoint}) {
        die "You must supply an endpoint";
    }
    my $endpoint = $args{endpoint};
    $endpoint =~ s{/$}{}; # strip trailing "/"
    my $self = bless {
        endpoint => $endpoint,
        user_agent => LWP::UserAgent->new(agent => "perl5/WebService::Octav"),
    }, $class;
    return $self;
}

sub last_error {
    my $self = shift;
    return $self->{last_error};
}
EOM

for my $link (@{$schema->{links}}) {
    my $name = camelize_title($link->{title});

    my $path = $link->{href};
    if ($schema->{pathStart}) {
        $path = $schema->{pathStart} . $path;
    }
    say $tmpout "\nsub $name {";
    say $tmpout '    my ($self, $payload) = @_;';

    if (my $link_schema = $link->{schema}) {
        my $required = $link_schema->{required};
        if ($required && scalar(@{$required}) > 0) {
            say $tmpout '    for my $required (qw(' . join(" ", @{$required}) . ')) {';
            say $tmpout '        if (!$payload->{$required}) {';
            say $tmpout '            die qq|property "$required" must be provided|;';
            say $tmpout '        }';
            say $tmpout '    }';
        }
    }

    say $tmpout '    my $uri = URI->new($self->{endpoint} . qq|' . $path . '|);';
    if (lc($link->{method}) eq 'post') {
        say $tmpout '    my $json_payload = JSON::encode_json($payload);';
        say $tmpout '    my $res = $self->{user_agent}->post($uri, Content => $json_payload);';
    } else {
        say $tmpout '    $uri->query_form($payload);';
        say $tmpout '    my $res = $self->{user_agent}->get($uri);';
    }
    say $tmpout '    if (!$res->is_success) {';
    say $tmpout '        $self->{last_error} = $res->status_line;';
    say $tmpout '        return;';
    say $tmpout '    }';
    if ($link->{targetSchema}) {
        say $tmpout '    return JSON::decode_json($res->content);';
    } else {
        say $tmpout '    return 1';
    }
    say $tmpout '}'
}

say $tmpout "\n1;";

my $dstfile = File::Spec->catfile("lib", "WebService", "Octav.pm");
my $dstdir  = File::Basename::dirname($dstfile);
if (! -e $dstdir) {
    if (! File::Path::make_path($dstdir)) {
        die "Could not create directory $dstdir: $!";
    }
}
open my $fh, '>', $dstfile or die "Could not open file $dstfile: $!";
print $fh $buf;

sub camelize_title {
    my $s = shift;
    $s =~ s/[\W+]([\w])/_$1/g;
    return lcfirst($s);
}