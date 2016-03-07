package octav

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/builderscon/octav/octav/db"
	"github.com/lestrrat/go-pdebug"
	"golang.org/x/net/context"
)

func httpJSON(w http.ResponseWriter, v interface{}) {
	buf := bytes.Buffer{}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		httpError(w, `encode json`, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}

func httpError(w http.ResponseWriter, message string, err error) {
	if pdebug.Enabled {
		pdebug.Printf("%s: %s", message, err)
	}
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func doCreateConference(ctx context.Context, w http.ResponseWriter, r *http.Request, payload *Conference) {
	c := db.Conference{}
	payload.ID = UUID()
	if err := payload.ToRow(&c); err != nil {
		httpError(w, `CreateConference`, err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `CreateConference`, err)
		return
	}
	defer tx.AutoRollback()

	if err := c.Create(tx); err != nil {
		httpError(w, `CreateConference`, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `CreateConference`, err)
		return
	}

	c2 := Conference{}
	if err := c2.FromRow(c); err != nil {
		httpError(w, `CreateConference`, err)
		return
	}

	httpJSON(w, c2)
}

func doCreateRoom(ctx context.Context, w http.ResponseWriter, r *http.Request, payload *Room) {
	c := db.Room{}
	payload.ID = UUID()
	if err := payload.ToRow(&c); err != nil {
		httpError(w, `CreateRoom`, err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `CreateRoom`, err)
		return
	}
	defer tx.AutoRollback()

	if err := c.Create(tx); err != nil {
		httpError(w, `CreateRoom`, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `CreateRoom`, err)
		return
	}

	c2 := Room{}
	if err := c2.FromRow(c); err != nil {
		httpError(w, `CreateRoom`, err)
		return
	}

	httpJSON(w, c2)
}

func doCreateSession(ctx context.Context, w http.ResponseWriter, r *http.Request, payload interface{}) {
}

func doCreateUser(ctx context.Context, w http.ResponseWriter, r *http.Request, payload interface{}) {
}

func doCreateVenue(ctx context.Context, w http.ResponseWriter, r *http.Request, payload *Venue) {
}

func doListRooms(ctx context.Context, w http.ResponseWriter, r *http.Request, payload map[string]interface{}) {
}

func doDeleteVenue(ctx context.Context, w http.ResponseWriter, r *http.Request, payload map[string]interface{}) {
	id, ok := payload["id"].(string)
	if !ok {
		httpError(w, `doDeleteVenue`, errors.New("invalid id"))
		return
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `doDeleteVenue`, err)
		return
	}
	defer tx.AutoRollback()

	s := db.Venue{EID: id}
	if err := s.Delete(tx); err != nil {
		httpError(w, `doDeleteVenue`, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `doDeleteVenue`, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func doLookupVenue(ctx context.Context, w http.ResponseWriter, r *http.Request, payload map[string]interface{}) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `doListVenues`, err)
		return
	}
	defer tx.AutoRollback()

	s := Venue{}
	if err := s.Load(tx, payload["id"].(string)); err != nil {
		httpError(w, `doLookupVenue`, err)
		return
	}

	httpJSON(w, s)
}
func doListVenues(ctx context.Context, w http.ResponseWriter, r *http.Request, payload map[string]interface{}) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `doListVenues`, err)
		return
	}
	defer tx.AutoRollback()

	vl := VenueList{}
	if err := vl.Load(tx, payload["since"].(string)); err != nil {
		httpError(w, `doListVenues`, err)
		return
	}

	httpJSON(w, vl)
}

func doLookupSession(ctx context.Context, w http.ResponseWriter, r *http.Request, payload map[string]interface{}) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `doListSession`, err)
		return
	}
	defer tx.AutoRollback()

	s := Session{}
	if err := s.Load(tx, payload["id"].(string)); err != nil {
		httpError(w, `doLookupSession`, err)
		return
	}

	httpJSON(w, s)
}

func doListSessionsByConference(ctx context.Context, w http.ResponseWriter, r *http.Request, payload map[string]interface{}) {
	cid := payload["conference_id"].(string)
	date := payload["date"].(string)

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `doListSessionsByConference`, err)
		return
	}
	defer tx.AutoRollback()

	sl := SessionList{}
	if err := sl.LoadByConference(tx, cid, date); err != nil {
		httpError(w, `doListSessionsByConference`, err)
		return
	}

	httpJSON(w, sl)
}
