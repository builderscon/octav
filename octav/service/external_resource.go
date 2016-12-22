package service

import (
	"time"
	"net/url"

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

	if payload.Description.Valid() {
		vdb.Description = payload.Description.String
	}

	// Parse the URL, and do away with the URL fragment, if any
	u, err := url.Parse(payload.URL)
	if err != nil {
		return errors.Wrap(err, "failed to parse URL")
	}
	u.Fragment = ""
	vdb.URL = u.String()

	return nil
}

func (v *ExternalResourceSvc) populateRowForUpdate(vdb *db.ExternalResource, payload *model.UpdateExternalResourceRequest) error {
	if payload.Description.Valid() {
		vdb.Description = payload.Description.String
	}
	if payload.Name.Valid() {
		vdb.Name = payload.Name.String
	}
	if payload.URL.Valid() {
		// Parse the URL, and do away with the URL fragment, if any
		u, err := url.Parse(payload.URL.String)
		if err != nil {
			return errors.Wrap(err, "failed to parse URL")
		}
		u.Fragment = ""
		vdb.URL = u.String()
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

func (v *ExternalResourceSvc) LoadByConferenceID(tx *db.Tx, result *model.ExternalResourceList, cid string, trustedCall bool, lang string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.ExternalResource.LoadByConferenceID %s", cid).BindError(&err)
		defer g.End()
	}

	c := Cache()
	key := c.Key("ExternalResource", "ListFromPayload", cid)
	x, err := c.GetOrSet(key, result, func() (interface{}, error) {
		if pdebug.Enabled {
			pdebug.Printf("CACHE MISS: Re-generating")
		}

		var vdbl db.ExternalResourceList
		if err := vdbl.LoadByConference(tx, cid); err != nil {
			return nil, errors.Wrap(err, "failed to load from database")
		}

		l := make(model.ExternalResourceList, len(vdbl))
		for i, vdb := range vdbl {
			if err := l[i].FromRow(&vdb); err != nil {
				return nil, errors.Wrap(err, "failed to populate model from database")
			}

			if err := v.Decorate(tx, &l[i], trustedCall, lang); err != nil {
				return nil, errors.Wrap(err, "failed to decorate ExternalResource with associated data")
			}
		}

		return &l, nil
	}, cache.WithExpires(time.Hour))

	if err != nil {
		return err
	}

	*result = *(x.(*model.ExternalResourceList))
	return nil
}
