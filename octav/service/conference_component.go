package service

import (
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	"github.com/pkg/errors"
)

func (v *ConferenceComponentSvc) Init() {}

func (v *ConferenceComponentSvc) populateRowForCreate(vdb *db.ConferenceComponent, payload model.CreateConferenceComponentRequest) error {
	vdb.EID = tools.UUID()
	vdb.Name = payload.Name
	vdb.Value = payload.Value
	return nil
}

func (v *ConferenceComponentSvc) populateRowForUpdate(vdb *db.ConferenceComponent, payload model.UpdateConferenceComponentRequest) error {
	vdb.EID = tools.UUID()
	if payload.Name.Valid() {
		vdb.Name = payload.Name.String
	}

	if payload.Name.Valid() {
		vdb.Value = payload.Value.String
	}
	return nil
}

func (v *ConferenceComponentSvc) DeleteByConferenceIDAndName(tx *db.Tx, conferenceID string, names ...string) error {
	if err := db.DeleteConferenceComponentsByIDAndName(tx, conferenceID, names...); err != nil {
		return errors.Wrap(err, "failed to delete from database")
	}

	return nil
}

func (v *ConferenceComponentSvc) UpsertByConferenceIDAndName(tx *db.Tx, conferenceID string, values map[string]string) error {
	if err := db.UpsertConferenceComponentsByIDAndName(tx, conferenceID, values); err != nil {
		return errors.Wrap(err, "failed to update/insert into database")
	}

	return nil
}

