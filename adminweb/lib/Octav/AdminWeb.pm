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
        my $c = shift;
        # FIXME endpoint should not be hardcoded
        my $host = $ENV{APISERVER_SERVICE_HOST};
        my $port = $ENV{APISERVER_SERVICE_PORT};
        my $endpoint = "http://$host:$port";
        warn "ENDPOINT = $endpoint";
        state $client;
        if (! $client) {
            $client = WebService::Octav->new(endpoint => $endpoint);
            my $client_key = $c->config->{client_key};
            my $client_secret = $c->config->{client_secret};
            if ($client_key && $client_secret) {
                $client->credentials($client_key, $client_secret);
            }
        }
        return $client;
    });

    $self->helper(config => sub {
        state $config = +{
            "github" => {
                "auth_endpoint" => "https://github.com/login/oauth/authorize",
                "client_id" => load_from_file($ENV{OCTAV_GITHUB_CLIENT_ID}),
                "client_secret" => load_from_file($ENV{OCTAV_GITHUB_CLIENT_SECRET}),
            },
            "googlemaps" => {
                "api_key" => load_from_file($ENV{OCTAV_GOOGLE_MAPS_API_KEY}),
            },
            "client_key" => load_from_file($ENV{OCTAV_API_CLIENT_KEY}),
            "client_secret" => load_from_file($ENV{OCTAV_API_CLIENT_SECRET}),
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

    $r->get("/auth/logout")->to("auth#logout");
    $r->get("/auth")->to("auth#index");
    for my $resource (qw(github)) {
        my $r_resource = $r->under("/auth");
        $r_resource->get("/$resource")->to("auth#$resource");
        $r_resource->get("/${resource}_cb")->to("auth#${resource}_cb");
    }

    for my $resource (qw(conference user venue room)) {
        my $r_resource = $r->under("/$resource");
        for my $action (qw(edit lookup list)) {
            $r_resource->get("/$action")->to("$resource#$action");
        }

        for my $action (qw(create update delete)) {
            $r_resource->post("/$action")->to("$resource#$action");
        }
    }

    $r->get("/user/dashboard")->to("user#dashboard");

    $self->hook(around_action => sub {
        my ($next, $c, $action, $last) = @_;

        # Not a great idea, but for now, we just want to avoid the silly
        # health checkers
        my $ua = $c->req->headers->user_agent();
        if ($ua =~ m{^GoogleHC/}) {
            $c->render(text => "Yes, I'm healthy");
            return
        }

        my $endpoint = $c->match->endpoint->to_string;
        if ($endpoint !~ m{^/auth(?:/|$)}) {
            my $session = $c->plack_session;
            if (! $session) {
                warn "Access to protected resource '$endpoint' detected, but no session found";
                $c->redirect_to($c->url_for("/auth"));
                return
            }

            if (! $session->{user}) {
                warn "Access to protected resource '$endpoint' detected, but no user information found in session";
                $c->redirect_to($c->url_for("/auth"));
                return
            }
            $c->stash(ui_user => $session->{user});

            if (! $session->{user}->{is_admin}) {
                warn "User is not an administrator";
                # Instead of just erroring out here, redirect to
                # dashboard so we can see the error in a user
                # friendly format
                if ($endpoint !~ m{^/user/dashboard$}) {
                    $c->redirect_to($c->url_for("/user/dashboard"));
                    return
                }
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

