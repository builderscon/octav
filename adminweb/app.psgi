use strict;
use Octav::AdminWeb;
use Mojo::Server::PSGI;


my $server = Mojo::Server::PSGI->new(app => Octav::AdminWeb->new());
return $server->to_psgi_app;