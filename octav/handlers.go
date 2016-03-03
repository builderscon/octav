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

