package model

// Automatically generated by genmodel utility. DO NOT EDIT!

import (
	"time"

	"github.com/builderscon/octav/octav/db"
	"github.com/lestrrat/go-pdebug"
)

var _ = time.Time{}

func (v *SessionType) Load(tx *db.Tx, id string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("model.SessionType.Load %s", id).BindError(&err)
		defer g.End()
	}
	vdb := db.SessionType{}
	if err := vdb.LoadByEID(tx, id); err != nil {
		return err
	}

	if err := v.FromRow(vdb); err != nil {
		return err
	}
	return nil
}

func (v *SessionType) FromRow(vdb db.SessionType) error {
	v.ID = vdb.EID
	v.ConferenceID = vdb.ConferenceID
	v.Name = vdb.Name
	v.Abstract = vdb.Abstract
	v.Duration = vdb.Duration
	if vdb.SubmissionStart.Valid {
		v.SubmissionStart = vdb.SubmissionStart.Time
	}
	if vdb.SubmissionEnd.Valid {
		v.SubmissionEnd = vdb.SubmissionEnd.Time
	}
	return nil
}

func (v *SessionType) ToRow(vdb *db.SessionType) error {
	vdb.EID = v.ID
	vdb.ConferenceID = v.ConferenceID
	vdb.Name = v.Name
	vdb.Abstract = v.Abstract
	vdb.Duration = v.Duration
	vdb.SubmissionStart.Valid = true
	vdb.SubmissionStart.Time = v.SubmissionStart
	vdb.SubmissionEnd.Valid = true
	vdb.SubmissionEnd.Time = v.SubmissionEnd
	return nil
}
