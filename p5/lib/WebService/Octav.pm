package WebService::Octav;
use strict;
use JSON;
use LWP::UserAgent;
use URI;

sub new {
    my ($class, %args) = @_;
    if (! $args{endpoint}) {
        die "You must supply an endpoint";
    }
    my $endpoint = $args{endpoint};
    $endpoint =~ s{/$}{}; # strip trailing "/"
    my $self = bless {
        endpoint => $endpoint,
        user_agent => LWP::UserAgent->new(agent => "perl5/WebService::Octav"),
    }, $class;
    return $self;
}

sub last_error {
    my $self = shift;
    return $self->{last_error};
}


sub create_user {
    my ($self, $payload) = @_;
    for my $required (qw(first_name last_name nickname email tshirt_size)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/user/create|);
    my $json_payload = JSON::encode_json($payload);
    my $res = $self->{user_agent}->post($uri, Content => $json_payload);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return JSON::decode_json($res->content);
}

sub lookup_user {
    my ($self, $payload) = @_;
    for my $required (qw(id)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/user/lookup|);
    $uri->query_form($payload);
    my $res = $self->{user_agent}->get($uri);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return JSON::decode_json($res->content);
}

sub update_user {
    my ($self, $payload) = @_;
    for my $required (qw(id)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/user/update|);
    my $json_payload = JSON::encode_json($payload);
    my $res = $self->{user_agent}->post($uri, Content => $json_payload);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return 1
}

sub delete_user {
    my ($self, $payload) = @_;
    for my $required (qw(id)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/user/delete|);
    my $json_payload = JSON::encode_json($payload);
    my $res = $self->{user_agent}->post($uri, Content => $json_payload);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return 1
}

sub create_venue {
    my ($self, $payload) = @_;
    for my $required (qw(name address)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/venue/create|);
    my $json_payload = JSON::encode_json($payload);
    my $res = $self->{user_agent}->post($uri, Content => $json_payload);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return JSON::decode_json($res->content);
}

sub list_venues {
    my ($self, $payload) = @_;
    my $uri = URI->new($self->{endpoint} . qq|/v1/venue/list|);
    $uri->query_form($payload);
    my $res = $self->{user_agent}->get($uri);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return JSON::decode_json($res->content);
}

sub lookup_venue {
    my ($self, $payload) = @_;
    for my $required (qw(id)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/venue/lookup|);
    $uri->query_form($payload);
    my $res = $self->{user_agent}->get($uri);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return JSON::decode_json($res->content);
}

sub update_venue {
    my ($self, $payload) = @_;
    for my $required (qw(id)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/venue/update|);
    my $json_payload = JSON::encode_json($payload);
    my $res = $self->{user_agent}->post($uri, Content => $json_payload);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return 1
}

sub delete_venue {
    my ($self, $payload) = @_;
    for my $required (qw(id)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/venue/delete|);
    my $json_payload = JSON::encode_json($payload);
    my $res = $self->{user_agent}->post($uri, Content => $json_payload);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return 1
}

sub create_room {
    my ($self, $payload) = @_;
    for my $required (qw(venue_id name)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/room/create|);
    my $json_payload = JSON::encode_json($payload);
    my $res = $self->{user_agent}->post($uri, Content => $json_payload);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return JSON::decode_json($res->content);
}

sub update_room {
    my ($self, $payload) = @_;
    for my $required (qw(id)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/room/update|);
    my $json_payload = JSON::encode_json($payload);
    my $res = $self->{user_agent}->post($uri, Content => $json_payload);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return 1
}

sub lookup_room {
    my ($self, $payload) = @_;
    for my $required (qw(id)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/room/lookup|);
    $uri->query_form($payload);
    my $res = $self->{user_agent}->get($uri);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return JSON::decode_json($res->content);
}

sub delete_room {
    my ($self, $payload) = @_;
    for my $required (qw(id)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/room/delete|);
    my $json_payload = JSON::encode_json($payload);
    my $res = $self->{user_agent}->post($uri, Content => $json_payload);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return 1
}

sub list_rooms {
    my ($self, $payload) = @_;
    for my $required (qw(venue_id)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/room/list|);
    $uri->query_form($payload);
    my $res = $self->{user_agent}->get($uri);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return JSON::decode_json($res->content);
}

