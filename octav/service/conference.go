package service

import (
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/tools"
)

func (v *Conference) Create(tx *db.Tx, payload CreateConferenceRequest, vdb *db.Conference) error {
	vdb.EID = tools.UUID()
	vdb.Slug = payload.Slug
	vdb.Title = payload.Title

	if payload.SubTitle.Valid() {
		vdb.SubTitle.Valid = true
		vdb.SubTitle.String = payload.SubTitle.String
	}

	if err := vdb.Create(tx); err != nil {
		return err
	}

	if err := payload.L10N.CreateLocalizedStrings(tx, "Conference", vdb.EID); err != nil {
		return err
	}
	return nil
}
