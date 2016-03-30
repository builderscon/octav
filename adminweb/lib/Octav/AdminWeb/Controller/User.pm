package Octav::AdminWeb::Controller::User;
use Mojo::Base qw(Mojolicious::Controller);

sub dashboard {
    my $self = shift;
}

sub list {
    my $self = shift;

    my $client = $self->client();
    my $users = $client->list_user();
    $self->stash(users => $users);
    $self->render(tx => "user/list");
}

sub lookup {
    my $self = shift;

    my $id = $self->param('id');
    if (!$id) {
        return $self->render(text => "not found", status => 404);
    }

    my $client = $self->client;
    my $user = $client->lookup_user({id => $id, lang => "all"});
    $self->stash(user => $user);
    $self->render(tx => "user/lookup");
}

sub update {
    my $self = shift;

    my $id = $self->param('id');
    if (!$id) {
        return $self->render(text => "not found", status => 404);
    }

    my $client = $self->client;
    my $user = $client->lookup_user({id => $id, lang => "all"});

    my %params = (id => $id);
    for my $pname (qw(FIXME)) {
        my $pvalue = $self->param($pname);
        if ($pvalue ne $user->{$pname}) {
            $params{$pname} = $pvalue;
        }
    }

    if (!$client->update_user(\%params)) {
        die "failed";
    }

    $self->redirect_to($self->url_for('lookup')->query(id => $id));
}

1;
