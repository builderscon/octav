package service

import (
	"time"

	"context"

	"github.com/builderscon/octav/octav/cache"
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/internal/errors"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	"github.com/lestrrat/go-pdebug"
)

func (v *SponsorSvc) Init() {
	v.mediaStorage = MediaStorage
}

func (v *SponsorSvc) populateRowForCreate(vdb *db.Sponsor, payload *model.CreateSponsorRequest) error {
	vdb.EID = tools.UUID()

	vdb.ConferenceID = payload.ConferenceID
	vdb.Name = payload.Name
	vdb.URL = payload.URL
	vdb.GroupName = payload.GroupName
	vdb.SortOrder = payload.SortOrder
	if payload.LogoURL.Valid() {
		vdb.LogoURL.Valid = true
		vdb.LogoURL.String = payload.LogoURL.String
	}

	return nil
}

func (v *SponsorSvc) populateRowForUpdate(vdb *db.Sponsor, payload *model.UpdateSponsorRequest) error {
	if payload.Name.Valid() {
		vdb.Name = payload.Name.String
	}

	if payload.LogoURL.Valid() {
		vdb.LogoURL.Valid = true
		vdb.LogoURL.String = payload.LogoURL.String
	}

	if payload.URL.Valid() {
		vdb.URL = payload.URL.String
	}

	if payload.GroupName.Valid() {
		vdb.GroupName = payload.GroupName.String
	}

	if payload.SortOrder.Valid() {
		vdb.SortOrder = int(payload.SortOrder.Int)
	}

	return nil
}

type finalizeFunc func() error

func (ff finalizeFunc) FinalizeFunc() func() error {
	return ff
}

// Ignorable always returns true, otherwise the caller will have to
// bail out immediately
func (ff finalizeFunc) Ignorable() bool {
	return true
}

func (ff finalizeFunc) Error() string {
	return "operation needs finalization"
}

func (v *SponsorSvc) CreateFromPayload(ctx context.Context, tx *db.Tx, payload *model.AddSponsorRequest, result *model.Sponsor) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Sponsor.CreateFromPayload").BindError(&err)
		defer g.End()
	}

	su := User()
	if err := su.IsConferenceAdministrator(tx, payload.ConferenceID, payload.UserID); err != nil {
		return errors.Wrap(err, "creating a featured sponsor requires conference administrator privilege")
	}

	var vdb db.Sponsor
	if err := v.Create(tx, &vdb, &model.CreateSponsorRequest{payload}); err != nil {
		return errors.Wrap(err, "failed to store in database")
	}

	var m model.Sponsor
	if err := m.FromRow(&vdb); err != nil {
		return errors.Wrap(err, "failed to populate model from database")
	}

	*result = m

	c := Cache()
	key := c.Key("Sponsor", "LoadByConferenceID", payload.ConferenceID)
	c.Delete(key)
	if pdebug.Enabled {
		pdebug.Printf("CACHE DEL %s", key)
	}

	return nil
}

/*
func (v *SponsorSvc) UpdateFromPayload(ctx context.Context, tx *db.Tx, payload *model.UpdateSponsorRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Sponsor.UpdateFromPayload").BindError(&err)
		defer g.End()
	}

	var vdb db.Sponsor
	if err := vdb.LoadByEID(tx, payload.ID); err != nil {
		return errors.Wrap(err, "failed to load featured sponsor from database")
	}

	su := User()
	if err := su.IsConferenceAdministrator(tx, vdb.ConferenceID, payload.UserID); err != nil {
		return errors.Wrap(err, "updating a featured sponsor requires conference administrator privilege")
	}

	if err := v.Update(tx, &vdb); err != nil {
		return errors.Wrap(err, "failed to update sponsor in database")
	}

	return nil

}*/

