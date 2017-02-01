CREATE TABLE users (
    oid INTEGER UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    eid CHAR(64) CHARACTER SET latin1 NOT NULL,
    first_name TEXT,
    last_name TEXT,
    lang CHAR(16) NOT NULL DEFAULT 'en',
    nickname  CHAR(128) NOT NULL,
    email TEXT,
    auth_via CHAR(16) NOT NULL, /* github, facebook, twitter */
    auth_user_id TEXT NOT NULL, /* ID in the auth provider */
    avatar_url TEXT,
    is_admin TINYINT(1) NOT NULL DEFAULT 0,
    tshirt_size CHAR(4) CHARACTER SET latin1,
    timezone CHAR(32) NOT NULL DEFAULT 'UTC',
    created_on DATETIME NOT NULL,
    modified_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY(eid),
    UNIQUE KEY(auth_via, auth_user_id(191))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE venues (
    oid INTEGER UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    eid CHAR(64) CHARACTER SET latin1 NOT NULL,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    place_id TEXT,
    url TEXT,
    latitude DECIMAL(10,7),
    longitude DECIMAL(10,7),
    created_on DATETIME NOT NULL,
    modified_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY(eid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE rooms (
    oid INTEGER UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    venue_id CHAR(64) CHARACTER SET latin1 NOT NULL,
    eid CHAR(64) CHARACTER SET latin1 NOT NULL,
    name TEXT NOT NULL,
    capacity INT UNSIGNED NOT NULL DEFAULT 0,
    created_on DATETIME NOT NULL,
    modified_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY (eid),
    KEY(venue_id, eid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE conference_series (
    oid INTEGER UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    eid CHAR(64) CHARACTER SET latin1 NOT NULL,
    slug TEXT NOT NULL,
    title TEXT NOT NULL,
    created_on DATETIME NOT NULL,
    modified_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY (eid),
    UNIQUE KEY (slug(191))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
    
CREATE TABLE conference_series_administrators (
    oid INTEGER UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    series_id CHAR(64) CHARACTER SET latin1 NOT NULL,
    user_id CHAR(64) CHARACTER SET latin1 NOT NULL,
    created_on DATETIME NOT NULL,
    modified_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (series_id) REFERENCES conference_series(eid) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(eid) ON DELETE CASCADE,
    UNIQUE KEY(series_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE conferences (
    oid INTEGER UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    eid CHAR(64) CHARACTER SET latin1 NOT NULL,
    series_id CHAR(64) CHARACTER SET latin1 NOT NULL,
    slug TEXT NOT NULL,
    title TEXT NOT NULL,
    sub_title TEXT,
    cover_url TEXT,
    redirect_url TEXT,
    timetable_available TINYINT(1) NOT NULL DEFAULT 0,
    blog_feedback_available TINYINT(1) NOT NULL DEFAULT 0,
    status CHAR(64) CHARACTER SET latin1 NOT NULL default "private",
    timezone CHAR(32) NOT NULL DEFAULT 'UTC',
    created_by CHAR(64) CHARACTER SET latin1 NOT NULL,
    created_on DATETIME NOT NULL,
    modified_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (series_id) REFERENCES conference_series(eid),
    UNIQUE KEY(eid),
    UNIQUE KEY(series_id, slug(191)),
    KEY (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE tracks (
    oid INTEGER UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    eid CHAR(64) CHARACTER SET latin1 NOT NULL,
    conference_id CHAR(64) CHARACTER SET latin1 NOT NULL,
    room_id CHAR(64) CHARACTER SET latin1 NOT NULL,
    name TEXT NOT NULL,
    sort_order INTEGER UNSIGNED NOT NULL DEFAULT 0,
    created_on DATETIME NOT NULL,
    modified_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY (eid),
    UNIQUE KEY (conference_id, room_id),
    KEY(sort_order),
    FOREIGN KEY (room_id) REFERENCES rooms(eid) ON DELETE CASCADE,
    FOREIGN KEY (conference_id) REFERENCES conferences(eid) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- key, value pairs of long texts that go with a conference
CREATE TABLE conference_components (
    oid           INTEGER UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    eid           CHAR(64) CHARACTER SET latin1 NOT NULL,
    conference_id CHAR(64) CHARACTER SET latin1 NOT NULL,
    name          CHAR(64) CHARACTER SET latin1 NOT NULL,
    value         TEXT NOT NULL,
    created_on    DATETIME NOT NULL,
    modified_on   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY(eid),
    KEY(name),
    FOREIGN KEY (conference_id) REFERENCES conferences(eid) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE conference_dates (
    oid INTEGER UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    eid           CHAR(64) CHARACTER SET latin1 NOT NULL,
    conference_id CHAR(64) CHARACTER SET latin1 NOT NULL,
    open DATETIME,
    close DATETIME,
    KEY(open),
    FOREIGN KEY (conference_id) REFERENCES conferences(eid) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE external_resources (
    oid INTEGER UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    eid           CHAR(64) CHARACTER SET latin1 NOT NULL,
    conference_id CHAR(64) CHARACTER SET latin1 NOT NULL,
    description TEXT,
    image_url TEXT,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    sort_order int not null default 0,
    UNIQUE KEY(eid),
    FOREIGN KEY (conference_id) REFERENCES conferences(eid) ON DELETE CASCADE,
    KEY(sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE conference_administrators (
    oid INTEGER UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    conference_id CHAR(64) CHARACTER SET latin1 NOT NULL,
    user_id CHAR(64) CHARACTER SET latin1 NOT NULL,
    sort_order int not null default 0,
    created_on DATETIME NOT NULL,
    modified_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (conference_id) REFERENCES conferences(eid) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(eid) ON DELETE CASCADE,
    KEY(sort_order),
    UNIQUE KEY(conference_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE conference_staff LIKE conference_administrators;

CREATE TABLE conference_venues (
    oid INTEGER UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    conference_id CHAR(64) CHARACTER SET latin1 NOT NULL,
    venue_id CHAR(64) CHARACTER SET latin1 NOT NULL,
    created_on DATETIME NOT NULL,
    modified_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (conference_id) REFERENCES conferences(eid) ON DELETE CASCADE,
    FOREIGN KEY (venue_id) REFERENCES venues(eid) ON DELETE CASCADE,
    UNIQUE KEY(conference_id, venue_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- per conference session types. When a new conference is created,
-- a standard set of types are automatically created.
CREATE TABLE session_types (
    oid              INTEGER UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    eid              CHAR(64) CHARACTER SET latin1 NOT NULL,
    conference_id    CHAR(64) CHARACTER SET latin1,
    name             TEXT NOT NULL, -- "Lightning Talk"
    abstract         TEXT NOT NULL, -- "5 minute talks about anything you want"
    duration         INTEGER UNSIGNED NOT NULL,
    is_default       TINYINT(1) NOT NULL DEFAULT 0,
    sort_order       INTEGER DEFAULT 0,
    submission_start DATETIME,
    submission_end   DATETIME,
    created_on       DATETIME NOT NULL,
    modified_on      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY (eid),
    KEY (sort_order),
    FOREIGN KEY (conference_id) REFERENCES conferences(eid) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE sessions (
    oid INTEGER UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    eid CHAR(64) CHARACTER SET latin1 NOT NULL,
    conference_id CHAR(64) CHARACTER SET latin1,
    room_id CHAR(64) CHARACTER SET latin1,
    speaker_id CHAR(64) CHARACTER SET latin1,
    session_type_id CHAR(64) CHARACTER SET latin1,
    title TEXT,
    abstract TEXT,
    memo TEXT,
    starts_on DATETIME,
    duration INTEGER UNSIGNED,
    material_level TEXT,
    tags TEXT,
    category TEXT,
    selection_result_sent TINYINT(1) NOT NULL DEFAULT 0,
    spoken_language TEXT,
    slide_language TEXT,
    slide_subtitles TEXT,
    slide_url TEXT,
    video_url TEXT,
    photo_release CHAR(16) NOT NULL DEFAULT "allow",
    recording_release CHAR(16) NOT NULL DEFAULT "allow",
    materials_release CHAR(16) NOT NULL DEFAULT "allow",
    has_interpretation TINYINT(1) NOT NULL DEFAULT 0,
    status CHAR(16) NOT NULL DEFAULT "pending",
    sort_order INTEGER NOT NULL DEFAULT 0,
    confirmed TINYINT(1) NOT NULL DEFAULT 0,
    created_on DATETIME NOT NULL,
    modified_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY (session_type_id),
    KEY (starts_on),
    FOREIGN KEY (session_type_id) REFERENCES session_types(eid),
    FOREIGN KEY (speaker_id) REFERENCES users(eid) ON DELETE SET NULL,
    FOREIGN KEY (conference_id) REFERENCES conferences(eid) ON DELETE SET NULL,
    UNIQUE KEY (eid),
    KEY(eid, conference_id, room_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE localized_strings (
    oid INTEGER UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    parent_id CHAR(64) CHARACTER SET latin1 NOT NULL,
    parent_type CHAR(64) CHARACTER SET latin1 NOT NULL,
    name CHAR(128) BINARY NOT NULL,
    language CHAR(32) BINARY NOT NULL,
    localized TEXT NOT NULL,
    UNIQUE KEY (parent_id, parent_type, name, language)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE questions (
    oid INTEGER UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    eid CHAR(64) CHARACTER SET latin1 NOT NULL,
    session_id CHAR(64) CHARACTER SET latin1 NOT NULL,
    user_id CHAR(64) CHARACTER SET latin1 NOT NULL,
    body TEXT NOT NULL,
    KEY (eid, session_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- session survey responses.
CREATE TABLE session_survey_responses (
    oid INTEGER UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    eid CHAR(64) CHARACTER SET latin1 NOT NULL,
    session_id CHAR(64) CHARACTER SET latin1 NOT NULL,
    user_id CHAR(64) CHARACTER SET latin1 NOT NULL,
    user_prior_knowledge SMALLINT UNSIGNED NOT NULL DEFAULT 0,
    speaker_knowledge SMALLINT UNSIGNED NOT NULL DEFAULT 0,
    speaker_presentation SMALLINT UNSIGNED NOT NULL DEFAULT 0,
    material_quality SMALLINT UNSIGNED NOT NULL DEFAULT 0,
    overall_rating SMALLINT UNSIGNED NOT NULL DEFAULT 0,
    comment_good TEXT,
    comment_improvement TEXT,
    created_on DATETIME NOT NULL,
    modified_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY (eid),
    KEY(eid, session_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- client stores data about clients that use our API
CREATE TABLE clients (
    oid INTEGER UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    eid CHAR(64) CHARACTER SET latin1 NOT NULL, -- client ID
    secret CHAR(64) BINARY CHARACTER SET latin1 NOT NULL,
    name TEXT NOT NULL, -- name of the client
    created_on DATETIME NOT NULL,
    modified_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY (eid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Featured speakers are those speakers that you want have a separate
-- section in your page, to describe who they are and what value they
-- bring to your conference.
-- These entries may be associated with an actual user, but often times
-- you want to put these stories *BEFORE* they have the actual talk
-- materials, or even, before they can commit to dates.
-- Therefore this table replicates some of `users` table. Consumers
-- should use the actual user data if the `user_id` field is populated
CREATE TABLE featured_speakers (
    oid INTEGER UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    eid CHAR(64) CHARACTER SET latin1 NOT NULL,
    conference_id CHAR(64) CHARACTER SET latin1 NOT NULL, -- featured speakers are bound to a conference.
    speaker_id CHAR(64) CHARACTER SET latin1, -- If non-null, is linked to an actual user
    display_name TEXT NOT NULL, -- consolidated because we just need a name to show
    description TEXT NOT NULL, -- text to be displayed in the featured section
    avatar_url TEXT, -- if null, we should provide a sane default
    created_on DATETIME NOT NULL,
    modified_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY(eid),
    FOREIGN KEY (conference_id) REFERENCES conferences(eid) ON DELETE CASCADE,
    FOREIGN KEY (speaker_id) REFERENCES users(eid) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE sponsors (
    oid INTEGER UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    eid CHAR(64) CHARACTER SET latin1 NOT NULL,
    conference_id CHAR(64) CHARACTER SET latin1 NOT NULL, -- sponsors are bound to a conference.
    name TEXT NOT NULL,
    logo_url1 TEXT, -- it is up to the consumer to choose which logo to use
    logo_url2 TEXT,
    logo_url3 TEXT,
    url TEXT NOT NULL,
    group_name CHAR(64) CHARACTER SET latin1 NOT NULL,
    sort_order int not null default 0,
    created_on DATETIME NOT NULL,
    modified_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY(eid),
    KEY(sort_order, group_name),
    FOREIGN KEY (conference_id) REFERENCES conferences(eid) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE temporary_emails (
    oid INTEGER UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    user_id CHAR(64) CHARACTER SET latin1 NOT NULL,
    email TEXT CHARACTER SET latin1 NOT NULL,
    confirmation_key CHAR(64) BINARY NOT NULL,
    expires_on DATETIME NOT NULL,
    UNIQUE KEY(confirmation_key),
    UNIQUE KEY(user_id),
    FOREIGN KEY (user_id) REFERENCES users(eid) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE blog_entries (
    oid INTEGER UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    eid CHAR(64) CHARACTER SET latin1 NOT NULL,
    conference_id CHAR(64) CHARACTER SET latin1 NOT NULL,
    title TEXT,
    url TEXT NOT NULL,
    url_hash CHAR(64) CHARACTER SET latin1 NOT NULL,
    status CHAR(16) NOT NULL DEFAULT 'private',
    created_on DATETIME NOT NULL,
    modified_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY(eid),
    UNIQUE KEY(url_hash),
    KEY(status),
    FOREIGN KEY (conference_id) REFERENCES conferences(eid) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

