package service

import (
	"fmt"
	"time"

	"github.com/builderscon/octav/octav/cache"
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	pdebug "github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

func (v *VenueSvc) Init() {}

func (v *VenueSvc) populateRowForCreate(vdb *db.Venue, payload *model.CreateVenueRequest) error {
	vdb.EID = tools.UUID()
	vdb.Name = payload.Name.String
	vdb.Address = payload.Address.String
	vdb.Latitude = payload.Latitude.Float
	vdb.Longitude = payload.Longitude.Float
	return nil
}

func (v *VenueSvc) populateRowForUpdate(vdb *db.Venue, payload *model.UpdateVenueRequest) error {
	if payload.Name.Valid() {
		vdb.Name = payload.Name.String
	}

	if payload.Address.Valid() {
		vdb.Address = payload.Address.String
	}

	if payload.Latitude.Valid() {
		vdb.Latitude = payload.Latitude.Float
	}

	if payload.Longitude.Valid() {
		vdb.Longitude = payload.Longitude.Float
	}
	return nil
}

func (v *VenueSvc) LoadRooms(tx *db.Tx, cdl *model.RoomList, venueID string) error {
	var vdbl db.RoomList
	if err := db.LoadVenueRooms(tx, &vdbl, venueID); err != nil {
		return err
	}

	res := make(model.RoomList, len(vdbl))
	for i, vdb := range vdbl {
		var u model.Room
		if err := u.FromRow(vdb); err != nil {
			return err
		}
		res[i] = u
	}
	*cdl = res
	return nil
}

func (v *VenueSvc) Decorate(tx *db.Tx, venue *model.Venue, trustedCall bool, lang string) error {
	if err := v.LoadRooms(tx, &venue.Rooms, venue.ID); err != nil {
		return err
	}

	if lang != "" {
		sr := Room()
		for i := range venue.Rooms {
			if err := sr.ReplaceL10NStrings(tx, &venue.Rooms[i], lang); err != nil {
				return errors.Wrap(err, "failed to replace L10N strings")
			}
		}

		if err := v.ReplaceL10NStrings(tx, venue, lang); err != nil {
			return errors.Wrap(err, "failed to replace L10N strings")
		}
	}

	return nil
}

func (v *VenueSvc) CreateFromPayload(tx *db.Tx, venue *model.Venue, payload *model.CreateVenueRequest) error {
	su := User()
	if err := su.IsAdministrator(tx, payload.UserID); err != nil {
		return errors.Wrap(err, "creating venues require administrator privileges")
	}

	var vdb db.Venue
	if err := v.Create(tx, &vdb, payload); err != nil {
		return errors.Wrap(err, "failed to store in database")
	}

	var r model.Venue
	if err := r.FromRow(vdb); err != nil {
		return errors.Wrap(err, "failed to populate model from database")
	}
	*venue = r

	return nil
}

func (v *VenueSvc) DeleteFromPayload(tx *db.Tx, payload *model.DeleteVenueRequest) error {
	su := User()
	if err := su.IsAdministrator(tx, payload.UserID); err != nil {
		return errors.Wrap(err, "deleting venues require administrator privileges")
	}

	return errors.Wrap(v.Delete(tx, payload.ID), "failed to delete from database")
}

func (v *VenueSvc) ListFromPayload(tx *db.Tx, result *model.VenueList, payload *model.ListVenueRequest) error {
	var vdbl db.VenueList
	if err := vdbl.LoadSinceEID(tx, payload.Since.String, int(payload.Limit.Int)); err != nil {
		return errors.Wrap(err, "failed to load from database")
	}

	l := make(model.VenueList, len(vdbl))
	for i, vdb := range vdbl {
		if err := (l[i]).FromRow(vdb); err != nil {
			return errors.Wrap(err, "failed to populate model from database")
		}

		if err := v.Decorate(tx, &l[i], payload.TrustedCall, payload.Lang.String); err != nil {
			return errors.Wrap(err, "failed to decorate venue with associated data")
		}
	}

	*result = l
	return nil
}

func (v *VenueSvc) LoadByConferenceID(tx *db.Tx, cdl *model.VenueList, cid, lang string, trustedCall bool) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Venue.LoadByConferenceID (%s,%s,%t)", cid, lang, trustedCall).BindError(&err)
		defer g.End()
		defer func() {
			pdebug.Printf("Loaded %d venues", len(*cdl))
		}()
	}

	c := Cache()
	key := c.Key("Venue", "LoadByConferenceID", cid, lang, fmt.Sprintf("%t", trustedCall))

	var ids []string
	if err := c.Get(key, &ids); err == nil {
		if pdebug.Enabled {
			pdebug.Printf("CACHE HIT %s", key)
		}
		m := make(model.VenueList, len(ids))
		r := model.LookupVenueRequest{
			TrustedCall: trustedCall,
		}
		r.Lang.Set(lang)
		for i, id := range ids {
			r.ID = id
			if err := v.LookupFromPayload(tx, &m[i], &r); err != nil {
				return errors.Wrap(err, "failed to lookup venue")
			}
		}

		*cdl = m
		return nil
	}

	if pdebug.Enabled {
		pdebug.Printf("CACHE MISS %s", key)
	}
	var vdbl db.VenueList
	if err := db.LoadConferenceVenues(tx, &vdbl, cid); err != nil {
		return errors.Wrap(err, "failed to load venues from database")
	}

	res := make(model.VenueList, len(vdbl))
	ids = make([]string, len(vdbl))
	for i, vdb := range vdbl {
		var u model.Venue
		if err := u.FromRow(vdb); err != nil {
			return err
		}
		ids[i] = vdb.EID
		res[i] = u
	}
	*cdl = res

	c.Set(key, ids, cache.WithExpires(15*time.Minute))
	return nil
}
