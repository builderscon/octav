package service

import (
	"database/sql"
	"net/url"
	"time"

	"github.com/builderscon/octav/octav/cache"
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/internal/context"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	"github.com/pkg/errors"

	pdebug "github.com/lestrrat/go-pdebug"
)

func (v *ExternalResourceSvc) Init() {}

func (v *ExternalResourceSvc) populateRowForCreate(ctx context.Context, vdb *db.ExternalResource, payload *model.CreateExternalResourceRequest) error {
	vdb.EID = tools.UUID()
	vdb.ConferenceID = payload.ConferenceID
	vdb.Title = payload.Title

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

	if payload.SortOrder.Valid() {
		vdb.SortOrder = int(payload.SortOrder.Int)
	}

	return nil
}

func (v *ExternalResourceSvc) populateRowForUpdate(ctx context.Context, vdb *db.ExternalResource, payload *model.UpdateExternalResourceRequest) error {
	if payload.Description.Valid() {
		vdb.Description = payload.Description.String
	}
	if payload.Title.Valid() {
		vdb.Title = payload.Title.String
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
	if payload.SortOrder.Valid() {
		vdb.SortOrder = int(payload.SortOrder.Int)
	}
	return nil
}

func (v *ExternalResourceSvc) CreateFromPayload(ctx context.Context, tx *sql.Tx, result *model.ExternalResource, payload *model.CreateExternalResourceRequest) error {
	var vdb db.ExternalResource
	if err := v.Create(ctx, tx, &vdb, payload); err != nil {
		return errors.Wrap(err, "failed to insert into database")
	}

	var m model.ExternalResource
	if err := m.FromRow(&vdb); err != nil {
		return errors.Wrap(err, "failed to populate model from database row")
	}
	*result = m

	return nil
}

func (v *ExternalResourceSvc) DeleteFromPayload(ctx context.Context, tx *sql.Tx, payload *model.DeleteExternalResourceRequest) error {
	su := User()
	if err := su.IsAdministrator(ctx, tx, context.GetUserID(ctx)); err != nil {
		return errors.Wrap(err, "deleting exrernal resources require administrator privileges")
	}

	return errors.Wrap(v.Delete(tx, payload.ID), "failed to delete from database")
}

func (v *ExternalResourceSvc) ListFromPayload(ctx context.Context, tx *sql.Tx, result *model.ExternalResourceList, payload *model.ListExternalResourceRequest) error {
	var vdbl db.ExternalResourceList
	if err := vdbl.LoadSinceEID(tx, payload.Since.String, int(payload.Limit.Int)); err != nil {
		return errors.Wrap(err, "failed to load from database")
	}

	l := make(model.ExternalResourceList, len(vdbl))
	for i, vdb := range vdbl {
		if err := (l[i]).FromRow(&vdb); err != nil {
			return errors.Wrap(err, "failed to populate model from database")
		}

		if err := v.Decorate(ctx, tx, &l[i], payload.TrustedCall, payload.Lang.String); err != nil {
			return errors.Wrap(err, "failed to decorate exrernal resource with associated data")
		}
	}

	*result = l
	return nil
}

func invalidateExternalResourceLoadByConferenceID(conferenceID string) error {
	c := Cache()
	key := c.Key("ExternalResource", "LoadByConferenceID", conferenceID)
	c.Delete(key)
	if pdebug.Enabled {
		pdebug.Printf("CACHE DEL: %s", key)
	}
	return nil
}

func (v *ExternalResourceSvc) PostCreateHook(ctx context.Context, _ *sql.Tx, vdb *db.ExternalResource) error {
	return invalidateExternalResourceLoadByConferenceID(vdb.ConferenceID)
}

func (v *ExternalResourceSvc) PostUpdateHook(_ *sql.Tx, vdb *db.ExternalResource) error {
	return invalidateExternalResourceLoadByConferenceID(vdb.ConferenceID)
}

func (v *ExternalResourceSvc) PostDeleteHook(_ *sql.Tx, vdb *db.ExternalResource) error {
	return invalidateExternalResourceLoadByConferenceID(vdb.ConferenceID)
}

func (v *ExternalResourceSvc) Decorate(ctx context.Context, tx *sql.Tx, c *model.ExternalResource, trustedCall bool, lang string) error {
	if lang == "" {
		return nil
	}
	if err := v.ReplaceL10NStrings(tx, c, lang); err != nil {
		return errors.Wrap(err, "failed to replace L10N strings")
	}
	return nil
}

func (v *ExternalResourceSvc) LoadByConferenceID(tx *sql.Tx, result *model.ExternalResourceList, cid string, trustedCall bool, lang string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.ExternalResource.LoadByConferenceID %s", cid).BindError(&err)
		defer g.End()
	}

	c := Cache()
	key := c.Key("ExternalResource", "LoadByConferenceID", cid)
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
		}

		return &l, nil
	}, cache.WithExpires(time.Hour))

	if err != nil {
		return err
	}

	*result = *(x.(*model.ExternalResourceList))
	return nil
}
