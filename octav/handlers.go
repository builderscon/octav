package octav

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/service"
	"github.com/lestrrat/go-pdebug"
	"golang.org/x/net/context"
)

func httpJSON(w http.ResponseWriter, v interface{}) {
	buf := bytes.Buffer{}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		httpError(w, `encode json`, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}

func doCreateConference(ctx context.Context, w http.ResponseWriter, r *http.Request, payload service.CreateConferenceRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doCreateConference")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `CreateConference`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	s := service.Conference{}
	vdb := db.Conference{}
	if err := s.Create(tx, payload, &vdb); err != nil {
		httpError(w, `CreateConference`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `CreateConference`, http.StatusInternalServerError, err)
		return
	}

	c := Conference{}
	if err := c.FromRow(vdb); err != nil {
		httpError(w, `CreateConference`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, c)
}

func doLookupConference(ctx context.Context, w http.ResponseWriter, r *http.Request, payload LookupConferenceRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doLookupConference")
		defer g.End()
	}
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `LookupConference`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	s := Conference{}
	if err := s.Load(tx, payload.ID); err != nil {
		httpError(w, `LookupConference`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, s)
}

func doUpdateConference(ctx context.Context, w http.ResponseWriter, r *http.Request, payload UpdateConferenceRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doUpdateConference")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `UpdateConference`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	v := Conference{}
	if err := v.Load(tx, payload.ID); err != nil {
		httpError(w, `UpdateConference`, http.StatusInternalServerError, err)
		return
	}

	if payload.Title.Valid() {
		v.Title = payload.Title.String
	}

	if payload.SubTitle.Valid() {
		v.SubTitle = payload.SubTitle.String
	}

	if payload.Slug.Valid() {
		v.Slug = payload.Slug.String
	}

	payload.L10N.Foreach(v.L10N.Set)

	if err := v.Update(tx); err != nil {
		httpError(w, `UpdateConference`, http.StatusInternalServerError, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `UpdateConference`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
}

func doDeleteConference(ctx context.Context, w http.ResponseWriter, r *http.Request, payload DeleteConferenceRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doDeleteConference")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `DeleteConference`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	v := Conference{ID: payload.ID}
	if err := v.Delete(tx); err != nil {
		httpError(w, `DeleteConference`, http.StatusInternalServerError, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `DeleteConference`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
}

func doCreateRoom(ctx context.Context, w http.ResponseWriter, r *http.Request, payload Room) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `CreateRoom`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	if err := payload.Create(tx); err != nil {
		httpError(w, `CreateRoom`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `CreateRoom`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, payload)
}

func doCreateSession(ctx context.Context, w http.ResponseWriter, r *http.Request, payload service.CreateSessionRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `CreateSession`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	s := service.Session{}
	vdb := db.Session{}
	if err := s.Create(tx, payload, &vdb); err != nil {
		httpError(w, `CreateSession`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `CreateSession`, http.StatusInternalServerError, err)
		return
	}

	v := Session{}
	if err := v.FromRow(vdb); err != nil {
		httpError(w, `CreateSession`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, v)
}

func doUpdateSession(ctx context.Context, w http.ResponseWriter, r *http.Request, payload UpdateSessionRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `UpdateSession`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	v := Session{}
	if err := v.Load(tx, payload.ID); err != nil {
		httpError(w, `UpdateConference`, http.StatusNotFound, err)
		return
	}

	// TODO: We must protect the API server from changing important
	// fields like conference_id, speaker_id, room_id, etc from regular
	// users, but allow administrators to do anything they want
	if payload.ConferenceID.Valid() {
		v.ConferenceID = payload.ConferenceID.String
	}

	if payload.SpeakerID.Valid() {
		v.SpeakerID = payload.SpeakerID.String
	}

	if payload.HasInterpretation.Valid() {
		v.HasInterpretation = payload.HasInterpretation.Bool
	}

	if payload.Status.Valid() {
		v.Status = payload.Status.String
	}

	if payload.SortOrder.Valid() {
		v.SortOrder = int(payload.SortOrder.Int)
	}

	if payload.Confirmed.Valid() {
		v.Confirmed = payload.Confirmed.Bool
	}

	// TODO: End of important stuff that regular users should not be
	// updating on their own

	if payload.Title.Valid() {
		v.Title = payload.Title.String
	}

	if payload.Abstract.Valid() {
		v.Abstract = payload.Abstract.String
	}

	if payload.Memo.Valid() {
		v.Memo = payload.Memo.String
	}

	if payload.Duration.Valid() {
		v.Duration = int(payload.Duration.Int)
	}

	if payload.MaterialLevel.Valid() {
		v.MaterialLevel = payload.MaterialLevel.String
	}

	if payload.Category.Valid() {
		v.Category = payload.Category.String
	}

	if payload.SpokenLanguage.Valid() {
		v.SpokenLanguage = payload.SpokenLanguage.String
	}

	if payload.SlideLanguage.Valid() {
		v.SlideLanguage = payload.SlideLanguage.String
	}

	if payload.SlideSubtitles.Valid() {
		v.SlideSubtitles = payload.SlideSubtitles.String
	}

	if payload.SlideURL.Valid() {
		v.SlideURL = payload.SlideURL.String
	}

	if payload.VideoURL.Valid() {
		v.VideoURL = payload.VideoURL.String
	}

	if payload.PhotoPermission.Valid() {
		v.PhotoPermission = payload.PhotoPermission.String
	}

	if payload.VideoPermission.Valid() {
		v.VideoPermission = payload.VideoPermission.String
	}

	if payload.Tags.Valid() {
		v.Tags = TagString(payload.Tags.String)
	}

	payload.L10N.Foreach(v.L10N.Set)

	if err := v.Update(tx); err != nil {
		httpError(w, `UpdateSession`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `UpdateSession`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, v)
}

func doCreateUser(ctx context.Context, w http.ResponseWriter, r *http.Request, payload CreateUserRequest) {
	u := User{
		Email:      payload.Email,
		FirstName:  payload.FirstName,
		LastName:   payload.LastName,
		Nickname:   payload.Nickname,
		TshirtSize: payload.TshirtSize,
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `CreateUser`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	if err := u.Create(tx); err != nil {
		httpError(w, `CreateUser`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `CreateUser`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, u)
}

func doDeleteUser(ctx context.Context, w http.ResponseWriter, r *http.Request, payload DeleteUserRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doDeleteUser")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `DeleteUser`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	v := User{ID: payload.ID}
	if err := v.Delete(tx); err != nil {
		httpError(w, `DeleteUser`, http.StatusInternalServerError, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `DeleteUser`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
}

func doLookupUser(ctx context.Context, w http.ResponseWriter, r *http.Request, payload LookupUserRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doLookupUser")
		defer g.End()
	}
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `LookupUser`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	s := User{}
	if err := s.Load(tx, payload.ID); err != nil {
		httpError(w, `LookupUser`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, s)
}

func doCreateVenue(ctx context.Context, w http.ResponseWriter, r *http.Request, payload Venue) {
	if pdebug.Enabled {
		g := pdebug.Marker("doCreateVenue")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `CreateVenue`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	if err := payload.Create(tx); err != nil {
		httpError(w, `CreateVenue`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `CreateVenue`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, payload)
}

func doListRooms(ctx context.Context, w http.ResponseWriter, r *http.Request, payload ListRoomRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `ListRoom`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	rl := RoomList{}
	if err := rl.LoadForVenue(tx, payload.VenueID, payload.Since.String, int(payload.Limit.Int)); err != nil {
		httpError(w, `ListRoom`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, rl)
}

func doLookupRoom(ctx context.Context, w http.ResponseWriter, r *http.Request, payload LookupRoomRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doLookupRoom")
		defer g.End()
	}
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `LookupRoom`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	s := Room{}
	if err := s.Load(tx, payload.ID); err != nil {
		httpError(w, `LookupRoom`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, s)
}

func doDeleteRoom(ctx context.Context, w http.ResponseWriter, r *http.Request, payload DeleteRoomRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doDeleteRoom")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `DeleteRoom`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	v := Room{ID: payload.ID}
	if err := v.Delete(tx); err != nil {
		httpError(w, `DeleteRoom`, http.StatusInternalServerError, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `DeleteRoom`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
}

func doDeleteVenue(ctx context.Context, w http.ResponseWriter, r *http.Request, payload DeleteVenueRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doDeleteVenue")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `DeleteVenue`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	v := Venue{ID: payload.ID}
	if err := v.Delete(tx); err != nil {
		httpError(w, `DeleteVenue`, http.StatusInternalServerError, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `DeleteVenue`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
}

func doLookupVenue(ctx context.Context, w http.ResponseWriter, r *http.Request, payload LookupVenueRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doLookupVenue")
		defer g.End()
	}
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `LookupVenue`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	s := Venue{}
	if err := s.Load(tx, payload.ID); err != nil {
		httpError(w, `LookupVenue`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, s)
}

func doListVenues(ctx context.Context, w http.ResponseWriter, r *http.Request, payload ListVenueRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `ListVenues`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	vl := VenueList{}
	if err := vl.Load(tx, payload.Since.String, int(payload.Limit.Int)); err != nil {
		httpError(w, `ListVenues`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, vl)
}

func doLookupSession(ctx context.Context, w http.ResponseWriter, r *http.Request, payload LookupSessionRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `LookupSession`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	s := Session{}
	if err := s.Load(tx, payload.ID); err != nil {
		httpError(w, `LookupSession`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, s)
}

func doListSessionsByConference(ctx context.Context, w http.ResponseWriter, r *http.Request, payload ListSessionsByConferenceRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `ListVenuesByConference`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	sl := SessionList{}
	if err := sl.LoadByConference(tx, payload.ConferenceID, payload.Date.String); err != nil {
		httpError(w, `ListVenuesByConference`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, sl)
}
