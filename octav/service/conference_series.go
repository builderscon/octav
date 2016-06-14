package service

import (
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
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
