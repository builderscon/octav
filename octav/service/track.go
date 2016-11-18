package service

import (
	"time"

	"github.com/builderscon/octav/octav/cache"
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	pdebug "github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

func (v *TrackSvc) Init() {}

func (v *TrackSvc) populateRowForCreate(vdb *db.Track, payload *model.CreateTrackRequest) error {
	vdb.EID = tools.UUID()
	vdb.ConferenceID = payload.ConferenceID
	vdb.RoomID = payload.RoomID
	if payload.Name.Valid() {
		vdb.Name = payload.Name.String
	}
	return nil
}

func (v *TrackSvc) populateRowForUpdate(vdb *db.Track, payload *model.UpdateTrackRequest) error {
	if payload.Name.Valid() {
		vdb.Name = payload.Name.String
	}
	return nil
}

func (v *TrackSvc) LookupByConferenceRoom(tx *db.Tx, m *model.Track, conferenceID, roomID string) (err error) {
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

func (v *TrackSvc) CreateFromPayload(tx *db.Tx, payload *model.CreateTrackRequest, result *model.Track) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Track.CreateFromPayload").BindError(&err)
		defer g.End()
	}

	var vdb db.Track
	if err := v.Create(tx, &vdb, payload); err != nil {
		return errors.Wrap(err, "failed to store in database")
	}

	var c model.Track
	if err := c.FromRow(&vdb); err != nil {
		return errors.Wrap(err, "failed to populate model")
	}

	if result != nil {
		*result = c
	}

	invalidateTrackLoadByConferenceID(payload.ConferenceID)
	return nil
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

func (v *TrackSvc) LoadByConferenceID(tx *db.Tx, result *model.TrackList, conferenceID string) (err error) {
	c := Cache()
	key := c.Key("Track", "LoadByConferenceID", conferenceID)

	var ids []string
	if err := c.Get(key, &ids); err == nil {
		if pdebug.Enabled {
			pdebug.Printf("CACHE HIT: %s", key)
		}

		m := make(model.TrackList, len(ids))
		for i, id := range ids {
			if err := v.LookupByConferenceRoom(tx, &m[i], conferenceID, id); err != nil {
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
		ids[i] = vdb.RoomID
		res[i] = u
	}
	*result = res

	c.Set(key, ids, cache.WithExpires(15*time.Minute))
	return nil
}

func (v *TrackSvc) Decorate(tx *db.Tx, track *model.Track, trustedCall bool, lang string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Track.Decorate (%s, %s, %t, %s)", track.ConferenceID, track.RoomID, trustedCall, lang).BindError(&err)
		defer g.End()
	}

	if lang != "" {
		if err := v.ReplaceL10NStrings(tx, track, lang); err != nil {
			return errors.Wrap(err, "failed to replace L10N strings")
		}
	}
	return nil
}
