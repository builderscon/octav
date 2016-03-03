package octav

import (
	"net/http"

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

func doListRooms(ctx context.Context, w http.ResponseWriter, r *http.Request, payload interface{}) {
}

func doListVenues(ctx context.Context, w http.ResponseWriter, r *http.Request, payload interface{}) {
}
