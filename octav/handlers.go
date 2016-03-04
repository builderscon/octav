package octav

import (
	"encoding/json"
	"net/http"

	"github.com/builderscon/octav/octav/db"
	"github.com/lestrrat/go-pdebug"
	"golang.org/x/net/context"
)

func httpError(w http.ResponseWriter, message string, err error) {
	if pdebug.Enabled {
		pdebug.Printf("%s: %s", message, err)
	}
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func doCreateConference(ctx context.Context, w http.ResponseWriter, r *http.Request, payload *Conference) {
	c := db.Conference{}
	payload.ID = uuid()
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

	if err := json.NewEncoder(w).Encode(c2); err != nil {
		httpError(w, `CreateConference`, err)
		return
	}
}

func doCreateRoom(ctx context.Context, w http.ResponseWriter, r *http.Request, payload *Room) {
	c := db.Room{}
	payload.ID = uuid()
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

	if err := json.NewEncoder(w).Encode(c2); err != nil {
		httpError(w, `CreateRoom`, err)
		return
	}
}

func doCreateSession(ctx context.Context, w http.ResponseWriter, r *http.Request, payload interface{}) {
}

func doCreateUser(ctx context.Context, w http.ResponseWriter, r *http.Request, payload interface{}) {
}

func doCreateVenue(ctx context.Context, w http.ResponseWriter, r *http.Request, payload *Venue) {
}

func doListRooms(ctx context.Context, w http.ResponseWriter, r *http.Request, payload map[string]interface{}) {
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

	if err := json.NewEncoder(w).Encode(s); err != nil {
		httpError(w, `doLookupVenue`, err)
		return
	}
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

	if err := json.NewEncoder(w).Encode(vl); err != nil {
		httpError(w, `doListVenues`, err)
		return
	}
}

func doLookupSession(ctx context.Context, w http.ResponseWriter, r *http.Request, payload map[string]interface{}) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `doListVenues`, err)
		return
	}
	defer tx.AutoRollback()

	s := Session{}
	if err := s.Load(tx, payload["id"].(string)); err != nil {
		httpError(w, `doLookupSession`, err)
		return
	}

	if err := json.NewEncoder(w).Encode(s); err != nil {
		httpError(w, `doLookupSession`, err)
		return
	}
}
