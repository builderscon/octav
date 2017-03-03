package service

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/builderscon/octav/octav/cache"
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	"github.com/lestrrat/go-pdebug"
	urlenc "github.com/lestrrat/go-urlenc"
	"github.com/pkg/errors"
)

func (v *SessionSvc) Init() {}

func (v *SessionSvc) populateRowForCreate(vdb *db.Session, payload *model.CreateSessionRequest) error {
	vdb.EID = tools.UUID()
	vdb.ConferenceID = payload.ConferenceID
	vdb.SpeakerID = payload.SpeakerID.String
	vdb.Duration = payload.Duration
	vdb.HasInterpretation = false
	vdb.Status = "pending"
	vdb.SortOrder = 0
	vdb.Confirmed = false
	vdb.SessionTypeID = payload.SessionTypeID

	// At least one of the English or Japanese titles must be
	// non-empty
	var hasTitle bool
	if payload.Title.Valid() {
		hasTitle = true
		vdb.Title.Valid = true
		vdb.Title.String = payload.Title.String
	}

	if s, ok := payload.LocalizedFields.Get("ja", "title"); ok && s != "" {
		hasTitle = true
	}

	if !hasTitle {
		return errors.New("missing title")
	}

	// At least one of the English or Japanese abstracts must be
	// non-empty
	var hasAbstract bool
	if payload.Abstract.Valid() {
		hasAbstract = true
		vdb.Abstract.Valid = true
		vdb.Abstract.String = payload.Abstract.String
	}

	if s, ok := payload.LocalizedFields.Get("ja", "abstract"); ok && s != "" {
		hasAbstract = true
	}

	if !hasAbstract {
		return errors.New("missing abstract")
	}

	if payload.Memo.Valid() {
		vdb.Memo.Valid = true
		vdb.Memo.String = payload.Memo.String
	}

	if payload.MaterialLevel.Valid() {
		vdb.MaterialLevel.Valid = true
		vdb.MaterialLevel.String = payload.MaterialLevel.String
	}

	if payload.Category.Valid() {
		vdb.Category.Valid = true
		vdb.Category.String = payload.Category.String
	}

	if payload.SpokenLanguage.Valid() {
		vdb.SpokenLanguage.Valid = true
		vdb.SpokenLanguage.String = payload.SpokenLanguage.String
	}

	if payload.SlideLanguage.Valid() {
		vdb.SlideLanguage.Valid = true
		vdb.SlideLanguage.String = payload.SlideLanguage.String
	}

	if payload.SlideSubtitles.Valid() {
		vdb.SlideSubtitles.Valid = true
		vdb.SlideSubtitles.String = payload.SlideSubtitles.String
	}

	if payload.SlideURL.Valid() {
		vdb.SlideURL.Valid = true
		vdb.SlideURL.String = payload.SlideURL.String
	}

	if payload.VideoURL.Valid() {
		vdb.VideoURL.Valid = true
		vdb.VideoURL.String = payload.VideoURL.String
	}

	if payload.PhotoRelease.Valid() {
		vdb.PhotoRelease.Valid = true
		vdb.PhotoRelease.String = payload.PhotoRelease.String
	}

	if payload.RecordingRelease.Valid() {
		vdb.RecordingRelease.Valid = true
		vdb.RecordingRelease.String = payload.RecordingRelease.String
	}

	if payload.MaterialsRelease.Valid() {
		vdb.MaterialsRelease.Valid = true
		vdb.MaterialsRelease.String = payload.MaterialsRelease.String
	}

	if payload.Tags.Valid() {
		vdb.Tags.Valid = true
		vdb.Tags.String = string(payload.Tags.String)
	}

	return nil
}

