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

var conferenceComponentSvc *ConferenceComponentSvc
var conferenceComponentOnce sync.Once

func ConferenceComponent() *ConferenceComponentSvc {
	conferenceComponentOnce.Do(conferenceComponentSvc.Init)
	return conferenceComponentSvc
}

func (v *ConferenceComponentSvc) LookupFromPayload(tx *db.Tx, m *model.ConferenceComponent, payload model.LookupConferenceComponentRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.ConferenceComponent.LookupFromPayload").BindError(&err)
		defer g.End()
	}
	if err = v.Lookup(tx, m, payload.ID); err != nil {
		return errors.Wrap(err, "failed to load model.ConferenceComponent from database")
	}
	return nil
}
func (v *ConferenceComponentSvc) Lookup(tx *db.Tx, m *model.ConferenceComponent, id string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.ConferenceComponent.Lookup").BindError(&err)
		defer g.End()
	}

	r := model.ConferenceComponent{}
	if err = r.Load(tx, id); err != nil {
		return errors.Wrap(err, "failed to load model.ConferenceComponent from database")
	}
	*m = r
	return nil
}

// Create takes in the transaction, the incoming payload, and a reference to
// a database row. The database row is initialized/populated so that the
// caller can use it afterwards.
func (v *ConferenceComponentSvc) Create(tx *db.Tx, vdb *db.ConferenceComponent, payload model.CreateConferenceComponentRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.ConferenceComponent.Create").BindError(&err)
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

func (v *ConferenceComponentSvc) Update(tx *db.Tx, vdb *db.ConferenceComponent, payload model.UpdateConferenceComponentRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.ConferenceComponent.Update (%s)", vdb.EID).BindError(&err)
		defer g.End()
	}

	if vdb.EID == "" {
		return errors.New("vdb.EID is required (did you forget to call vdb.Load(tx) before hand?)")
	}

	if err := v.populateRowForUpdate(vdb, payload); err != nil {
		return err
	}

	if err := vdb.Update(tx); err != nil {
		return err
	}
	return nil
}

func (v *ConferenceComponentSvc) Delete(tx *db.Tx, id string) error {
	if pdebug.Enabled {
		g := pdebug.Marker("ConferenceComponent.Delete (%s)", id)
		defer g.End()
	}

	vdb := db.ConferenceComponent{EID: id}
	if err := vdb.Delete(tx); err != nil {
		return err
	}
	return nil
}

func (v *ConferenceComponentSvc) LoadList(tx *db.Tx, vdbl *db.ConferenceComponentList, since string, limit int) error {
	return vdbl.LoadSinceEID(tx, since, limit)
}
