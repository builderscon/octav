package octav

import (
	"encoding/json"
	"net/http"

	"github.com/builderscon/octav/octav/db"
	"github.com/lestrrat/go-pdebug"
	"golang.org/x/net/context"
)

func doCreateConference(ctx context.Context, w http.ResponseWriter, r *http.Request, payload interface{}) {
}

func doCreateRoom(ctx context.Context, w http.ResponseWriter, r *http.Request, payload interface{}) {
}

func doCreateSession(ctx context.Context, w http.ResponseWriter, r *http.Request, payload interface{}) {
}

func doCreateUser(ctx context.Context, w http.ResponseWriter, r *http.Request, payload interface{}) {
}

func doCreateVenue(ctx context.Context, w http.ResponseWriter, r *http.Request, payload *Venue) {
}

func doListRooms(ctx context.Context, w http.ResponseWriter, r *http.Request, payload map[string]interface{}) {
}

func httpError(w http.ResponseWriter, message string, err error) {
	if pdebug.Enabled {
		pdebug.Printf("%s: %s", message, err)
	}
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func doListVenues(ctx context.Context, w http.ResponseWriter, r *http.Request, payload map[string]interface{}) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `doListVenues`, err)
		return
	}

	var since uint64
	if id := payload["since"].(string); id != "" {
		v := db.Venue{}
		if err := v.LoadByEID(tx, id); err != nil {
			httpError(w, `doListVenues`, err)
			return
		}

		since = v.OID
	}

	rows, err := tx.Query(`SELECT eid, name FROM venue WHERE oid > ? ORDER BY oid LIMIT 10`, since)
	if err != nil {
		httpError(w, `doListVenues`, err)
		return
	}

	// Not using db.Venue here
	res := make([]Venue, 0, 10)
	for rows.Next() {
		v := Venue{}
		if err := rows.Scan(&v.ID, &v.Name); err != nil {
			if pdebug.Enabled {
				pdebug.Printf("doListVenues: %s", err)
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		res = append(res, v)
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		httpError(w, `doListVenues`, err)
		return
	}
}
