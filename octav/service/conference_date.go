package service

import (
	"context"
	"database/sql"

	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	"github.com/pkg/errors"
)

func (v *ConferenceDateSvc) Init() {}

func (v *ConferenceDateSvc) populateRowForCreate(vdb *db.ConferenceDate, payload *model.CreateConferenceDateRequest) error {
	vdb.EID = tools.UUID()
	vdb.ConferenceID = payload.ConferenceID
	vdb.Open.Time = payload.Date.Open
	vdb.Open.Valid = true
	vdb.Close.Time = payload.Date.Close
	vdb.Close.Valid = true
	return nil
}

func (v *ConferenceDateSvc) CreateFromPayload(ctx context.Context, tx *sql.Tx, payload *model.CreateConferenceDateRequest, result *model.ConferenceDate) error {
	var vdb db.ConferenceDate
	if err := v.Create(ctx, tx, &vdb, payload); err != nil {
		return errors.Wrap(err, "failed to insert into database")
	}

	var m model.ConferenceDate
	if err := m.FromRow(&vdb); err != nil {
		return errors.Wrap(err, "failed to populate model from database row")
	}
	*result = m

	return nil
}
