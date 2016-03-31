use strict;
use Test::More 0.98;

use_ok $_ for qw(
    Octav::AdminWeb::Controller::Auth
    Octav::AdminWeb::Controller::Conference
    Octav::AdminWeb::Controller::Room
    Octav::AdminWeb::Controller::Root
    Octav::AdminWeb::Controller::User
    Octav::AdminWeb::Controller::Venue
    Octav::AdminWeb
);

done_testing;

