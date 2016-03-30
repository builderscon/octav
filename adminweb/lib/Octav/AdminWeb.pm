package Octav::AdminWeb;
use 5.008001;
use strict;
use warnings;
use Mojo::Base qw(Mojolicious);
use Mojolicious::Plugin::XslateRenderer;
use Redis::Jet;
use WebService::Octav;

our $VERSION = "0.01";

sub redis {
    state $redis = Redis::Jet->new(server => $ENV{OCTAV_REDIS});
    return $redis;
}

sub load_from_file {
    my $file = shift;
    my $fh;
    if (! open $fh, '<', $file) {
        warn "Could not load from file '$file': $!";
        return
    }
    local $/;
    return <$fh>;
}

sub startup {
    my $self = shift;

    $self->helper(plack_session => sub {
        my $c = shift;
        my $session = $c->req->env->{"psgix.session"};
        if (! $session) {
            $session = {};
            $c->req->env->{"psgix.session"} = $session;
        }
        return $session;
    });

    $self->helper(redis => \&redis);

    $self->helper(client => sub {
        # FIXME endpoint should not be hardcoded
        state $client = WebService::Octav->new(endpoint => "http://130.211.245.78:8080");
        return $client;
    });

    $self->helper(config => sub {
        state $config = +{
            "github" => {
                "auth_endpoint" => "https://github.com/login/oauth/authorize",
                "client_id" => load_from_file($ENV{OCTAV_GITHUB_CLIENT_ID}),
                "client_secret" => load_from_file($ENV{OCTAV_GITHUB_CLIENT_SECRET}),
            },
        };
        return $config;
    });
    $self->plugin('xslate_renderer' => {
        template_options => {
            module => [ "Text::Xslate::Bridge::TT2Like" ],
            syntax    => 'TTerse',
            tag_start => '[%',
            tag_end   => '%]',
        }
    });
    my $r = $self->routes;
    $r->get("/")->to("Root#index");

    $r->get("/auth")->to("auth#index");
    for my $resource (qw(github)) {
        my $r_resource = $r->under("/auth");
        $r_resource->get("/$resource")->to("auth#$resource");
        $r_resource->get("/${resource}_cb")->to("auth#${resource}_cb");
    }

    for my $resource (qw(conference user venue room)) {
        my $r_resource = $r->under("/$resource");
        for my $action (qw(lookup list)) {
            $r_resource->get("/$action")->to("$resource#$action");
        }

        for my $action (qw(create update delete)) {
            $r_resource->post("/$action")->to("$resource#$action");
        }
    }

    $r->get("/user/dashboard")->to("user#dashboard");

    $self->hook(around_action => sub {
        my ($next, $c, $action, $last) = @_;
        my $endpoint = $c->match->endpoint->to_string;
        if ($endpoint !~ m{^/auth(?:/|$)}) {
            my $session = $c->plack_session;
            if (! $session) {
                $c->redirect_to($c->url_for("/auth"));
                return
            }

            if (! $session->{user}) {
                $c->redirect_to($c->url_for("/auth"));
                return
            }
        }

        return $next->();
    });
}

1;
__END__

=encoding utf-8

=head1 NAME

Octav::AdminWeb - Web Frontend for Builderscon Admin Interface

=head1 SYNOPSIS

    # app.psgi
    use Octav::AdminWeb;
    return Octav::AdminWeb->psgi_app;
    

=head1 DESCRIPTION

=head1 LICENSE

Copyright (C) builderscon.

See LICENSE file for details

=head1 AUTHOR

Daisuke Maki E<lt>lestrrat@gmail.comE<gt>

=cut