func (v *SponsorSvc) DeleteFromPayload(ctx context.Context, tx *db.Tx, payload *model.DeleteSponsorRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Sponsor.DeleteFromPayload").BindError(&err)
		defer g.End()
	}

	var m db.Sponsor
	if err := m.LoadByEID(tx, payload.ID); err != nil {
		return errors.Wrap(err, "failed to load featured sponsor from database")
	}

	su := User()
	if err := su.IsConferenceAdministrator(tx, m.ConferenceID, payload.UserID); err != nil {
		return errors.Wrap(err, "deleting venues require administrator privileges")
	}

	if err := v.Delete(tx, m.EID); err != nil {
		return errors.Wrap(err, "failed to delete from database")
	}

	c := Cache()
	key := c.Key("Sponsor", "LoadByConferenceID", m.ConferenceID)
	c.Delete(key)
	if pdebug.Enabled {
		pdebug.Printf("CACHE DEL %s", key)
	}

	// For (current) testing purposes, we don't want to actually
	// access the Google storage backend.
	if InTesting {
		return
	}

	// This operation need not necessarily succeed. Spawn goroutines and deal with
	// it asynchronously
	go func() {
		if pdebug.Enabled {
			g := pdebug.Marker("service.Sponsor.DeleteFromPayload cleanup")
			defer g.End()
		}
		prefix := "conferences/" + m.ConferenceID + "/" + m.EID + "-logo"
		if pdebug.Enabled {
			pdebug.Printf("Listing objects that match query '%s'", prefix)
		}

		cl := v.mediaStorage
		list, err := cl.List(ctx, WithQueryPrefix(prefix))
		if err != nil {
			if pdebug.Enabled {
				pdebug.Printf("Error while listing '%s'", prefix)
			}
			return
		}
		cl.DeleteObjects(ctx, list)
	}()

	return nil
}

func (v *SponsorSvc) ListFromPayload(tx *db.Tx, result *model.SponsorList, payload *model.ListSponsorsRequest) error {
	var vdbl db.SponsorList
	if err := vdbl.LoadByConferenceSinceEID(tx, payload.ConferenceID, payload.Since.String, int(payload.Limit.Int)); err != nil {
		return errors.Wrap(err, "failed to load featured sponsor from database")
	}

	l := make(model.SponsorList, len(vdbl))
	for i, vdb := range vdbl {
		if err := (l[i]).FromRow(&vdb); err != nil {
			return errors.Wrap(err, "failed to populate model from database")
		}

		if err := v.Decorate(tx, &l[i], payload.TrustedCall, payload.Lang.String); err != nil {
			return errors.Wrap(err, "failed to decorate venue with associated data")
		}
	}

	*result = l
	return nil
}

func (v *SponsorSvc) Decorate(tx *db.Tx, sponsor *model.Sponsor, trustedCall bool, lang string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Sponsor.Decorate").BindError(&err)
		defer g.End()
	}

	if sponsor.LogoURL == "" {
		sponsor.LogoURL = "https://storage.googleapis.com/media-builderscon-1248/system/nophoto_600.png"
	}

	if lang == "" {
		return nil
	}

	if err := v.ReplaceL10NStrings(tx, sponsor, lang); err != nil {
		return errors.Wrap(err, "failed to replace L10N strings")
	}

	return nil
}

func (v *SponsorSvc) LoadByConferenceID(tx *db.Tx, cdl *model.SponsorList, cid string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("serviec.Sponsor.LoadByConferenceID %s", cid).BindError(&err)
		defer g.End()
	}

	var ids []string
	c := Cache()
	key := c.Key("Sponsor", "LoadByConferenceID", cid)
	if err := c.Get(key, &ids); err == nil {
		if pdebug.Enabled {
			pdebug.Printf("CACHE HIT: %s", key)
		}
		m := make(model.SponsorList, len(ids))
		for i, id := range ids {
			if err := v.Lookup(tx, &m[i], id); err != nil {
				return errors.Wrap(err, "failed to load from database")
			}
		}

		*cdl = m
		return nil
	}

	if pdebug.Enabled {
		pdebug.Printf("CACHE MISS: %s", key)
	}
	var vdbl db.SponsorList
	if err := db.LoadSponsors(tx, &vdbl, cid); err != nil {
		return err
	}

	ids = make([]string, len(vdbl))
	res := make(model.SponsorList, len(vdbl))
	for i, vdb := range vdbl {
		var u model.Sponsor
		if err := u.FromRow(&vdb); err != nil {
			return err
		}
		ids[i] = vdb.EID
		res[i] = u
	}
	*cdl = res

	c.Set(key, ids, cache.WithExpires(15*time.Minute))
	return nil
}
