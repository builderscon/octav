package service

import (
	"database/sql"
	"time"

	"github.com/builderscon/octav/octav/cache"
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/internal/context"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	pdebug "github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

func (v *TrackSvc) Init() {}

func (v *TrackSvc) populateRowForCreate(ctx context.Context, vdb *db.Track, payload *model.CreateTrackRequest) error {
	vdb.EID = tools.UUID()
	vdb.ConferenceID = payload.ConferenceID
	vdb.RoomID = payload.RoomID
	if payload.Name.Valid() {
		vdb.Name = payload.Name.String
	}
	if payload.SortOrder.Valid() {
		vdb.SortOrder = int(payload.SortOrder.Int)
	}
	return nil
}

func (v *TrackSvc) populateRowForUpdate(ctx context.Context, vdb *db.Track, payload *model.UpdateTrackRequest) error {
	if payload.Name.Valid() {
		vdb.Name = payload.Name.String
	}
	if payload.SortOrder.Valid() {
		vdb.SortOrder = int(payload.SortOrder.Int)
	}
	return nil
}

func (v *TrackSvc) LookupByConferenceRoom(tx *sql.Tx, m *model.Track, conferenceID, roomID string) (err error) {
	var r model.Track
	c := Cache()
	key := c.Key("Track", conferenceID, roomID)
	var cacheMiss bool
	_, err = c.GetOrSet(key, &r, func() (interface{}, error) {
		if pdebug.Enabled {
			cacheMiss = true
		}
		if err := r.LoadByConferenceRoom(tx, conferenceID, roomID); err != nil {
			return nil, errors.Wrap(err, "failed to load from database")
		}
		return &r, nil
	}, cache.WithExpires(time.Hour))

	if pdebug.Enabled {
		cacheSt := `HIT`
		if cacheMiss {
			cacheSt = `MISS`
		}
		pdebug.Printf(`CACHE %s: %s`, cacheSt, key)
	}
	*m = r

	return nil
}

func (v *TrackSvc) CreateFromPayload(ctx context.Context, tx *sql.Tx, payload *model.CreateTrackRequest, result *model.Track) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Track.CreateFromPayload").BindError(&err)
		defer g.End()
	}

	su := User()
	if err := su.IsConferenceAdministrator(ctx, tx, payload.ConferenceID, context.GetUserID(ctx)); err != nil {
		return errors.Wrap(err, "creating a track requires conference administrator privilege")
	}

	// If the payload name doesn't exist, we must populate it
	// using the default room name
	if !payload.Name.Valid() || payload.Name.String == "" {
		var m model.Room
		sr := Room()
		if err := sr.Lookup(ctx, tx, &m, payload.RoomID); err != nil {
			return errors.Wrap(err, "failed to load room")
		}
		payload.Name.Set(m.Name)
	}

	var vdb db.Track
	if err := v.Create(ctx, tx, &vdb, payload); err != nil {
		return errors.Wrap(err, "failed to store in database")
	}

	if result != nil {
		var c model.Track
		if err := c.FromRow(&vdb); err != nil {
			return errors.Wrap(err, "failed to populate model")
		}
		*result = c
	}

	return nil
}

func (v *TrackSvc) DeleteFromPayload(ctx context.Context, tx *sql.Tx, payload *model.DeleteTrackRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Track.DeleteFromPayload").BindError(&err)
		defer g.End()
	}

	var vdb db.Track
	if err := vdb.LoadByEID(tx, payload.ID); err != nil {
		return errors.Wrap(err, "failed to load track")
	}

	su := User()
	if err := su.IsConferenceAdministrator(ctx, tx, vdb.ConferenceID, context.GetUserID(ctx)); err != nil {
		return errors.Wrap(err, "deleting a track requires conference administrator privilege")
	}

	if err := v.Delete(tx, payload.ID); err != nil {
		return errors.Wrap(err, "failed to store in database")
	}
	return nil
}

func (v *TrackSvc) PostCreateHook(ctx context.Context, _ *sql.Tx, vdb *db.Track) error {
	return invalidateTrackLoadByConferenceID(vdb.ConferenceID)
}

func (v *TrackSvc) PostUpdateHook(_ *sql.Tx, vdb *db.Track) error {
	return invalidateTrackLoadByConferenceID(vdb.ConferenceID)
}

func (v *TrackSvc) PostDeleteHook(_ *sql.Tx, vdb *db.Track) error {
	return invalidateTrackLoadByConferenceID(vdb.ConferenceID)
}

func invalidateTrackLoadByConferenceID(conferenceID string) error {
	c := Cache()
	key := c.Key("Track", "LoadByConferenceID", conferenceID)
	c.Delete(key)
	if pdebug.Enabled {
		pdebug.Printf("CACHE DEL: %s", key)
	}
	return nil
}

func (v *TrackSvc) LoadByConferenceID(ctx context.Context, tx *sql.Tx, result *model.TrackList, conferenceID string) (err error) {
	c := Cache()
	key := c.Key("Track", "LoadByConferenceID", conferenceID)

	var ids []string
	if err := c.Get(key, &ids); err == nil {
		if pdebug.Enabled {
			pdebug.Printf("CACHE HIT: %s", key)
		}

		m := make(model.TrackList, len(ids))
		for i, id := range ids {
			if err := v.Lookup(ctx, tx, &m[i], id); err != nil {
				// Something fishy. Thro away this cache
				c.Delete(key)
				return errors.Wrap(err, "failed to load from database")
			}
		}

		*result = m
		return nil
	}

	if pdebug.Enabled {
		pdebug.Printf("CACHE MISS: %s", key)
	}

	var vdbl db.TrackList
	if err := vdbl.LoadByConferenceID(tx, conferenceID); err != nil {
		return err
	}

	ids = make([]string, len(vdbl))
	res := make(model.TrackList, len(vdbl))
	for i, vdb := range vdbl {
		var u model.Track
		if err := u.FromRow(&vdb); err != nil {
			return err
		}
		ids[i] = vdb.EID
		res[i] = u
	}
	*result = res

	c.Set(key, ids, cache.WithExpires(15*time.Minute))
	return nil
}

func (v *TrackSvc) Decorate(ctx context.Context, tx *sql.Tx, track *model.Track, verifiedCall bool, lang string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Track.Decorate (%s, %s, %t, %s)", track.ConferenceID, track.RoomID, verifiedCall, lang).BindError(&err)
		defer g.End()
	}

	if lang != "" {
		if err := v.ReplaceL10NStrings(tx, track, lang); err != nil {
			return errors.Wrap(err, "failed to replace L10N strings")
		}
	}
	return nil
}
