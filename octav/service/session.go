package service

import (
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	"github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

func (v *Session) populateRowForCreate(vdb *db.Session, payload model.CreateSessionRequest) error {
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

	if s, ok := payload.L10N.Get("ja", "title"); ok && s != ""  {
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

	if s, ok := payload.L10N.Get("ja", "abstract"); ok && s != ""  {
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

	if payload.PhotoPermission.Valid() {
		vdb.PhotoPermission.Valid = true
		vdb.PhotoPermission.String = payload.PhotoPermission.String
	}

	if payload.VideoPermission.Valid() {
		vdb.VideoPermission.Valid = true
		vdb.VideoPermission.String = payload.VideoPermission.String
	}

	if payload.Tags.Valid() {
		vdb.Tags.Valid = true
		vdb.Tags.String = string(payload.Tags.String)
	}

	return nil
}

func (v *Session) populateRowForUpdate(vdb *db.Session, payload model.UpdateSessionRequest) error {
	if payload.ConferenceID.Valid() {
		vdb.ConferenceID = payload.ConferenceID.String
	}

	if payload.SpeakerID.Valid() {
		vdb.SpeakerID = payload.SpeakerID.String
	}

	if payload.SessionTypeID.Valid() {
		vdb.SessionTypeID = payload.SessionTypeID.String
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

	if payload.PhotoPermission.Valid() {
		vdb.PhotoPermission.Valid = true
		vdb.PhotoPermission.String = payload.PhotoPermission.String
	}

	if payload.VideoPermission.Valid() {
		vdb.VideoPermission.Valid = true
		vdb.VideoPermission.String = payload.VideoPermission.String
	}

	if payload.Tags.Valid() {
		vdb.Tags.Valid = true
		vdb.Tags.String = string(payload.Tags.String)
	}

	return nil
}

func (v *Session) LoadByConference(tx *db.Tx, vdbl *db.SessionList, cid string, date string) error {
	if err := vdbl.LoadByConference(tx, cid, "", date, ""); err != nil {
		return err
	}
	return nil
}

func (v *Session) Decorate(tx *db.Tx, session *model.Session, lang string) error {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Session.Decorate")
		defer g.End()
	}
	// session must be associated with a conference
	if session.ConferenceID != "" {
		var cs Conference
		var mc model.Conference
		if err := cs.Lookup(tx, &mc, session.ConferenceID); err != nil {
			return errors.Wrap(err, "failed to load conference")
		}
		if err := cs.Decorate(tx, &mc, lang); err != nil {
			return errors.Wrap(err, "failed to decorate conference")
		}
		session.Conference = &mc
	}

	// ... but not necessarily with a room
	if session.RoomID != "" {
		var room model.Room
		if err := room.Load(tx, session.RoomID); err != nil {
			return errors.Wrap(err, "failed to load room")
		}
		session.Room = &room
	}

	if session.SpeakerID != "" {
		var speaker model.User
		if err := speaker.Load(tx, session.SpeakerID); err != nil {
			return errors.Wrapf(err, "failed to load speaker '%s'", session.SpeakerID)
		}
		session.Speaker = &speaker
	}

	if session.SessionTypeID != "" {
		var sessionType model.SessionType
		if err := sessionType.Load(tx, session.SessionTypeID); err != nil {
			return errors.Wrapf(err, "failed to load session type '%s'", session.SessionTypeID)
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

func (v *Session) CreateFromPayload(tx *db.Tx, result *model.Session, payload model.CreateSessionRequest) error {
	var u model.User
	su := User{}
	if err := su.Lookup(tx, &u, payload.UserID); err != nil {
		return errors.Wrapf(err, "failed to load user %s", payload.UserID)
	}

	// Check if this session type is allowed to be submitted right now
	sst := SessionType{}
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
	if err := m.FromRow(vdb); err != nil {
		return errors.Wrap(err, "failed to populate model from database")
	}

	*result = m
	return nil
}

func (v *Session) UpdateFromPayload(tx *db.Tx, result *model.Session, payload model.UpdateSessionRequest) error {
	su := User{}
	if err := su.IsSessionOwner(tx, payload.ID, payload.UserID); err != nil {
		return errors.Wrap(err, "updating sessions require session owner privileges")
	}

	vdb := db.Session{}
	if err := vdb.LoadByEID(tx, payload.ID); err != nil {
		return errors.Wrap(err, "failed to load from database")
	}

	// TODO: We must protect the API server from changing important
	// fields like conference_id, speaker_id, room_id, etc from regular
	// users, but allow administrators to do anything they want
	if err := v.Update(tx, &vdb, payload); err != nil {
		return errors.Wrap(err, "failed to update database")
	}

	m := model.Session{}
	if err := m.FromRow(vdb); err != nil {
		return errors.Wrap(err, "failed to populate model from database")
	}

	*result = m
	return nil
}

func (v *Session) ListSessionFromPayload(tx *db.Tx, result *model.SessionList, payload model.ListSessionsRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Session.ListSessionFromPayload").BindError(&err)
		defer g.End()
	}

	// Make sure that we have at least one of the arguments
	var conferenceID, speakerID, date, status string
	var hasQuery bool
	if payload.ConferenceID.Valid() {
		conferenceID = payload.ConferenceID.String
		hasQuery = true
	}

	if payload.SpeakerID.Valid() {
		speakerID = payload.SpeakerID.String
		hasQuery = true
	}

	if payload.Date.Valid() {
		date = payload.Date.String
		// Don't set the hasQuery flag, as this alone doesn't work
	}

	if payload.Status.Valid() {
		status = payload.Status.String
		// Don't set the hasQuery flag, as this alone doesn't work
	}

	if !hasQuery {
		return errors.New("no query specified (one of conference_id/speaker_id is required)")
	}

	var vdbl db.SessionList
	if err := vdbl.LoadByConference(tx, conferenceID, speakerID, date, status); err != nil {
		return errors.Wrap(err, "failed to load from database")
	}

	l := make(model.SessionList, len(vdbl))
	for i, vdb := range vdbl {
		if err := l[i].FromRow(vdb); err != nil {
			return errors.Wrap(err, "failed to populate model from database")
		}

		if err := v.Decorate(tx, &l[i], payload.Lang.String); err != nil {
			return errors.Wrap(err, "failed to decorate session with associated data")
		}
	}

	*result = l
	return nil
}
