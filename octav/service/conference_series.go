package service

import (
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	"github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

func (v *ConferenceSeries) populateRowForCreate(vdb *db.ConferenceSeries, payload model.CreateConferenceSeriesRequest) error {
	vdb.EID = tools.UUID()
	vdb.Slug = payload.Slug
	vdb.Title = payload.Title
	return nil
}

func (v *ConferenceSeries) populateRowForUpdate(vdb *db.ConferenceSeries, payload model.UpdateConferenceSeriesRequest) error {
	if payload.Slug.Valid() {
		vdb.Slug = payload.Slug.String
	}

	if payload.Title.Valid() {
		vdb.Title = payload.Title.String
	}

	return nil
}

func (v *ConferenceSeries) DeleteFromPayload(tx *db.Tx, payload model.DeleteConferenceSeriesRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.ConferenceSeries.DeleteFromPayload").BindError(&err)
		defer g.End()
	}

	u := model.User{}
	su := User{}
	if err := su.Lookup(tx, &u, payload.UserID); err != nil {
		return errors.Wrap(err, "failed to load user from database")
	}

	// The user must be a system admin
	if !u.IsAdmin {
		return errors.New("user lacks system administrator privileges")
	}

	return errors.Wrap(v.Delete(tx, payload.ID), "failed to delete from database")
}

// CreateFromPayload adds extra logic around Create to verify data
// and create accessory data.
func (v *ConferenceSeries) CreateFromPayload(tx *db.Tx, result *model.ConferenceSeries, payload model.CreateConferenceSeriesRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.ConferenceSeries.CreateFromPayload").BindError(&err)
		defer g.End()
	}

	su := User{}
	if err := su.IsSystemAdmin(tx, payload.UserID); err != nil {
		return errors.Wrap(err, "creating a conference series requires system administrator privilege")
	}

	vdb := db.ConferenceSeries{}
	if err := v.Create(tx, &vdb, payload); err != nil {
		return errors.Wrap(err, "failed to store conference series in database")
	}

	csa := db.ConferenceSeriesAdministrator{
		SeriesID: vdb.EID,
		UserID:   payload.UserID,
	}
	if err := csa.Create(tx); err != nil {
		return errors.Wrap(err, "failed to store conference series administrator in database")
	}

	c := model.ConferenceSeries{}
	if err := c.FromRow(vdb); err != nil {
		return errors.Wrap(err, "failed to populate model from database")
	}

	*result = c
	return nil
}

func (v *ConferenceSeries) LoadByRange(tx *db.Tx, l *[]model.ConferenceSeries, since string, limit int) error {
	vdbl := db.ConferenceSeriesList{}
	if err := vdbl.LoadSinceEID(tx, since, limit); err != nil {
		return errors.Wrap(err, "failed to load list of conference series")
	}

	csl := make([]model.ConferenceSeries, len(vdbl))
	for i, row := range vdbl {
		if err := (csl[i]).FromRow(row); err != nil {
			return errors.Wrap(err, "failed to convert row to model")
		}
	}

	*l = csl
	return nil
}

func (v *ConferenceSeries) AddAdministratorFromPayload(tx *db.Tx, payload model.AddConferenceSeriesAdminRequest) error {
	if err := db.IsConferenceSeriesAdministrator(tx, payload.SeriesID, payload.UserID); err != nil {
		return errors.Wrap(err, "adding a conference series administrator requires conference series administrator privilege")
	}
	return errors.Wrap(v.AddAdministrator(tx, payload.SeriesID, payload.AdminID), "failed to add administrator")
}

func (v *ConferenceSeries) AddAdministrator(tx *db.Tx, seriesID, userID string) error {
	c := db.ConferenceSeriesAdministrator{
		SeriesID: seriesID,
		UserID:   userID,
	}
	return c.Create(tx, db.WithInsertIgnore(true))
}

func (v *ConferenceSeries) Decorate(tx *db.Tx, c *model.ConferenceSeries, trustedCall bool, lang string) error {
	if lang == "" {
		return nil
	}
	if err := v.ReplaceL10NStrings(tx, c, lang); err != nil {
		return errors.Wrap(err, "failed to replace L10N strings")
	}
	return nil
}