func (v *SessionSvc) populateRowForUpdate(vdb *db.Session, payload *model.UpdateSessionRequest) error {
	if vdb.EID != payload.ID {
		return errors.New("ID mismatched for Session.populdateRowForUpdate")
	}

	if payload.ConferenceID.Valid() {
		vdb.ConferenceID = payload.ConferenceID.String
	}

	if payload.RoomID.Valid() {
		vdb.RoomID.Valid = true
		vdb.RoomID.String = payload.RoomID.String
	}

	if payload.SpeakerID.Valid() {
		vdb.SpeakerID = payload.SpeakerID.String
	}

	if payload.SessionTypeID.Valid() {
		vdb.SessionTypeID = payload.SessionTypeID.String
	}

	if payload.StartsOn.Valid() {
		vdb.StartsOn.Valid = true
		vdb.StartsOn.Time = payload.StartsOn.Time
	}

	if payload.Duration.Valid() {
		vdb.Duration = int(payload.Duration.Int)
	}

	if payload.HasInterpretation.Valid() {
		vdb.HasInterpretation = payload.HasInterpretation.Bool
	}

	if payload.Status.Valid() {
		vdb.Status = payload.Status.String
	}

	if payload.SortOrder.Valid() {
		vdb.SortOrder = int(payload.SortOrder.Int)
	}

	if payload.Confirmed.Valid() {
		vdb.Confirmed = payload.Confirmed.Bool
	}

	if payload.SelectionResultSent.Valid() {
		vdb.SelectionResultSent = payload.SelectionResultSent.Bool
	}

	if payload.Title.Valid() {
		vdb.Title.Valid = true
		vdb.Title.String = payload.Title.String
	}

	if payload.Abstract.Valid() {
		vdb.Abstract.Valid = true
		vdb.Abstract.String = payload.Abstract.String
	}

	if payload.Memo.Valid() {
		vdb.Memo.Valid = true
		vdb.Memo.String = payload.Memo.String
	}

	if payload.MaterialLevel.Valid() {
		vdb.MaterialLevel.Valid = true
		vdb.MaterialLevel.String = payload.MaterialLevel.String
	}

	if payload.Category.Valid() {
		vdb.Category.Valid = true
		vdb.Category.String = payload.Category.String
	}

	if payload.SpokenLanguage.Valid() {
		vdb.SpokenLanguage.Valid = true
		vdb.SpokenLanguage.String = payload.SpokenLanguage.String
	}

	if payload.SlideLanguage.Valid() {
		vdb.SlideLanguage.Valid = true
		vdb.SlideLanguage.String = payload.SlideLanguage.String
	}

	if payload.SlideSubtitles.Valid() {
		vdb.SlideSubtitles.Valid = true
		vdb.SlideSubtitles.String = payload.SlideSubtitles.String
	}

	if payload.SlideURL.Valid() {
		vdb.SlideURL.Valid = true
		vdb.SlideURL.String = payload.SlideURL.String
	}

	if payload.VideoURL.Valid() {
		vdb.VideoURL.Valid = true
		vdb.VideoURL.String = payload.VideoURL.String
	}

	if payload.PhotoRelease.Valid() {
		vdb.PhotoRelease.Valid = true
		vdb.PhotoRelease.String = payload.PhotoRelease.String
	}

	if payload.RecordingRelease.Valid() {
		vdb.RecordingRelease.Valid = true
		vdb.RecordingRelease.String = payload.RecordingRelease.String
	}

	if payload.MaterialsRelease.Valid() {
		vdb.MaterialsRelease.Valid = true
		vdb.MaterialsRelease.String = payload.MaterialsRelease.String
	}

	if payload.Tags.Valid() {
		vdb.Tags.Valid = true
		vdb.Tags.String = string(payload.Tags.String)
	}

	return nil
}

/*
func (v *SessionSvc) LoadByConference(tx *db.Tx, vdbl *db.SessionList, cid string, date string) error {
	if err := vdbl.LoadByConference(tx, cid, "", date, nil, nil); err != nil {
		return err
	}
	return nil
}
*/

