package Octav::AdminWeb;
use 5.008001;
use strict;
use warnings;
use Mojo::Base qw(Mojolicious);
use Mojolicious::Plugin::XslateRenderer;
use WebService::Octav;

our $VERSION = "0.01";

sub startup {
    my $self = shift;
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


    for my $resource (qw(conference venue room)) {
        my $r_resource = $r->under("/$resource");
        for my $action (qw(lookup list)) {
            $r_resource->get("/$action")->to("$resource#$action");
        }

        for my $action (qw(create update delete)) {
            $r_resource->post("/$action")->to("$resource#$action");
        }
    }

    $self->helper(client => sub {
        WebService::Octav->new(endpoint => "http://130.211.245.78:8080");
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

