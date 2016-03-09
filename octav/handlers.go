package octav

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/builderscon/octav/octav/db"
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

func doCreateConference(ctx context.Context, w http.ResponseWriter, r *http.Request, payload CreateConferenceRequest) {
	c := Conference{
		Slug: payload.Slug,
		Title: payload.Title,
		SubTitle: payload.SubTitle,
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `CreateConference`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	if err := c.Create(tx); err != nil {
		httpError(w, `CreateConference`, http.StatusInternalServerError, err)
		return
	}
	if err := tx.Commit(); err != nil {
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

func doCreateSession(ctx context.Context, w http.ResponseWriter, r *http.Request, payload interface{}) {
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
	if err := vl.Load(tx, payload.Since, payload.Limit); err != nil {
		httpError(w, `ListVenues`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, vl)
}

func doLookupSession(ctx context.Context, w http.ResponseWriter, r *http.Request, payload map[string]interface{}) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `LookupSession`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	s := Session{}
	if err := s.Load(tx, payload["id"].(string)); err != nil {
		httpError(w, `LookupSession`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, s)
}

func doListSessionsByConference(ctx context.Context, w http.ResponseWriter, r *http.Request, payload map[string]interface{}) {
	cid := payload["conference_id"].(string)
	date := payload["date"].(string)

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `ListVenuesByConference`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	sl := SessionList{}
	if err := sl.LoadByConference(tx, cid, date); err != nil {
		httpError(w, `ListVenuesByConference`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, sl)
}
