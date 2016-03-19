package octav

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/service"
	"github.com/lestrrat/go-pdebug"
	"golang.org/x/net/context"
)

func httpJSON(w http.ResponseWriter, v interface{}) {
	buf := bytes.Buffer{}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		httpError(w, `encode json`, http.StatusInternalServerError, err)
		return
	}
pdebug.Printf("JSON -> '%s'", buf.String())

	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}

func doCreateConference(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.CreateConferenceRequest) {
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
	if err := s.Create(tx, &vdb, payload); err != nil {
		httpError(w, `CreateConference`, http.StatusInternalServerError, err)
		return
	}

	if err := s.AddAdministrator(tx, vdb.EID, payload.UserID); err != nil {
		httpError(w, `CreateConference`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `CreateConference`, http.StatusInternalServerError, err)
		return
	}

	c := model.Conference{}
	if err := c.FromRow(vdb); err != nil {
		httpError(w, `CreateConference`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, c)
}

func doLookupConference(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.LookupConferenceRequest) {
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

	c := model.Conference{}
	if err := c.Load(tx, payload.ID); err != nil {
		httpError(w, `LookupConference`, http.StatusInternalServerError, err)
		return
	}

	s := service.Conference{}
	if payload.Lang.Valid() {
		if err := s.ReplaceL10NStrings(tx, &c, payload.Lang.String); err != nil {
			httpError(w, `LookupConference`, http.StatusInternalServerError, err)
			return
		}
	}

	if err := s.LoadDates(tx, &c.Dates, c.ID); err != nil {
		httpError(w, `LookupConference`, http.StatusInternalServerError, err)
		return
	}

	if err := s.LoadAdmins(tx, &c.Administrators, c.ID); err != nil {
		httpError(w, `LookupConference`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, c)
}

func doUpdateConference(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.UpdateConferenceRequest) {
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

	vdb := db.Conference{}
	if err := vdb.LoadByEID(tx, payload.ID); err != nil {
		httpError(w, `UpdateConference`, http.StatusNotFound, err)
		return
	}

	s := service.Conference{}
	if err := s.Update(tx, &vdb, payload); err != nil {
		httpError(w, `UpdateConference`, http.StatusInternalServerError, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `UpdateConference`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, map[string]string{"status": "success"})
}

func doDeleteConference(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.DeleteConferenceRequest) {
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

	s := service.Conference{}
	if err := s.Delete(tx, payload.ID); err != nil {
		httpError(w, `DeleteConference`, http.StatusInternalServerError, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `DeleteConference`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
}

func doDeleteConferenceDates(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.DeleteConferenceDatesRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `DeleteConferenceDates`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	s := service.Conference{}
	if err := s.DeleteDates(tx, payload.ConferenceID, payload.Dates...); err != nil {
		httpError(w, `DeleteConferenceDates`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `DeleteConferenceDates`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
}

func doAddConferenceDates(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.AddConferenceDatesRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `AddConferenceDates`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	s := service.Conference{}
	if err := s.AddDates(tx, payload.ConferenceID, payload.Dates...); err != nil {
		httpError(w, `AddConferenceDates`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `AddConferenceDates`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
}

func doDeleteConferenceAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.DeleteConferenceAdminRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `DeleteConferenceAdmin`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	s := service.Conference{}
	if err := s.DeleteAdmin(tx, payload.ConferenceID, payload.UserID); err != nil {
		httpError(w, `DeleteConferenceAdmin`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `DeleteConferenceAdmin`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
}

func doAddConferenceAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.AddConferenceAdminRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `AddConferenceAdmin`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	s := service.Conference{}
	if err := s.AddAdmin(tx, payload.ConferenceID, payload.UserID); err != nil {
		httpError(w, `AddConferenceAdmin`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `AddConferenceAdmin`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
}

func doListConferences(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.ListConferencesRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `ListConferences`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	s := service.Conference{}
	vdbl := db.ConferenceList{}
	if err := s.LoadByRange(tx, &vdbl, payload.Since.String, payload.Lang.String, payload.RangeStart.String, payload.RangeEnd.String, int(payload.Limit.Int)); err != nil {
		httpError(w, `ListConferences`, http.StatusInternalServerError, err)
		return
	}

	l := make(model.ConferenceL10NList, len(vdbl))
	for i, vdb := range vdbl {
		v := model.Conference{}
		if err := v.FromRow(vdb); err != nil {
			httpError(w, `ListConferences`, http.StatusInternalServerError, err)
			return
		}
		vl := model.ConferenceL10N{Conference: v}
		if err := vl.LoadLocalizedFields(tx); err != nil {
			httpError(w, `ListConferences`, http.StatusInternalServerError, err)
			return
		}
		l[i] = vl
	}



	httpJSON(w, l)
}

func doCreateRoom(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.CreateRoomRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `CreateRoom`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	s := service.Room{}
	vdb := db.Room{}
	if err := s.Create(tx, &vdb, payload); err != nil {
		httpError(w, `CreateRoom`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `CreateRoom`, http.StatusInternalServerError, err)
		return
	}

	v := model.Room{}
	if err := v.FromRow(vdb); err != nil {
		httpError(w, `CreateRoom`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, v)
}

func doUpdateRoom(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.UpdateRoomRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doUpdateRoom")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `UpdateRoom`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	vdb := db.Room{}
	if err := vdb.LoadByEID(tx, payload.ID); err != nil {
		httpError(w, `UpdateRoom`, http.StatusNotFound, err)
		return
	}

	s := service.Room{}
	if err := s.Update(tx, &vdb, payload); err != nil {
		httpError(w, `UpdateRoom`, http.StatusInternalServerError, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `UpdateRoom`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, map[string]string{"status": "success"})
}

func doCreateSession(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.CreateSessionRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `CreateSession`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	s := service.Session{}
	vdb := db.Session{}
	if err := s.Create(tx, &vdb, payload); err != nil {
		httpError(w, `CreateSession`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `CreateSession`, http.StatusInternalServerError, err)
		return
	}

	v := model.Session{}
	if err := v.FromRow(vdb); err != nil {
		httpError(w, `CreateSession`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, v)
}

func doUpdateSession(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.UpdateSessionRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `UpdateSession`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	s := service.Session{}
	vdb := db.Session{}
	if err := vdb.LoadByEID(tx, payload.ID); err != nil {
		httpError(w, `UpdateConference`, http.StatusNotFound, err)
		return
	}

	// TODO: We must protect the API server from changing important
	// fields like conference_id, speaker_id, room_id, etc from regular
	// users, but allow administrators to do anything they want
	if err := s.Update(tx, &vdb, payload); err != nil {
		httpError(w, `UpdateSession`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `UpdateSession`, http.StatusInternalServerError, err)
		return
	}

	v := model.Session{}
	if err := v.FromRow(vdb); err != nil {
		httpError(w, `UpdateSession`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, v)
}

func doDeleteSession(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.DeleteSessionRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doDeleteSession")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `DeleteSession`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	s := service.Session{}
	if err := s.Delete(tx, payload.ID); err != nil {
		httpError(w, `DeleteSession`, http.StatusInternalServerError, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `DeleteSession`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
}

func doCreateUser(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.CreateUserRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `CreateUser`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	s := service.User{}
	vdb := db.User{}
	if err := s.Create(tx, &vdb, payload); err != nil {
		httpError(w, `CreateUser`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `CreateUser`, http.StatusInternalServerError, err)
		return
	}

	v := model.User{}
	if err := v.FromRow(vdb); err != nil {
		httpError(w, `CreateUser`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, v)
}

func doDeleteUser(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.DeleteUserRequest) {
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

	s := service.User{}
	if err := s.Delete(tx, payload.ID); err != nil {
		httpError(w, `DeleteUser`, http.StatusInternalServerError, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `DeleteUser`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
}

func doLookupUser(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.LookupUserRequest) {
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

	v := model.User{}
	if err := v.Load(tx, payload.ID); err != nil {
		httpError(w, `LookupUser`, http.StatusInternalServerError, err)
		return
	}

	if payload.Lang.Valid() {
		s := service.User{}
		if err := s.ReplaceL10NStrings(tx, &v, payload.Lang.String); err != nil {
			httpError(w, `LookupUser`, http.StatusInternalServerError, err)
			return
		}
	}

	httpJSON(w, v)
}

func doCreateVenue(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.CreateVenueRequest) {
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

	s := service.Venue{}
	vdb := db.Venue{}
	if err := s.Create(tx, &vdb, payload); err != nil {
		httpError(w, `CreateVenue`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `CreateVenue`, http.StatusInternalServerError, err)
		return
	}

	v := model.Venue{}
	if err := v.FromRow(vdb); err != nil {
		httpError(w, `CreateVenue`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, v)
}

func doUpdateUser(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.UpdateUserRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doUpdateUser")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `UpdateUser`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	vdb := db.User{}
	if err := vdb.LoadByEID(tx, payload.ID); err != nil {
		httpError(w, `UpdateUser`, http.StatusNotFound, err)
		return
	}

	s := service.User{}
	if err := s.Update(tx, &vdb, payload); err != nil {
		httpError(w, `UpdateUser`, http.StatusInternalServerError, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `UpdateUser`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, map[string]string{"status": "success"})
}

func doListRooms(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.ListRoomRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `ListRoom`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	rl := model.RoomList{}
	if err := rl.LoadForVenue(tx, payload.VenueID, payload.Since.String, int(payload.Limit.Int)); err != nil {
		httpError(w, `ListRoom`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, rl)
}

func doLookupRoom(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.LookupRoomRequest) {
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

	v := model.Room{}
	if err := v.Load(tx, payload.ID); err != nil {
		httpError(w, `LookupRoom`, http.StatusInternalServerError, err)
		return
	}

	if payload.Lang.Valid() {
		s := service.Room{}
		if err := s.ReplaceL10NStrings(tx, &v, payload.Lang.String); err != nil {
			httpError(w, `LookupRoom`, http.StatusInternalServerError, err)
			return
		}
	}

	httpJSON(w, v)
}

func doDeleteRoom(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.DeleteRoomRequest) {
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

	s := service.Room{}
	if err := s.Delete(tx, payload.ID); err != nil {
		httpError(w, `DeleteRoom`, http.StatusInternalServerError, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `DeleteRoom`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
}

func doDeleteVenue(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.DeleteVenueRequest) {
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

	s := service.Venue{}
	if err := s.Delete(tx, payload.ID); err != nil {
		httpError(w, `DeleteVenue`, http.StatusInternalServerError, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `DeleteVenue`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
}

func doLookupVenue(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.LookupVenueRequest) {
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

	v := model.Venue{}
	if err := v.Load(tx, payload.ID); err != nil {
		httpError(w, `LookupVenue`, http.StatusInternalServerError, err)
		return
	}

	if payload.Lang.Valid() {
		s := service.Venue{}
		if err := s.ReplaceL10NStrings(tx, &v, payload.Lang.String); err != nil {
			httpError(w, `LookupVenue`, http.StatusInternalServerError, err)
			return
		}
	}

	httpJSON(w, v)
}

func doUpdateVenue(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.UpdateVenueRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `UpdateVenue`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	s := service.Venue{}
	vdb := db.Venue{}
	if err := vdb.LoadByEID(tx, payload.ID); err != nil {
		httpError(w, `UpdateConference`, http.StatusNotFound, err)
		return
	}

	// TODO: We must protect the API server from changing important
	// fields like conference_id, speaker_id, room_id, etc from regular
	// users, but allow administrators to do anything they want
	if err := s.Update(tx, &vdb, payload); err != nil {
		httpError(w, `UpdateVenue`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `UpdateVenue`, http.StatusInternalServerError, err)
		return
	}

	v := model.Venue{}
	if err := v.FromRow(vdb); err != nil {
		httpError(w, `UpdateVenue`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, v)
}

func doListVenues(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.ListVenueRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `ListVenues`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	s := service.Venue{}
	vdbl := db.VenueList{}
	if err := s.LoadList(tx, &vdbl, payload.Since.String, int(payload.Limit.Int)); err != nil {
		httpError(w, `ListVenues`, http.StatusInternalServerError, err)
		return
	}

	l := make(model.VenueL10NList, len(vdbl))
	for i, vdb := range vdbl {
		v := model.Venue{}
		if err := v.FromRow(vdb); err != nil {
			httpError(w, `ListVenues`, http.StatusInternalServerError, err)
			return
		}
		vl := model.VenueL10N{Venue: v}
		if err := vl.LoadLocalizedFields(tx); err != nil {
			httpError(w, `ListVenues`, http.StatusInternalServerError, err)
			return
		}
		l[i] = vl
	}
	httpJSON(w, l)
}

func doLookupSession(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.LookupSessionRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `LookupSession`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	v := model.Session{}
	if err := v.Load(tx, payload.ID); err != nil {
		httpError(w, `LookupSession`, http.StatusInternalServerError, err)
		return
	}

	if payload.Lang.Valid() {
		s := service.Session{}
		if err := s.ReplaceL10NStrings(tx, &v, payload.Lang.String); err != nil {
			httpError(w, `LookupSession`, http.StatusInternalServerError, err)
			return
		}
	}

	httpJSON(w, v)
}

func doListSessionsByConference(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.ListSessionsByConferenceRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `ListSessionsByConference`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	s := service.Session{}
	vdbl := db.SessionList{}
	if err := s.LoadByConference(tx, &vdbl, payload.ConferenceID, payload.Date.String); err != nil {
		httpError(w, `ListSessionsByConference`, http.StatusInternalServerError, err)
		return
	}

	l := make(model.SessionL10NList, len(vdbl))
	for i, vdb := range vdbl {
		v := model.Session{}
		if err := v.FromRow(vdb); err != nil {
			httpError(w, `ListSessionsByConference`, http.StatusInternalServerError, err)
			return
		}
		vl := model.SessionL10N{Session: v}
		if err := vl.LoadLocalizedFields(tx); err != nil {
			httpError(w, `ListSessionsByConference`, http.StatusInternalServerError, err)
			return
		}
		l[i] = vl
	}
	httpJSON(w, l)
}
