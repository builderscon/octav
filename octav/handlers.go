package octav

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/internal/errors"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/service"
	"github.com/lestrrat/go-apache-logformat"
	"github.com/lestrrat/go-pdebug"
	"context"
)

var mwset middlewareSet
type middlewareSet struct{}

func (m middlewareSet) Wrap(h http.Handler) http.Handler {
	l := apachelog.CombinedLog
	return apachelog.WrapLoggingWriter(h, l)
}

const trustedCall = "octav.api.trustedCall"
func isTrustedCall(ctx context.Context) bool {
	allow, ok := ctx.Value(trustedCall).(bool)
	return ok && allow
}

func init() {
	httpError = httpErrorAsJSON
	mwset = middlewareSet{}
}

type httpCoder interface {
	HTTPCode() int
}

func httpCodeFromError(err error) int {
	if v, ok := err.(httpCoder); ok {
		return v.HTTPCode()
	}
	return http.StatusInternalServerError
}

func httpWithOptionalBasicAuth(h HandlerWithContext) HandlerWithContext {
	return wrapBasicAuth(h, true)
}
func httpWithBasicAuth(h HandlerWithContext) HandlerWithContext {
	return wrapBasicAuth(h, false)
}

func wrapBasicAuth(h HandlerWithContext, authIsOptional bool) HandlerWithContext {
	return HandlerWithContext(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		// Verify access token in the Basic-Auth
		clientID, clientSecret, ok := r.BasicAuth()
		if !ok {
			if pdebug.Enabled {
				pdebug.Printf("clientID and/or clientSecret not provided")
			}

			if authIsOptional {
				// if the authentication is optional, then we can just proceed
				if pdebug.Enabled {
					pdebug.Printf("authentication is optional, allowing regular access")
				}
				h(ctx, w, r)
				return
			}
			w.Header().Set("WWW-Authenticate", `Basic realm="octav"`)
			httpError(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized, nil)
			return
		}

		s := service.Client{}
		if err := s.Authenticate(clientID, clientSecret); err != nil {
			if pdebug.Enabled {
				pdebug.Printf("Failed to authenticate client: %s", err)
			}
			code := httpCodeFromError(err)
			httpError(w, http.StatusText(code), code, err)
			return
		}

		if pdebug.Enabled {
			pdebug.Printf("Authentication succeeded, proceeding to call handler")
		}
		ctx = context.WithValue(ctx, trustedCall, true)
		h(ctx, w, r)
	})
}

func httpJSONWithStatus(w http.ResponseWriter, v interface{}, st int) {
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		httpError(w, `encode json`, http.StatusInternalServerError, err)
		return
	}

	if pdebug.Enabled {
		pdebug.Printf("response buffer: %s", buf.String())
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(st)
	buf.WriteTo(w)
}

type jsonerr struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func httpErrorAsJSON(w http.ResponseWriter, message string, st int, err error) {
	v := jsonerr{
		Message: message,
	}
	if err != nil {
		v.Error = err.Error()
	}
	httpJSONWithStatus(w, v, st)
}

func httpJSON(w http.ResponseWriter, v interface{}) {
	httpJSONWithStatus(w, v, http.StatusOK)
}

