package Octav::AdminWeb::Controller::Venue;
use Mojo::Base qw(Mojolicious::Controller);
use JSON::Types ();

sub dashboard {
    my $self = shift;
}

sub list {
    my $self = shift;

    my $client = $self->client();
    my $venues = $client->list_venue();
    $self->stash(venues => $venues);
    $self->render(tx => "venue/list");
}

sub lookup {
    my $self = shift;

    my $id = $self->param('id');
    if (!$id) {
        return $self->render(text => "not found", status => 404);
    }

    my $client = $self->client;
    my $venue = $client->lookup_venue({id => $id, lang => "all"});
    $self->stash(venue => $venue);
    $self->stash(api_key => $self->config->{googlemaps}->{api_key});
    $self->render(tx => "venue/lookup");
}

sub update {
    my $self = shift;

    my $id = $self->param('id');
    if (!$id) {
        return $self->render(text => "not found", status => 404);
    }

    my $client = $self->client;
    my $venue = $client->lookup_venue({id => $id, lang => "all"});

    my %params = (id => $id);
    for my $pname (qw(name name#ja latitude longitude)) {
        my $pvalue = $self->param($pname);
        if ($pvalue ne $venue->{$pname}) {
            if ($pname =~ /^(?:latitude|longitude)$/) {
                $pvalue = JSON::Types::number($pvalue);
            }
            $params{$pname} = $pvalue;
        }
    }

    if (!$client->update_venue(\%params)) {
        die "failed";
    }

    $self->redirect_to($self->url_for('venue/lookup')->query(id => $id));
}

1;