func (v *SessionSvc) Decorate(tx *db.Tx, session *model.Session, trustedCall bool, lang string) error {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Session.Decorate")
		defer g.End()
	}
	// session must be associated with a conference
	if session.ConferenceID != "" {
		var cs ConferenceSvc
		var mc model.Conference
		if err := cs.Lookup(tx, &mc, session.ConferenceID); err != nil {
			return errors.Wrap(err, "failed to load conference")
		}
		if err := cs.Decorate(tx, &mc, trustedCall, lang); err != nil {
			return errors.Wrap(err, "failed to decorate conference")
		}
		session.Conference = &mc

		if mc.Timezone != "" && !session.StartsOn.IsZero() {
			loc, err := time.LoadLocation(mc.Timezone)
			if err == nil {
				session.StartsOn = session.StartsOn.In(loc)
			}
		}
	}

	// ... but not necessarily with a room
	if session.RoomID != "" {
		var rs RoomSvc
		var room model.Room
		if err := rs.Lookup(tx, &room, session.RoomID); err != nil {
			return errors.Wrap(err, "failed to load room")
		}
		if err := rs.Decorate(tx, &room, trustedCall, lang); err != nil {
			return errors.Wrap(err, "failed to decorate room")
		}
		session.Room = &room
	}

	if session.SpeakerID != "" {
		var su UserSvc
		var speaker model.User
		if err := su.Lookup(tx, &speaker, session.SpeakerID); err != nil {
			return errors.Wrapf(err, "failed to load speaker '%s'", session.SpeakerID)
		}
		if err := su.Decorate(tx, &speaker, trustedCall, lang); err != nil {
			return errors.Wrap(err, "failed to decorate speaker")
		}
		session.Speaker = &speaker
	}

	if session.SessionTypeID != "" {
		var sts SessionTypeSvc
		var sessionType model.SessionType
		if err := sts.Lookup(tx, &sessionType, session.SessionTypeID); err != nil {
			return errors.Wrapf(err, "failed to load session type '%s'", session.SessionTypeID)
		}
		if err := sts.Decorate(tx, &sessionType, trustedCall, lang); err != nil {
			return errors.Wrap(err, "failed to decorate session type")
		}
		session.SessionType = &sessionType
	}

	if lang != "" {
		if err := v.ReplaceL10NStrings(tx, session, lang); err != nil {
			return errors.Wrap(err, "failed to replace localized strings")
		}
	}

	return nil
}

func (v *SessionSvc) CreateFromPayload(tx *db.Tx, result *model.Session, payload *model.CreateSessionRequest) error {
	var u model.User
	su := User()
	if err := su.Lookup(tx, &u, payload.UserID); err != nil {
		return errors.Wrapf(err, "failed to load user %s", payload.UserID)
	}

	// Check if this session type is allowed to be submitted right now
	sst := SessionType()
	if err := sst.IsAcceptingSubmissions(tx, payload.SessionTypeID); err != nil {
		return errors.Wrap(err, "not accepting submissions for this session type")
	}

	// Load the session type, so we can populate payload.Duration
	var mst model.SessionType
	if err := sst.Lookup(tx, &mst, payload.SessionTypeID); err != nil {
		return errors.Wrap(err, "failed to lookup session type")
	}

	payload.Duration = mst.Duration

	var vdb db.Session
	if err := v.Create(tx, &vdb, payload); err != nil {
		return errors.Wrap(err, "failed to insert into database")
	}

	var m model.Session
	if err := m.FromRow(&vdb); err != nil {
		return errors.Wrap(err, "failed to populate model from database")
	}

	*result = m
	return nil
}

func (v *SessionSvc) PreUpdateFromPayloadHook(ctx context.Context, tx *db.Tx, vdb *db.Session, payload *model.UpdateSessionRequest) (err error) {
	su := User()
	if err := su.IsSessionOwner(tx, payload.ID, payload.UserID); err != nil {
		return errors.Wrap(err, "updating sessions require session owner privileges")
	}

	// We must protect the API server from changing important
	// fields like conference_id, speaker_id, room_id, etc from regular
	// users, but allow administrators to do anything they want
	if err := su.IsAdministrator(tx, payload.UserID); err != nil {
		// Reset the payload, whatever it is
		payload.ConferenceID.Set(vdb.ConferenceID)
		payload.SpeakerID.Set(vdb.SpeakerID)
		payload.RoomID.Set(vdb.RoomID)
		payload.SelectionResultSent.Set(vdb.SelectionResultSent)
	}
	return nil
}

