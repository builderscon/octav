package service

import (
	"context"
	"crypto/sha1"
	"io"

	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	pdebug "github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

func (v *BlogEntrySvc) Init() {}

func (v *BlogEntrySvc) populateRowForCreate(vdb *db.BlogEntry, payload *model.CreateBlogEntryRequest) error {
	vdb.ConferenceID = payload.ConferenceID
	vdb.Title = payload.Title
	vdb.URL = payload.URL

	h := sha1.New()
	io.WriteString(h, payload.URL)
	vdb.URLHash = string(h.Sum(nil))
	return nil
}

func (v *BlogEntrySvc) populateRowForUpdate(vdb *db.BlogEntry, payload *model.UpdateBlogEntryRequest) error {
	if payload.Title.Valid() {
		vdb.Title = payload.Title.String
	}
	if payload.URL.Valid() {
		vdb.URL = payload.URL.String
	}
	return nil
}

func (v *BlogEntrySvc) CreateFromPayload(ctx context.Context, tx *db.Tx, result *model.BlogEntry, payload *model.CreateBlogEntryRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("BlogEntrySvc.CreateFromPayload").BindError(&err)
		defer g.End()
	}

	su := User()
	if err := su.IsConferenceAdministrator(tx, payload.ConferenceID, payload.UserID); err != nil {
		return errors.Wrap(err, "creating blog entries require conference administrator privileges")
	}

	var vdb db.BlogEntry
	if err := v.Create(tx, &vdb, payload); err != nil {
		return errors.Wrap(err, "failed to insert into database")
	}

	if result != nil {
		var m model.BlogEntry
		if err := m.FromRow(&vdb); err != nil {
			return errors.Wrap(err, "failed to populate model from database")
		}

		*result = m
	}
	return nil
}

func (v *BlogEntrySvc) DeleteFromPayload(ctx context.Context, tx *db.Tx, payload *model.DeleteBlogEntryRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("BlogEntrySvc.DeleteFromPayload").BindError(&err)
		defer g.End()
	}

	var m model.BlogEntry
	if err := v.Lookup(tx, &m, payload.ID); err != nil {
		return errors.Wrap(err, "failed to look up blog entry")
	}

	su := User()
	if err := su.IsConferenceAdministrator(tx, m.ConferenceID, payload.UserID); err != nil {
		return errors.Wrap(err, "deleting blog entries requires conference administrator privilege")
	}

	if err := v.Delete(tx, payload.ID); err != nil {
		return errors.Wrap(err, "failed to delete session")
	}
	return nil
}
