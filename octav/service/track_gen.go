package service

// Automatically generated by genmodel utility. DO NOT EDIT!

import (
	"context"
	"sync"
	"time"

	"github.com/builderscon/octav/octav/cache"

	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/internal/errors"
	"github.com/builderscon/octav/octav/model"
	"github.com/lestrrat/go-pdebug"
)

var _ = time.Time{}
var _ = cache.WithExpires(time.Minute)
var _ = context.Background
var _ = errors.Wrap
var _ = model.Track{}
var _ = db.Track{}
var _ = pdebug.Enabled

var trackSvc TrackSvc
var trackOnce sync.Once

func Track() *TrackSvc {
	trackOnce.Do(trackSvc.Init)
	return &trackSvc
}

func (v *TrackSvc) LookupFromPayload(tx *db.Tx, m *model.Track, payload *model.LookupTrackRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Track.LookupFromPayload").BindError(&err)
		defer g.End()
	}
	if err = v.Lookup(tx, m, payload.ID); err != nil {
		return errors.Wrap(err, "failed to load model.Track from database")
	}
	if err := v.Decorate(tx, m, payload.TrustedCall, payload.Lang.String); err != nil {
		return errors.Wrap(err, "failed to load associated data for model.Track from database")
	}
	return nil
}

func (v *TrackSvc) Lookup(tx *db.Tx, m *model.Track, id string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Track.Lookup").BindError(&err)
		defer g.End()
	}

	var r model.Track
	c := Cache()
	key := c.Key("Track", id)
	var cacheMiss bool
	_, err = c.GetOrSet(key, &r, func() (interface{}, error) {
		if pdebug.Enabled {
			cacheMiss = true
		}
		if err := r.Load(tx, id); err != nil {
			return nil, errors.Wrap(err, "failed to load model.Track from database")
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

// Create takes in the transaction, the incoming payload, and a reference to
// a database row. The database row is initialized/populated so that the
// caller can use it afterwards.
func (v *TrackSvc) Create(tx *db.Tx, vdb *db.Track, payload *model.CreateTrackRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Track.Create").BindError(&err)
		defer g.End()
	}

	if err := v.populateRowForCreate(vdb, payload); err != nil {
		return errors.Wrap(err, `failed to populate row`)
	}

	if err := vdb.Create(tx, payload.DatabaseOptions...); err != nil {
		return errors.Wrap(err, `failed to insert into database`)
	}

	if err := payload.LocalizedFields.CreateLocalizedStrings(tx, "Track", vdb.EID); err != nil {
		return errors.Wrap(err, `failed to populate localized strings`)
	}
	if err := v.PostCreateHook(tx, vdb); err != nil {
		return errors.Wrap(err, `post create hook failed`)
	}
	return nil
}

func (v *TrackSvc) Update(tx *db.Tx, vdb *db.Track) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Track.Update (%s)", vdb.EID).BindError(&err)
		defer g.End()
	}

	if vdb.EID == `` {
		return errors.New("vdb.EID is required (did you forget to call vdb.Load(tx) before hand?)")
	}

	if err := vdb.Update(tx); err != nil {
		return errors.Wrap(err, `failed to update database`)
	}
	c := Cache()
	key := c.Key("Track", vdb.EID)
	if pdebug.Enabled {
		pdebug.Printf(`CACHE DEL %s`, key)
	}
	cerr := c.Delete(key)
	if pdebug.Enabled {
		if cerr != nil {
			pdebug.Printf(`CACHE ERR: %s`, cerr)
		}
	}
	if err := v.PostUpdateHook(tx, vdb); err != nil {
		return errors.Wrap(err, `post update hook failed`)
	}
	return nil
}

func (v *TrackSvc) UpdateFromPayload(ctx context.Context, tx *db.Tx, payload *model.UpdateTrackRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Track.UpdateFromPayload (%s)", payload.ID).BindError(&err)
		defer g.End()
	}
	var vdb db.Track
	if err := vdb.LoadByEID(tx, payload.ID); err != nil {
		return errors.Wrap(err, `failed to load from database`)
	}

	if err := v.populateRowForUpdate(&vdb, payload); err != nil {
		return errors.Wrap(err, `failed to populate row data`)
	}

	if err := v.Update(tx, &vdb); err != nil {
		return errors.Wrap(err, `failed to update row in database`)
	}

	ls := LocalizedString()
	if err := ls.UpdateFields(tx, "Track", vdb.EID, payload.LocalizedFields); err != nil {
		return errors.Wrap(err, `failed to update localized fields`)
	}
	return nil
}

func (v *TrackSvc) ReplaceL10NStrings(tx *db.Tx, m *model.Track, lang string) error {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Track.ReplaceL10NStrings lang = %s", lang)
		defer g.End()
	}
	ls := LocalizedString()
	list := make([]db.LocalizedString, 0, 1)
	switch lang {
	case "", "en":
		if len(m.Name) > 0 {
			return nil
		}
		for _, extralang := range []string{`ja`} {
			list = list[:0]
			if err := ls.LookupFields(tx, "Track", m.ID, extralang, &list); err != nil {
				return errors.Wrap(err, `failed to lookup localized fields`)
			}

			for _, l := range list {
				switch l.Name {
				case "name":
					if len(m.Name) == 0 {
						if pdebug.Enabled {
							pdebug.Printf("Replacing for key 'name' (fallback en -> %s", l.Language)
						}
						m.Name = l.Localized
					}
				}
			}
		}
		return nil
	case "all":
		for _, extralang := range []string{`ja`} {
			list = list[:0]
			if err := ls.LookupFields(tx, "Track", m.ID, extralang, &list); err != nil {
				return errors.Wrap(err, `failed to lookup localized fields`)
			}

			for _, l := range list {
				if pdebug.Enabled {
					pdebug.Printf("Adding key '%s#%s'", l.Name, l.Language)
				}
				m.LocalizedFields.Set(l.Language, l.Name, l.Localized)
			}
		}
	default:
		for _, extralang := range []string{`ja`} {
			list = list[:0]
			if err := ls.LookupFields(tx, "Track", m.ID, extralang, &list); err != nil {
				return errors.Wrap(err, `failed to lookup localized fields`)
			}

			for _, l := range list {
				switch l.Name {
				case "name":
					if pdebug.Enabled {
						pdebug.Printf("Replacing for key 'name'")
					}
					m.Name = l.Localized
				}
			}
		}
	}
	return nil
}

func (v *TrackSvc) Delete(tx *db.Tx, id string) error {
	if pdebug.Enabled {
		g := pdebug.Marker("Track.Delete (%s)", id)
		defer g.End()
	}
	original := db.Track{EID: id}
	if err := original.LoadByEID(tx, id); err != nil {
		return errors.Wrap(err, `failed load before delete`)
	}

	vdb := db.Track{EID: id}
	if err := vdb.Delete(tx); err != nil {
		return errors.Wrap(err, `failed to delete from database`)
	}
	c := Cache()
	key := c.Key("Track", id)
	c.Delete(key)
	if pdebug.Enabled {
		pdebug.Printf(`CACHE DEL %s`, key)
	}
	if err := db.DeleteLocalizedStringsForParent(tx, id, "Track"); err != nil {
		return errors.Wrap(err, `failed to delete localized strings`)
	}
	if err := v.PostDeleteHook(tx, &original); err != nil {
		return errors.Wrap(err, `post delete hook failed`)
	}
	return nil
}