func (v *SessionSvc) ListFromPayload(tx *db.Tx, result *model.SessionList, payload *model.ListSessionsRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Session.ListFromPayload").BindError(&err)
		defer g.End()
	}

	// Make sure that we have at least one of the arguments
	var conferenceID, speakerID string
	var hasQuery bool
	if payload.ConferenceID.Valid() {
		conferenceID = payload.ConferenceID.String
		hasQuery = true
	}

	if payload.SpeakerID.Valid() {
		speakerID = payload.SpeakerID.String
		hasQuery = true
	}

	var rangeStart, rangeEnd time.Time
	if payload.RangeStart.Valid() {
		dt, err := time.Parse(time.RFC3339, payload.RangeStart.String)
		if err == nil {
			rangeStart = dt.UTC()
		}
		// Don't set the hasQuery flag, as this alone doesn't work
	}
	if payload.RangeEnd.Valid() {
		dt, err := time.Parse(time.RFC3339, payload.RangeEnd.String)
		if err == nil {
			rangeEnd = dt.UTC()
		}
		// Don't set the hasQuery flag, as this alone doesn't work
	}

	status := payload.Status
	if len(status) == 0 {
		status = append(status, model.StatusAccepted)
	}

	confirmed := payload.Confirmed

	if !hasQuery {
		return errors.New("no query specified (one of conference_id/speaker_id is required)")
	}

	keybytes, err := urlenc.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "failed to marshal payload")
	}
	c := Cache()
	key := c.Key("Session", "ListFromPayload", string(keybytes))
	x, err := c.GetOrSet(key, result, func() (interface{}, error) {
		if pdebug.Enabled {
			pdebug.Printf("CACHE MISS: Re-generating")
		}

		var vdbl db.SessionList
		if err := vdbl.LoadByConference(tx, conferenceID, speakerID, rangeStart, rangeEnd, status, confirmed); err != nil {
			return nil, errors.Wrap(err, "failed to load from database")
		}

		l := make(model.SessionList, len(vdbl))
		for i, vdb := range vdbl {
			if err := l[i].FromRow(&vdb); err != nil {
				return nil, errors.Wrap(err, "failed to populate model from database")
			}

			if err := v.Decorate(tx, &l[i], payload.TrustedCall, payload.Lang.String); err != nil {
				return nil, errors.Wrap(err, "failed to decorate session with associated data")
			}
		}

		return &l, nil
	}, cache.WithExpires(10*time.Minute))

	if err != nil {
		return err
	}

	*result = *(x.(*model.SessionList))
	return nil
}

func (v *SessionSvc) DeleteFromPayload(tx *db.Tx, payload *model.DeleteSessionRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Session.DeleteFromPayload %s", payload.ID).BindError(&err)
		defer g.End()
	}

	// First, we need to load the target session
	var s model.Session
	if err := v.Lookup(tx, &s, payload.ID); err != nil {
		return errors.Wrap(err, "failed to lookup session")
	}

	if pdebug.Enabled {
		pdebug.Printf("Session status is %s", s.Status)
	}

	if s.Status == model.StatusAccepted {
		if pdebug.Enabled {
			pdebug.Printf("Session is already accepted, this required conference administrator privileges")
		}
		// The only user(s) that can delete an accepted session is an administrator.
		su := User()
		if err := su.IsConferenceAdministrator(tx, s.ConferenceID, payload.UserID); err != nil {
			return errors.Wrap(err, "deleting accepted sessions require administrator privileges")
		}
	} else {
		if s.SpeakerID != payload.UserID {
			if pdebug.Enabled {
				pdebug.Printf("User is not the owner of session, requires conference administrator privileges")
			}

			su := User()
			if err := su.IsConferenceAdministrator(tx, s.ConferenceID, payload.UserID); err != nil {
				return errors.Wrap(err, "deleting sessions require operation from speaker or user with administrator privileges")
			}
		}
	}

	if err := v.Delete(tx, payload.ID); err != nil {
		return errors.Wrap(err, "failed to delete session")
	}
	return nil
}

func (v *SessionSvc) PostSocialServices(session *model.Session) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("SessionSvc.PostSocialServices %s", session.ID).BindError(&err)
		defer g.End()
	}

	if InTesting {
		return errors.New("skipped during testing")
	}

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to start database transaction")
	}
	defer tx.AutoRollback()

	var conf model.Conference
	if err := conf.Load(tx, session.ConferenceID); err != nil {
		return errors.Wrap(err, "failed to load conference")
	}

	var series model.ConferenceSeries
	if err := series.Load(tx, conf.SeriesID); err != nil {
		return errors.Wrap(err, "failed to load conference series")
	}

	if err := v.ReplaceL10NStrings(tx, session, "all"); err != nil {
		return errors.Wrap(err, "failed to replace localized strings")
	}

	tweet, err := formatSessionTweet(session, &conf, &series)
	if err != nil {
		return errors.Wrap(err, "failed to format session before tweeting")
	}

	return Twitter().TweetAsConference(session.ConferenceID, tweet)
}

