package service

// Automatically generated by genmodel utility. DO NOT EDIT!

import (
	"sync"
	"time"

	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/internal/errors"
	"github.com/builderscon/octav/octav/model"
	"github.com/lestrrat/go-pdebug"
)

var _ = time.Time{}

var conferenceDateSvc ConferenceDateSvc
var conferenceDateOnce sync.Once

func ConferenceDate() *ConferenceDateSvc {
	conferenceDateOnce.Do(conferenceDateSvc.Init)
	return &conferenceDateSvc
}

func (v *ConferenceDateSvc) Lookup(tx *db.Tx, m *model.ConferenceDate, id string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.ConferenceDate.Lookup").BindError(&err)
		defer g.End()
	}

	r := model.ConferenceDate{}
	if err = r.Load(tx, id); err != nil {
		return errors.Wrap(err, "failed to load model.ConferenceDate from database")
	}
	*m = r
	return nil
}

// Create takes in the transaction, the incoming payload, and a reference to
// a database row. The database row is initialized/populated so that the
// caller can use it afterwards.
func (v *ConferenceDateSvc) Create(tx *db.Tx, vdb *db.ConferenceDate, payload model.CreateConferenceDateRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.ConferenceDate.Create").BindError(&err)
		defer g.End()
	}

	if err := v.populateRowForCreate(vdb, payload); err != nil {
		return err
	}

	if err := vdb.Create(tx); err != nil {
		return err
	}

	return nil
}

func (v *ConferenceDateSvc) Delete(tx *db.Tx, id string) error {
	if pdebug.Enabled {
		g := pdebug.Marker("ConferenceDate.Delete (%s)", id)
		defer g.End()
	}

	vdb := db.ConferenceDate{EID: id}
	if err := vdb.Delete(tx); err != nil {
		return err
	}
	return nil
}

func (v *ConferenceDateSvc) LoadList(tx *db.Tx, vdbl *db.ConferenceDateList, since string, limit int) error {
	return vdbl.LoadSinceEID(tx, since, limit)
}