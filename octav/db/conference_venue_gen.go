package db

// Automatically generated by gendb utility. DO NOT EDIT!

import (
	"bytes"
	"database/sql"
	"time"

	"github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

const ConferenceVenueStdSelectColumns = "conference_venues.oid, conference_venues.conference_id, conference_venues.venue_id, conference_venues.created_on, conference_venues.modified_on"
const ConferenceVenueTable = "conference_venues"

type ConferenceVenueList []ConferenceVenue

func (c *ConferenceVenue) Scan(scanner interface {
	Scan(...interface{}) error
}) error {
	return scanner.Scan(&c.OID, &c.ConferenceID, &c.VenueID, &c.CreatedOn, &c.ModifiedOn)
}

func (c *ConferenceVenue) Create(tx *sql.Tx, opts ...InsertOption) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.ConferenceVenue.Create").BindError(&err)
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
	stmt.WriteString(ConferenceVenueTable)
	stmt.WriteString(` (conference_id, venue_id, created_on, modified_on) VALUES (?, ?, ?, ?)`)
	result, err := Exec(tx, stmt.String(), c.ConferenceID, c.VenueID, c.CreatedOn, c.ModifiedOn)
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

func (c ConferenceVenue) Update(tx *sql.Tx) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker(`ConferenceVenue.Update`).BindError(&err)
		defer g.End()
	}
	if c.OID != 0 {
		if pdebug.Enabled {
			pdebug.Printf(`Using OID (%d) as key`, c.OID)
		}
		const sqltext = `UPDATE conference_venues SET conference_id = ?, venue_id = ? WHERE oid = ?`
		if _, err := Exec(tx, sqltext, c.ConferenceID, c.VenueID, c.OID); err != nil {
			return errors.Wrap(err, `failed to execute statement`)
		}
		return nil
	}
	return errors.New("OID must be filled")
}

func (c ConferenceVenue) Delete(tx *sql.Tx) error {
	if c.OID != 0 {
		const sqltext = `DELETE FROM conference_venues WHERE oid = ?`
		if _, err := Exec(tx, sqltext, c.OID); err != nil {
			return errors.Wrap(err, `failed to execute statement`)
		}
		return nil
	}
	return errors.New("column OID must be filled")
}

func (v *ConferenceVenueList) FromRows(rows *sql.Rows, capacity int) error {
	var res []ConferenceVenue
	if capacity > 0 {
		res = make([]ConferenceVenue, 0, capacity)
	} else {
		res = []ConferenceVenue{}
	}

	for rows.Next() {
		vdb := ConferenceVenue{}
		if err := vdb.Scan(rows); err != nil {
			return err
		}
		res = append(res, vdb)
	}
	*v = res
	return nil
}