func doCreateConferenceSeries(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.CreateConferenceSeriesRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doCreateConferenceSeries")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `CreateConferenceSeries`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.ConferenceSeries
	var c model.ConferenceSeries
	if err := s.CreateFromPayload(tx, &c, payload); err != nil {
		httpError(w, `CreateConferenceSeries`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `CreateConferenceSeries`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, c)
}

func doDeleteConferenceSeries(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.DeleteConferenceSeriesRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doDeleteConferenceSeries")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `DeleteConferenceSeries`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.ConferenceSeries
	if err := s.DeleteFromPayload(tx, payload); err != nil {
		httpError(w, `DeleteConferenceSeries`, http.StatusInternalServerError, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `DeleteConferenceSeries`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
}

func doListConferenceSeries(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.ListConferenceSeriesRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `ListConferencesSeries`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.ConferenceSeries
	l := []model.ConferenceSeries{}
	if err := s.LoadByRange(tx, &l, payload.Since.String, int(payload.Limit.Int)); err != nil {
		httpError(w, `ListConferenceSeries`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, l)
}

func doAddConferenceSeriesAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.AddConferenceSeriesAdminRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `AddConferenceSeriesAdmin`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.ConferenceSeries
	if err := s.AddAdministratorFromPayload(tx, payload); err != nil {
		httpError(w, `AddConferenceSeriesAdmin`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `AddConferenceSeriesAdmin`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
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

	var s service.Conference
	var c model.Conference
	if err := s.CreateFromPayload(tx, payload, &c); err != nil {
		httpError(w, `CreateConference`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `CreateConference`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, model.ObjectID{ID: c.ID, Type: "conference"})
}

func doLookupConferenceBySlug(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.LookupConferenceBySlugRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doLookupConferenceBySlug")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `LookupConferenceBySlug`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.Conference
	var c model.Conference
	if err := s.LookupBySlug(tx, &c, payload); err != nil {
		httpError(w, `LookupConferenceBySlug`, http.StatusInternalServerError, err)
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

	var s service.Conference
	var c model.Conference
	if err := s.LookupFromPayload(tx, &c, payload); err != nil {
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

	var s service.Conference
	var updateErr error
	if updateErr = s.UpdateFromPayload(ctx, tx, payload); !errors.IsIgnorable(updateErr) {
		httpError(w, `UpdateConference`, http.StatusInternalServerError, updateErr)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `UpdateConference`, http.StatusInternalServerError, err)
		return
	}

	// This extra bit is for finalizing the image upload
	if cb, ok := errors.IsFinalizationRequired(updateErr); ok {
		if err := cb(); err != nil {
			httpError(w, `Failed to finalize image uploads`, http.StatusInternalServerError, err)
			return
		}
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

	var s service.Conference
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

	var s service.Conference
	if err := s.DeleteDatesFromPayload(tx, payload); err != nil {
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

	var s service.Conference
	if err := s.AddDatesFromPayload(tx, payload); err != nil {
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

	var s service.Conference
	if err := s.DeleteAdministratorFromPayload(tx, payload); err != nil {
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

	var s service.Conference
	if err := s.AddAdministratorFromPayload(tx, payload); err != nil {
		httpError(w, `AddConferenceAdmin`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `AddConferenceAdmin`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
}

func doDeleteConferenceVenue(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.DeleteConferenceVenueRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `DeleteConferenceVenue`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.Conference
	if err := s.DeleteVenueFromPayload(tx, payload); err != nil {
		httpError(w, `DeleteConferenceVenue`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `DeleteConferenceVenue`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
}

func doAddConferenceVenue(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.AddConferenceVenueRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `AddConferenceVenue`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.Conference
	if err := s.AddVenueFromPayload(tx, payload); err != nil {
		httpError(w, `AddConferenceVenue`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `AddConferenceVenue`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
}

func doListConference(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.ListConferenceRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `ListConferences`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.Conference
	var v model.ConferenceList
	if err := s.ListFromPayload(tx, &v, payload); err != nil {
		httpError(w, `ListConference`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, v)
}

func doCreateRoom(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.CreateRoomRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `CreateRoom`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.Room
	var v model.Room
	if err := s.CreateFromPayload(tx, &v, payload); err != nil {
		httpError(w, `CreateRoom`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
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

	var s service.Room
	if err := s.UpdateFromPayload(tx, payload); err != nil {
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

	var s service.Session
	var v model.Session
	if err := s.CreateFromPayload(tx, &v, payload); err != nil {
		httpError(w, `CreateSession`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
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

	var s service.Session
	var v model.Session
	if err := s.UpdateFromPayload(tx, &v, payload); err != nil {
		httpError(w, `UpdateConference`, http.StatusNotFound, err)
		return
	}
	if err := tx.Commit(); err != nil {
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

	var s service.Session
	if err := s.DeleteFromPayload(tx, payload); err != nil {
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

	var s service.User
	var v model.User
	if err := s.CreateFromPayload(tx, &v, payload); err != nil {
		httpError(w, `CreateUser`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
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

	var s service.User
	if err := s.DeleteFromPayload(tx, payload); err != nil {
		httpError(w, `DeleteUser`, http.StatusInternalServerError, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `DeleteUser`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
}

func doListUser(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.ListUserRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `ListUsers`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.User
	var v model.UserList
	payload.TrustedCall = isTrustedCall(ctx)
	if err := s.ListFromPayload(tx, &v, payload); err != nil {
		httpError(w, `ListUsers`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, v)
}

func doLookupUserByAuthUserID(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.LookupUserByAuthUserIDRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doLookupUserByAuthUserID")
		defer g.End()
	}
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `LookupUserByAuthUserID`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.User
	var v model.User
	payload.TrustedCall = isTrustedCall(ctx)
	if err := s.LookupUserByAuthUserIDFromPayload(tx, &v, payload); err != nil {
		httpError(w, `LookupUserByAuthUserID`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, v)
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

	var s service.User
	var v model.User
	payload.TrustedCall = isTrustedCall(ctx)
	if err := s.LookupFromPayload(tx, &v, payload); err != nil {
		httpError(w, `LookupUser`, http.StatusInternalServerError, err)
		return
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

	var s service.Venue
	var v model.Venue
	if err := s.CreateFromPayload(tx, &v, payload); err != nil {
		httpError(w, `CreateVenue`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
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

	s := service.User{}
	if err := s.UpdateFromPayload(tx, payload); err != nil {
		httpError(w, `UpdateUser`, http.StatusInternalServerError, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `UpdateUser`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, map[string]string{"status": "success"})
}

func doListRoom(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.ListRoomRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `ListRoom`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.Room
	var v model.RoomList
	if err := s.ListFromPayload(tx, &v, payload); err != nil {
		httpError(w, `ListRoom`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, v)
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

	var s service.Room
	var v model.Room
	if err := s.LookupFromPayload(tx, &v, payload); err != nil {
		httpError(w, `LookupRoom`, http.StatusInternalServerError, err)
		return
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

	var s service.Room
	if err := s.DeleteFromPayload(tx, payload); err != nil {
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

	var s service.Venue
	if err := s.DeleteFromPayload(tx, payload); err != nil {
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

	var s service.Venue
	var v model.Venue
	if err := s.LookupFromPayload(tx, &v, payload); err != nil {
		httpError(w, `LookupVenue`, http.StatusInternalServerError, err)
		return
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

	var s service.Venue
	var v model.Venue
	if err := s.UpdateFromPayload(tx, &v, payload); err != nil {
		httpError(w, `UpdateVenue`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `UpdateVenue`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, v)
}

func doListVenue(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.ListVenueRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `ListVenues`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.Venue
	var v model.VenueList
	if err := s.ListFromPayload(tx, &v, payload); err != nil {
		httpError(w, `ListVenues`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, v)
}

func doLookupSession(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.LookupSessionRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doLookupSession")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `LookupSession`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.Session
	var v model.Session
	payload.TrustedCall = isTrustedCall(ctx)
	if err := s.LookupFromPayload(tx, &v, payload); err != nil {
		httpError(w, `LookupSession`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, v)
}

func doListSessions(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.ListSessionsRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `ListSessions`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.Session
	var v model.SessionList
	if err := s.ListSessionFromPayload(tx, &v, payload); err != nil {
		httpError(w, `ListSessions`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, v)
}

func doCreateQuestion(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.CreateQuestionRequest) {
}

func doDeleteQuestion(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.DeleteQuestionRequest) {
}

func doListQuestion(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.ListQuestionRequest) {
}

func doCreateSessionSurveyResponse(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.CreateSessionSurveyResponseRequest) {

}

func doAddFeaturedSpeaker(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.AddFeaturedSpeakerRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doAddFeaturedSpeaker")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `AddFeaturedSpeaker`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.FeaturedSpeaker
	var c model.FeaturedSpeaker
	if err := s.CreateFromPayload(tx, payload, &c); err != nil {
		httpError(w, `AddFeaturedSpeaker`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `CreateConference`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, c)
}

func doDeleteFeaturedSpeaker(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.DeleteFeaturedSpeakerRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doDeleteFeaturedSpeaker")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `DeleteFeaturedSpeaker`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.FeaturedSpeaker
	if err := s.DeleteFromPayload(tx, payload); err != nil {
		httpError(w, `DeleteFeaturedSpeaker`, http.StatusInternalServerError, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `DeleteFeaturedSpeaker`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
}

func doUpdateFeaturedSpeaker(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.UpdateFeaturedSpeakerRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doUpdateFeaturedSpeaker")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `UpdateFeaturedSpeaker`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.FeaturedSpeaker
	if err := s.UpdateFromPayload(tx, payload); err != nil {
		httpError(w, `UpdateFeaturedSpeaker`, http.StatusInternalServerError, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `UpdateFeaturedSpeaker`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, map[string]string{"status": "success"})
}

func doLookupFeaturedSpeaker(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.LookupFeaturedSpeakerRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doLookupFeaturedSpeaker")
		defer g.End()
	}
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `LookupFeaturedSpeaker`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.FeaturedSpeaker
	var c model.FeaturedSpeaker
	if err := s.LookupFromPayload(tx, &c, payload); err != nil {
		httpError(w, `LookupFeaturedSpeaker`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, c)
}

func doListFeaturedSpeakers(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.ListFeaturedSpeakersRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `ListConferencesSeries`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.FeaturedSpeaker
	var l model.FeaturedSpeakerList
	if err := s.ListFromPayload(tx, &l, payload); err != nil {
		httpError(w, `ListFeaturedSpeakers`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, l)
}

func doAddSponsor(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.AddSponsorRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doAddSponsor")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `AddSponsor`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.Sponsor
	var c model.Sponsor

	if err := s.CreateFromPayload(ctx, tx, payload, &c); err != nil {
		httpError(w, `Faild to create sponsor from payload`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `Failed to commit transaction`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, c)
}

func doDeleteSponsor(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.DeleteSponsorRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doDeleteSponsor")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `DeleteSponsor`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.Sponsor
	if err := s.DeleteFromPayload(ctx, tx, payload); err != nil {
		httpError(w, `DeleteSponsor`, http.StatusInternalServerError, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `DeleteSponsor`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
}

func doUpdateSponsor(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.UpdateSponsorRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doUpdateSponsor")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `UpdateSponsor`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.Sponsor
	updateErr := s.UpdateFromPayload(ctx, tx, payload)
	if !errors.IsIgnorable(updateErr) {
		httpError(w, `Failed to update data from payload`, http.StatusInternalServerError, updateErr)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `Failed to commit data`, http.StatusInternalServerError, err)
		return
	}

	// This extra bit is for finalizing the image upload
	if cb, ok := errors.IsFinalizationRequired(updateErr); ok {
		if err := cb(); err != nil {
			httpError(w, `Failed to finalize image uploads`, http.StatusInternalServerError, err)
			return
		}
	}

	httpJSON(w, map[string]string{"status": "success"})
}

func doLookupSponsor(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.LookupSponsorRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doLookupSponsor")
		defer g.End()
	}
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `LookupSponsor`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.Sponsor
	var c model.Sponsor
	if err := s.LookupFromPayload(tx, &c, payload); err != nil {
		httpError(w, `LookupSponsor`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, c)
}

func doListSponsors(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.ListSponsorsRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `ListConferencesSeries`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.Sponsor
	var l model.SponsorList
	if err := s.ListFromPayload(tx, &l, payload); err != nil {
		httpError(w, `ListSponsors`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, l)
}

func doAddSessionType(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.AddSessionTypeRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doAddSessionType")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `AddSessionType`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.SessionType
	var c model.SessionType

	if err := s.CreateFromPayload(ctx, tx, payload, &c); err != nil {
		httpError(w, `Faild to create sponsor from payload`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `Failed to commit transaction`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, c)
}

func doDeleteSessionType(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.DeleteSessionTypeRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doDeleteSessionType")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `DeleteSessionType`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.SessionType
	if err := s.DeleteFromPayload(ctx, tx, payload); err != nil {
		httpError(w, `DeleteSessionType`, http.StatusInternalServerError, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, `DeleteSessionType`, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, map[string]string{"status": "success"})
}

func doUpdateSessionType(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.UpdateSessionTypeRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doUpdateSessionType")
		defer g.End()
	}

	tx, err := db.Begin()
	if err != nil {
		httpError(w, `UpdateSessionType`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.SessionType
	if err := s.UpdateFromPayload(ctx, tx, payload); err != nil {
		httpError(w, `Failed to update data from payload`, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		httpError(w, `Failed to commit data`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, map[string]string{"status": "success"})
}

func doLookupSessionType(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.LookupSessionTypeRequest) {
	if pdebug.Enabled {
		g := pdebug.Marker("doLookupSessionType")
		defer g.End()
	}
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `LookupSessionType`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.SessionType
	var c model.SessionType
	if err := s.LookupFromPayload(tx, &c, payload); err != nil {
		httpError(w, `LookupSessionType`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, c)
}

func doListSessionTypesByConference(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.ListSessionTypesByConferenceRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `doListSessionTypesbyConference`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.SessionType
	var l model.SessionTypeList
	if err := s.ListFromPayload(tx, &l, payload); err != nil {
		httpError(w, `doListSessionTypesbyConference`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, l)
}

func doListConferencesByOrganizer(ctx context.Context, w http.ResponseWriter, r *http.Request, payload model.ListConferencesByOrganizerRequest) {
	tx, err := db.Begin()
	if err != nil {
		httpError(w, `doListConferencesByOrganizer`, http.StatusInternalServerError, err)
		return
	}
	defer tx.AutoRollback()

	var s service.Conference
	var l model.ConferenceList
	if err := s.ListByOrganizerFromPayload(tx, &l, payload); err != nil {
		httpError(w, `doListConferencesByOrganizer`, http.StatusInternalServerError, err)
		return
	}

	httpJSON(w, l)
}