sub create_conference {
    my ($self, $payload) = @_;
    for my $required (qw(title slug user_id)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/conference/create|);
    my $json_payload = JSON::encode_json($payload);
    my $res = $self->{user_agent}->post($uri, Content => $json_payload);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return JSON::decode_json($res->content);
}

sub add_conference_dates {
    my ($self, $payload) = @_;
    for my $required (qw(conference_id dates)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/conference/date/add|);
    my $json_payload = JSON::encode_json($payload);
    my $res = $self->{user_agent}->post($uri, Content => $json_payload);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return 1
}

sub delete_conference_dates {
    my ($self, $payload) = @_;
    for my $required (qw(conference_id dates)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/conference/date/delete|);
    my $json_payload = JSON::encode_json($payload);
    my $res = $self->{user_agent}->post($uri, Content => $json_payload);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return 1
}

sub add_conference_admin {
    my ($self, $payload) = @_;
    for my $required (qw(conference_id user_id)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/conference/admin/add|);
    my $json_payload = JSON::encode_json($payload);
    my $res = $self->{user_agent}->post($uri, Content => $json_payload);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return 1
}

sub delete_conference_admin {
    my ($self, $payload) = @_;
    for my $required (qw(conference_id user_id)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/conference/admin/delete|);
    my $json_payload = JSON::encode_json($payload);
    my $res = $self->{user_agent}->post($uri, Content => $json_payload);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return 1
}

sub lookup_conference {
    my ($self, $payload) = @_;
    for my $required (qw(id)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/conference/lookup|);
    $uri->query_form($payload);
    my $res = $self->{user_agent}->get($uri);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return JSON::decode_json($res->content);
}

sub list_conferences {
    my ($self, $payload) = @_;
    my $uri = URI->new($self->{endpoint} . qq|/v1/conference/list|);
    $uri->query_form($payload);
    my $res = $self->{user_agent}->get($uri);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return JSON::decode_json($res->content);
}

sub update_conference {
    my ($self, $payload) = @_;
    for my $required (qw(id)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/conference/update|);
    my $json_payload = JSON::encode_json($payload);
    my $res = $self->{user_agent}->post($uri, Content => $json_payload);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return 1
}

sub delete_conference {
    my ($self, $payload) = @_;
    for my $required (qw(id)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/conference/delete|);
    my $json_payload = JSON::encode_json($payload);
    my $res = $self->{user_agent}->post($uri, Content => $json_payload);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return 1
}

sub create_session {
    my ($self, $payload) = @_;
    for my $required (qw(conference_id speaker_id title abstract duration)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/session/create|);
    my $json_payload = JSON::encode_json($payload);
    my $res = $self->{user_agent}->post($uri, Content => $json_payload);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return JSON::decode_json($res->content);
}

sub lookup_session {
    my ($self, $payload) = @_;
    for my $required (qw(id)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/session/lookup|);
    $uri->query_form($payload);
    my $res = $self->{user_agent}->get($uri);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return JSON::decode_json($res->content);
}

sub delete_session {
    my ($self, $payload) = @_;
    for my $required (qw(id)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/session/delete|);
    my $json_payload = JSON::encode_json($payload);
    my $res = $self->{user_agent}->post($uri, Content => $json_payload);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return 1
}

sub update_session {
    my ($self, $payload) = @_;
    for my $required (qw(id)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/session/update|);
    my $json_payload = JSON::encode_json($payload);
    my $res = $self->{user_agent}->post($uri, Content => $json_payload);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return 1
}

sub list_sessions_by_conference {
    my ($self, $payload) = @_;
    for my $required (qw(conference_id)) {
        if (!$payload->{$required}) {
            die qq|property "$required" must be provided|;
        }
    }
    my $uri = URI->new($self->{endpoint} . qq|/v1/schedule/list|);
    $uri->query_form($payload);
    my $res = $self->{user_agent}->get($uri);
    if (!$res->is_success) {
        $self->{last_error} = $res->status_line;
        return;
    }
    return JSON::decode_json($res->content);
}

1;
