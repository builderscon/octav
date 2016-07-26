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

sub credentials {
    my $self = shift;
    my $uri = URI->new($self->{endpoint});
    $self->{user_agent}->credentials($uri->host_port, "octav", @_[0], @_[1])
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
        say $tmpout '    my @request_args;';
        if ($link->{encType} eq 'multipart/form-data') {
            # If the encType is multipart/form-data, we must generate
            # a different type of request

            say $tmpout '    my @content;';
            # first, remove and push into content args the files that
            # should be uploaded
            if (my @files = @{$link->{"hsup.multipartFiles"} || []}) {
                say $tmpout '    for my $file (qw(' . join(" ", @files) . ')) {';
                say $tmpout '        if (my $fn = delete $payload->{$file}) {';
                say $tmpout '            push @content, ($file => [$fn]);';
                say $tmpout '        }';
                say $tmpout '    }';
            }
            say $tmpout '    push @content, (payload => JSON::encode_json($payload));';
            say $tmpout '    push @request_args, (Content_Type => "form-data");';
            say $tmpout '    push @request_args, (Content => \@content);';
        } else {
            say $tmpout '    push @request_args, (Content_Type => "application/json", Content => JSON::encode_json($payload));';
        }
        say $tmpout '    my $res = $self->{user_agent}->post($uri, @request_args);';
    } else {
        say $tmpout '    $uri->query_form($payload);';
        say $tmpout '    my $res = $self->{user_agent}->get($uri);';
    }
    say $tmpout '    if (!$res->is_success) {';
    say $tmpout '        $self->{last_error} = $res->status_line;';
    say $tmpout '        return;';
    say $tmpout '    }';
    if ($link->{targetSchema}) {
        say $tmpout '    return JSON::decode_json($res->decoded_content);';
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
    return lc($s);
}