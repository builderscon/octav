package Octav::AdminWeb::Controller::Room;
use Mojo::Base qw(Mojolicious::Controller);
use JSON::Types ();

sub list {
    my $self = shift;

    my $log = $self->app->log;
    my $venue_id = $self->param('venue_id');
    if (!$venue_id) {
        $log->debug("No 'venue_id' available in query");
        return $self->render(text => "not found", status => 404);
    }

    my $client = $self->client();
    my $rooms = $client->list_room({venue_id => $venue_id});
    $self->stash(rooms => $rooms);
    $self->render(tx => "room/list");
}

sub _lookup {
    my $self = shift;

    my $log = $self->app->log;
    my $id = $self->param('id');
    if (!$id) {
        $log->debug("No 'id' available in query");
        $self->render(text => "not found", status => 404);
        return
    }

    my $client = $self->client;
    my $room = $client->lookup_room({id => $id, lang => "all"});
    if (! $room) {
        $log->debug("No such room '$id'");
        $self->render(text => "not found", status => 404);
        return
    }
    $self->stash(room => $room);
    $self->stash(api_key => $self->config->{googlemaps}->{api_key});
    return 1
}

sub lookup {
    my $self = shift;
    if (! $self->_lookup()) {
        return
    }
    $self->render(tx => "room/lookup");
}

sub edit {
    my $self = shift;
    if (! $self->_lookup()) {
        return
    }
    $self->render(tx => "room/edit");
}

sub update {
    my $self = shift;

    my $id = $self->param('id');
    if (!$id) {
        return $self->render(text => "not found", status => 404);
    }

    my $client = $self->client;
    my $room = $client->lookup_room({id => $id, lang => "all"});

    my %params = (id => $id);
    for my $pname (qw(name name#ja capacity)) {
        my $pvalue = $self->param($pname);
        if ($pvalue ne $room->{$pname}) {
            if ($pname =~ /^(?:latitude|longitude)$/) {
                $pvalue = JSON::Types::number($pvalue);
            }
            $params{$pname} = $pvalue;
        }
    }

    if (!$client->update_room(\%params)) {
        die "failed";
    }

    $self->redirect_to($self->url_for('room/lookup')->query(id => $id));
}

1;