package db

// Automatically generated by gendb utility. DO NOT EDIT!

import (
	"bytes"
	"database/sql"
	"time"

	"github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

const ConferenceSeriesAdministratorStdSelectColumns = "conference_series_administrators.oid, conference_series_administrators.series_id, conference_series_administrators.user_id, conference_series_administrators.created_on, conference_series_administrators.modified_on"
const ConferenceSeriesAdministratorTable = "conference_series_administrators"

type ConferenceSeriesAdministratorList []ConferenceSeriesAdministrator

func (c *ConferenceSeriesAdministrator) Scan(scanner interface {
	Scan(...interface{}) error
}) error {
	return scanner.Scan(&c.OID, &c.SeriesID, &c.UserID, &c.CreatedOn, &c.ModifiedOn)
}

func (c *ConferenceSeriesAdministrator) Create(tx *sql.Tx, opts ...InsertOption) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.ConferenceSeriesAdministrator.Create").BindError(&err)
		defer g.End()
		pdebug.Printf("%#v", c)
	}
	c.CreatedOn = time.Now()
	doIgnore := false
	for _, opt := range opts {
		switch opt.(type) {
		case insertIgnoreOption:
			doIgnore = true
		}
	}

	stmt := bytes.Buffer{}
	stmt.WriteString("INSERT ")
	if doIgnore {
		stmt.WriteString("IGNORE ")
	}
	stmt.WriteString("INTO ")
	stmt.WriteString(ConferenceSeriesAdministratorTable)
	stmt.WriteString(` (series_id, user_id, created_on, modified_on) VALUES (?, ?, ?, ?)`)
	result, err := Exec(tx, stmt.String(), c.SeriesID, c.UserID, c.CreatedOn, c.ModifiedOn)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}

	lii, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, `failed to fetch last insert ID`)
	}

	c.OID = lii
	return nil
}

func (c ConferenceSeriesAdministrator) Update(tx *sql.Tx) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker(`ConferenceSeriesAdministrator.Update`).BindError(&err)
		defer g.End()
	}
	if c.OID != 0 {
		if pdebug.Enabled {
			pdebug.Printf(`Using OID (%d) as key`, c.OID)
		}
		const sqltext = `UPDATE conference_series_administrators SET series_id = ?, user_id = ? WHERE oid = ?`
		if _, err := Exec(tx, sqltext, c.SeriesID, c.UserID, c.OID); err != nil {
			return errors.Wrap(err, `failed to execute statement`)
		}
		return nil
	}
	return errors.New("OID must be filled")
}

func (c ConferenceSeriesAdministrator) Delete(tx *sql.Tx) error {
	if c.OID != 0 {
		const sqltext = `DELETE FROM conference_series_administrators WHERE oid = ?`
		if _, err := Exec(tx, sqltext, c.OID); err != nil {
			return errors.Wrap(err, `failed to execute statement`)
		}
		return nil
	}
	return errors.New("column OID must be filled")
}

func (v *ConferenceSeriesAdministratorList) FromRows(rows *sql.Rows, capacity int) error {
	var res []ConferenceSeriesAdministrator
	if capacity > 0 {
		res = make([]ConferenceSeriesAdministrator, 0, capacity)
	} else {
		res = []ConferenceSeriesAdministrator{}
	}

	for rows.Next() {
		vdb := ConferenceSeriesAdministrator{}
		if err := vdb.Scan(rows); err != nil {
			return err
		}
		res = append(res, vdb)
	}
	*v = res
	return nil
}
