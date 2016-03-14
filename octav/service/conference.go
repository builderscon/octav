package service

import (
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
)

func (v *Conference) populateRowForCreate(vdb *db.Conference, payload model.CreateConferenceRequest) error {
	vdb.EID = tools.UUID()
	vdb.Slug = payload.Slug
	vdb.Title = payload.Title

	if payload.SubTitle.Valid() {
		vdb.SubTitle.Valid = true
		vdb.SubTitle.String = payload.SubTitle.String
	}
	return nil
}

func (v *Conference) populateRowForUpdate(vdb *db.Conference, payload model.UpdateConferenceRequest) error {
	if payload.Slug.Valid() {
		vdb.Slug = payload.Slug.String
	}

	if payload.Slug.Valid() {
		vdb.Title = payload.Title.String
	}

	if payload.SubTitle.Valid() {
		vdb.SubTitle.Valid = true
		vdb.SubTitle.String = payload.SubTitle.String
	}
	return nil
}

