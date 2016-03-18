package service

import (
	"database/sql"
	"strings"
	"time"

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

func (v *Conference) AddAdministrator(tx *db.Tx, cid, uid string) error {
	c := db.ConferenceAdministrator{
		ConferenceID: cid,
		UserID:       uid,
	}
	return c.Create(tx)
}

const datefmt = `2006-01-02`

func (v *Conference) LoadByRange(tx *db.Tx, vdbl *db.ConferenceList, since, lang, rangeStart, rangeEnd string, limit int) error {
	var rs time.Time
	var re time.Time
	var err error

	if rangeStart != "" {
		rs, err = time.Parse(datefmt, rangeStart)
		if err != nil {
			return err
		}
	}

	if rangeEnd != "" {
		re, err = time.Parse(datefmt, rangeEnd)
		if err != nil {
			return err
		}
	}

	if err := vdbl.LoadByRange(tx, since, rs, re, limit); err != nil {
		return err
	}

	return nil
}

func (v *Conference) AddDates(tx *db.Tx, cid string, dates ...model.ConferenceDate) error {
	for _, date := range dates {
		cd := db.ConferenceDate{
			ConferenceID: cid,
			Date:         date.Date.String(),
			Open:         sql.NullString{String: date.Open.String(), Valid: true},
			Close:        sql.NullString{String: date.Close.String(), Valid: true},
		}
		if err := cd.Create(tx, db.WithInsertIgnore(true)); err != nil {
			return err
		}
	}

	return nil
}

func (v *Conference) DeleteDates(tx *db.Tx, cid string, dates ...string) error {
	vdb := db.ConferenceDate{}
	return vdb.DeleteDates(tx, cid, dates...)
}

func (v *Conference) LoadDates(tx *db.Tx, cdl *model.ConferenceDateList, cid string) error {
	vdbl := db.ConferenceDateList{}
	if err := vdbl.LoadByConferenceID(tx, cid); err != nil {
		return err
	}

	res := make(model.ConferenceDateList, len(vdbl))
	for i, vdb := range vdbl {
		dt := vdb.Date
		if i := strings.IndexByte(dt, 'T'); i > -1 { // Cheat. Loading from DB contains time....!!!!
			dt = dt[:i]
		}
		if err := res[i].Date.Parse(dt); err != nil {
			return err
		}

		if vdb.Open.Valid {
			t := vdb.Open.String
			if len(t) > 5 {
				t = t[:5]
			}
			if err := res[i].Open.Parse(t); err != nil {
				return err
			}
		}

		if vdb.Close.Valid {
			t := vdb.Close.String
			if len(t) > 5 {
				t = t[:5]
			}
			if err := res[i].Close.Parse(t); err != nil {
				return err
			}
		}
	}
	*cdl = res
	return nil
}
