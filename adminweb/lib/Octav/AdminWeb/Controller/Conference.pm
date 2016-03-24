package Octav::AdminWeb::Controller::Conference;
use Mojo::Base qw(Mojolicious::Controller);

sub list {
    my $self = shift;

    my $client = $self->client();
    my $conferences = $client->list_conference();
    $self->stash(conferences => $conferences);
    $self->render(tx => "conference/list");
}

sub lookup {
    my $self = shift;

    my $id = $self->param('id');
    if (!$id) {
        return $self->render(text => "not found", status => 404);
    }

    my $client = $self->client;
    my $conference = $client->lookup_conference({id => $id, lang => "all"});
    $self->stash(conference => $conference);
    $self->render(tx => "conference/lookup");
}

sub update {
    my $self = shift;

    my $id = $self->param('id');
    if (!$id) {
        return $self->render(text => "not found", status => 404);
    }

    my $client = $self->client;
    my $conference = $client->lookup_conference({id => $id, lang => "all"});

    my %params = (id => $id);
    for my $pname (qw(title sub_title title#ja sub_title#ja slug)) {
        my $pvalue = $self->param($pname);
        if ($pvalue ne $conference->{$pname}) {
            $params{$pname} = $pvalue;
        }
    }

    if (!$client->update_conference(\%params)) {
        die "failed";
    }

    $self->redirect_to($self->url_for('lookup')->query(id => $id));
}

1;