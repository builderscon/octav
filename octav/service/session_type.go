package service

import (
	"time"

	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
)

func (v *SessionType) populateRowForCreate(vdb *db.SessionType, payload model.CreateSessionTypeRequest) error {
	vdb.EID = tools.RandomString(64)
	vdb.Name = payload.Name
	vdb.ConferenceID = payload.ConferenceID
	vdb.Abstract = payload.Abstract
	vdb.Duration = payload.Duration

	if payload.SubmissionStart.Valid() {
		vdb.SubmissionStart.Valid = true
		vdb.SubmissionStart.Time = time.Time(payload.SubmissionStart.JSONTime)
	}

	if payload.SubmissionEnd.Valid() {
		vdb.SubmissionEnd.Valid = true
		vdb.SubmissionEnd.Time = time.Time(payload.SubmissionEnd.JSONTime)
	}

	return nil
}

func (v *SessionType) populateRowForUpdate(vdb *db.SessionType, payload model.UpdateSessionTypeRequest) error {
	if payload.Name.Valid() {
		vdb.Name = payload.Name.String
	}

	if payload.Abstract.Valid() {
		vdb.Abstract = payload.Abstract.String
	}

	if payload.Duration.Valid() {
		vdb.Duration = int(payload.Duration.Int)
	}

	if payload.SubmissionStart.Valid() {
		vdb.SubmissionStart.Valid = true
		vdb.SubmissionStart.Time = time.Time(payload.SubmissionStart.JSONTime)
	}

	if payload.SubmissionEnd.Valid() {
		vdb.SubmissionEnd.Valid = true
		vdb.SubmissionEnd.Time = time.Time(payload.SubmissionEnd.JSONTime)
	}

	return nil
}
