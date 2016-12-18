package service

import (
	"time"

	"github.com/builderscon/octav/octav/cache"
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	"github.com/pkg/errors"

	pdebug "github.com/lestrrat/go-pdebug"
)

func (v *ExternalResourceSvc) Init() {}

func (v *ExternalResourceSvc) populateRowForCreate(vdb *db.ExternalResource, payload *model.CreateExternalResourceRequest) error {
	vdb.EID = tools.UUID()
	vdb.ConferenceID = payload.ConferenceID
	vdb.Name = payload.Name
	vdb.URL = payload.URL
	return nil
}

func (v *ExternalResourceSvc) populateRowForUpdate(vdb *db.ExternalResource, payload *model.UpdateExternalResourceRequest) error {
	if payload.Name.Valid() {
		vdb.Name = payload.Name.String
	}
	if payload.URL.Valid() {
		vdb.URL = payload.URL.String
	}
	return nil
}

func (v *ExternalResourceSvc) CreateFromPayload(tx *db.Tx, result *model.ExternalResource, payload *model.CreateExternalResourceRequest) error {
	var vdb db.ExternalResource
	if err := v.Create(tx, &vdb, payload); err != nil {
		return errors.Wrap(err, "failed to insert into database")
	}

	var m model.ExternalResource
	if err := m.FromRow(&vdb); err != nil {
		return errors.Wrap(err, "failed to populate model from database row")
	}
	*result = m

	return nil
}

func (v *ExternalResourceSvc) DeleteFromPayload(tx *db.Tx, payload *model.DeleteExternalResourceRequest) error {
	su := User()
	if err := su.IsAdministrator(tx, payload.UserID); err != nil {
		return errors.Wrap(err, "deleting exrernal resources require administrator privileges")
	}

	return errors.Wrap(v.Delete(tx, payload.ID), "failed to delete from database")
}

func (v *ExternalResourceSvc) ListFromPayload(tx *db.Tx, result *model.ExternalResourceList, payload *model.ListExternalResourceRequest) error {
	var vdbl db.ExternalResourceList
	if err := vdbl.LoadSinceEID(tx, payload.Since.String, int(payload.Limit.Int)); err != nil {
		return errors.Wrap(err, "failed to load from database")
	}

	l := make(model.ExternalResourceList, len(vdbl))
	for i, vdb := range vdbl {
		if err := (l[i]).FromRow(&vdb); err != nil {
			return errors.Wrap(err, "failed to populate model from database")
		}

		if err := v.Decorate(tx, &l[i], payload.TrustedCall, payload.Lang.String); err != nil {
			return errors.Wrap(err, "failed to decorate exrernal resource with associated data")
		}
	}

	*result = l
	return nil
}

func (v *ExternalResourceSvc) Decorate(tx *db.Tx, c *model.ExternalResource, trustedCall bool, lang string) error {
	if lang == "" {
		return nil
	}
	if err := v.ReplaceL10NStrings(tx, c, lang); err != nil {
		return errors.Wrap(err, "failed to replace L10N strings")
	}
	return nil
}

func (v *ExternalResourceSvc) LoadByConferenceID(tx *db.Tx, cdl *model.ExternalResourceList, cid string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.ExternalResource.LoadByConferenceID %s", cid).BindError(&err)
		defer g.End()
	}

	var ids []string
	c := Cache()
	key := c.Key("ExternalResource", "LoadByConferenceID", cid)
	if err := c.Get(key, &ids); err == nil {
		if pdebug.Enabled {
			pdebug.Printf("CACHE HIT: %s", key)
		}
		m := make(model.ExternalResourceList, len(ids))
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
	var vdbl db.ExternalResourceList
	if err := db.LoadExternalResources(tx, &vdbl, cid); err != nil {
		return err
	}

	ids = make([]string, len(vdbl))
	res := make(model.ExternalResourceList, len(vdbl))
	for i, vdb := range vdbl {
		var u model.ExternalResource
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
