package Octav::AdminWeb::Controller::Auth;
use Mojo::Base qw(Mojolicious::Controller);
use JSON ();
use JSON::Types ();
use Mojo::UserAgent;
use URI;
use UUID::Tiny qw(:std);

sub index {
    my $self = shift;
    if (my $error = $self->param('error')) {
        $self->stash(error => $error);
    }
    $self->render(tx => "auth/index");
}

sub logout {
    my $self = shift;
    delete $self->plack_session->{user};
    $self->redirect_to("/");
}

sub github {
    my $self = shift;

    my $log = $self->app->log;
    my $github_config = $self->config->{github};
    my $uri = URI->new($github_config->{auth_endpoint});
    my $redirect_uri = "https://admin.builderscon.io/auth/github_cb";
    my $state = unpack("H*", create_uuid(UUID_RANDOM));
    my $session = $self->plack_session;
    $session->{"github_state"} = $state;
    $session->{"github_state_expires"} = time() + 120;

    $log->debug("Saving github_state = $state");
    $uri->query_form(
        client_id    => $github_config->{client_id},
        redirect_uri => $redirect_uri,
        scope        => "user",
        state        => $state,
    );
    $self->redirect_to($uri->as_string);
}

sub github_cb {
    my $self = shift;

    my $log = $self->app->log;

    $log->debug("in github_cb");

    if (my $error = $self->param("error")) {
        $self->redirect_to($self->url_for("/auth")->query(error => $error));
    }

    my $code = $self->param("code");
    my $state = $self->param("state");
    my $github_config = $self->config->{github};

    my $session = $self->plack_session;
    my $github_state = delete $session->{github_state};
    if (!$github_state || $github_state ne $state) {
        $log->debug("No github_state key in session");
        $self->redirect_to($self->url_for("/auth"));
        return;
    }

    if (delete $session->{github_state_expires} < time()) {
        $self->redirect_to($self->url_for("/auth"));
        return;
    }

    my $ua = Mojo::UserAgent->new;
    my $tx = $ua->post("https://github.com/login/oauth/access_token", 
        {Accept => "application/json"},
        form => {
            client_id     => $github_config->{client_id},
            client_secret => $github_config->{client_secret},
            code          => $code,
        },
    );

    my $res = $tx->success;
    if (!$res) {
        my $err = $tx->error;
        die "$err->{code} response: $err->{message}" if $err->{code};
        die "Connection error: $err->{message}";
    }

    my $oauth_response = JSON::decode_json($res->body);
    if (my $error = $oauth_response->{error}) {
        if ($error eq "bad_verification_code") {
            $self->redirect_to($self->url_for("/auth"));
            return
        } else {
            die "real problem: $error";
        }
    }
    
    $self->render(text => $oauth_response->{access_token});
    $tx = $ua->get("https://api.github.com/user", {
        Authorization => "token $oauth_response->{access_token}",
    });
    $res = $tx->success;
    if (!$res) {
        die "Failed to get user data from github";
    }

    my $user = JSON::decode_json($res->body);
    my $auth_user_id = JSON::Types::string($user->{id});
    my %data = (
        auth_via     => "github",
        auth_user_id => $auth_user_id,
        avatar_url   => $user->{avatar_url},
        nickname     => $user->{login},
    );
    
    # name from github *could* be first/last name
    if (my $name = $user->{name}) {
        my ($first_name, $last_name) = split(/\s+/, $name, 2);
        if ($first_name && $last_name) {
            $data{first_name} = $first_name;
            $data{last_name} = $last_name;
        }
    }

    if (my $email = $user->{email}) {
        $data{email} = $email;
    }

    # Race condition ahead...
    my $client = $self->client();
    my $registered = $client->lookup_user_by_auth_user_id({auth_via => "github", auth_user_id => $auth_user_id});
    if (!$registered) {
        $registered = $client->create_user(\%data);
    }

    $log->debug("Setting session->{user}");
    $session->{user} = $registered;

    $self->redirect_to($self->url_for("/user/dashboard"));
}

1;