func formatSessionTweet(session *model.Session, conf *model.Conference, series *model.ConferenceSeries) (string, error) {
	prefix := "New submission "
	tweetLen := len(prefix) + 2 + 1 // prefix + 2 quotes + 1 space

	// we can post at most 140 - tweetLen
	var title string
	session.LocalizedFields.Foreach(func(lang, lk, lv string) error {
		if lk == "title" {
			title = lv
			return errors.New("stop")
		}
		return nil
	})
	if title == "" {
		title = session.Title
	}
	if title == "" {
		title = "(null)"
	}

	// XXX https://builderscon.io should probably be configurable
	u := "https://builderscon.io/" + series.Slug + "/" + conf.Slug + "/session/" + session.ID
	tweetLen = tweetLen + 23 // will be shortened

	if remain := 140 - tweetLen; utf8.RuneCountInString(title) > remain {
		var truncated bytes.Buffer
		for len(title) > 0 && remain > 1 {
			r, n := utf8.DecodeRuneInString(title)
			if r == utf8.RuneError {
				break
			}
			remain = remain - 1
			title = title[n:]
			truncated.WriteRune(r)
		}
		truncated.WriteRune('…')
		title = truncated.String()
	}

	return fmt.Sprintf("New submission %s %s",
		strconv.Quote(title),
		u,
	), nil
}

func (v *SessionSvc) SendSelectionResultNotificationFromPayload(ctx context.Context, tx *db.Tx, payload *model.SendSelectionResultNotificationRequest) error {
	var m model.Session
	if err := v.Lookup(tx, &m, payload.SessionID); err != nil {
		return errors.Wrap(err, "failed to load model.Session from database")
	}

	su := User()

	// We don't send email if it has been sent before, UNLESS the force
	// flag is specified
	if m.SelectionResultSent {
		if payload.Force {
			if err := su.IsAdministrator(tx, payload.UserID); err != nil {
				return errors.New("must be administrator to force send notification")
			}
		} else {
			return errors.New("selection result has already been sent")
		}
	}

	// Load the user
	var u model.User
	if err := su.Lookup(tx, &u, m.SpeakerID); err != nil {
		return errors.Wrap(err, "failed to load model.User from database")
	}

	// Now, based on the user's language, decorate the session
	if err := v.Decorate(tx, &m, payload.TrustedCall, u.Lang); err != nil {
		return errors.Wrap(err, "failed to declorate mode.Session")
	}

	var subject string
	var tname string
	switch m.Status {
	case model.StatusAccepted:
		subject = "[" + m.Conference.Title + "] Your proposal has been accepted"
		tname = "templates/" + u.Lang + "/eml/proposal-accepted.eml"
	case model.StatusRejected:
		subject = "[" + m.Conference.Title + "] Your proposal was not accepted"
		tname = "templates/" + u.Lang + "/eml/proposal-rejected.eml"
	default:
		return errors.New("can only send email for accepted/rejected sessions")
	}

	t, err := Template().Get(tname)
	if err != nil {
		return errors.Wrap(err, "failed to fetch template")
	}

	tz := time.UTC
	if xtz, err := time.LoadLocation(m.Conference.Timezone); err != nil {
		tz = xtz
	}

	vars := struct {
		Session  *model.Session
		Timezone *time.Location
	}{
		Session:  &m,
		Timezone: tz,
	}

	var msg bytes.Buffer
	if err := t.Execute(&msg, &vars); err != nil {
		return errors.Wrap(err, "failed to render notification template")
	}

	// Record that we have sent this notification
	// We do this BEFORE we actually send the email, because
	// we might FAIL updating the database after sending the email,
	// and that's not cool.
	// Changes to database is only really recorded when Commit is called
	// in the caller
	var req model.UpdateSessionRequest
	req.ID = payload.SessionID
	req.UserID = payload.UserID
	req.SelectionResultSent.Set(true)

	if err := v.UpdateFromPayload(ctx, tx, &req); err != nil {
		return errors.Wrap(err, "failed to update database")
	}

	mm := MailMessage{
		Recipients: []string{m.Speaker.Email},
		Subject:    subject,
		Text:       msg.String(),
	}

	if pdebug.Enabled {
		pdebug.Printf("%#v", mm)
	}

	if !InTesting {
		if err := Mailgun().Send(&mm); err != nil {
			return errors.Wrap(err, "failed to send notification")
		}
	}

	return nil
}

var videoRx = regexp.MustCompile(`^https://youtube.com/watch?v=(.+)`)

func (v *SessionSvc) VideoID(s *model.Session) (string, error) {
	if s.VideoURL == "" {
		return "", errors.New(`video url is not initialized`)
	}

	matches := videoRx.FindStringSubmatch(s.VideoURL)
	if matches == nil {
		return "", errors.New(`could not match video url`)
	}
	return matches[1], nil
}
