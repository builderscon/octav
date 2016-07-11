package service

import (
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	"github.com/pkg/errors"
)

func (v *ConferenceSeries) populateRowForCreate(vdb *db.ConferenceSeries, payload model.CreateConferenceSeriesRequest) error {
	vdb.EID = tools.UUID()
	vdb.Slug = payload.Slug
	return nil
}

func (v *ConferenceSeries) populateRowForUpdate(vdb *db.ConferenceSeries, payload model.UpdateConferenceSeriesRequest) error {
	if payload.Slug.Valid() {
		vdb.Slug = payload.Slug.String
	}

	return nil
}

func (v *ConferenceSeries) CreateFromPayload(tx *db.Tx, payload model.CreateConferenceSeriesRequest, result *model.ConferenceSeries) error {
	vdb := db.ConferenceSeries{}
	if err := v.Create(tx, &vdb, payload); err != nil {
		return errors.Wrap(err, "failed to create store conference in database")
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

