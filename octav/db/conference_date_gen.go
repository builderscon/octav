package db

// Automatically generated by gendb utility. DO NOT EDIT!

import (
	"bytes"
	"database/sql"
	"strconv"

	"github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

const ConferenceDateStdSelectColumns = "conference_dates.oid, conference_dates.eid, conference_dates.conference_id, conference_dates.open, conference_dates.close"
const ConferenceDateTable = "conference_dates"

type ConferenceDateList []ConferenceDate

func (c *ConferenceDate) Scan(scanner interface {
	Scan(...interface{}) error
}) error {
	return scanner.Scan(&c.OID, &c.EID, &c.ConferenceID, &c.Open, &c.Close)
}

func (c *ConferenceDate) LoadByEID(tx *sql.Tx, eid string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker(`ConferenceDate.LoadByEID %s`, eid).BindError(&err)
		defer g.End()
	}
	const sqltext = `SELECT conference_dates.oid, conference_dates.eid, conference_dates.conference_id, conference_dates.open, conference_dates.close FROM conference_dates WHERE conference_dates.eid = ?`
	row, err := QueryRow(tx, sqltext, eid)
	if err != nil {
		return errors.Wrap(err, `failed to query row`)
	}
	if err := c.Scan(row); err != nil {
		return err
	}
	return nil
}

func (c *ConferenceDate) Create(tx *sql.Tx, opts ...InsertOption) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.ConferenceDate.Create").BindError(&err)
		defer g.End()
		pdebug.Printf("%#v", c)
	}
	if c.EID == "" {
		return errors.New("create: non-empty EID required")
	}

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
	stmt.WriteString(ConferenceDateTable)
	stmt.WriteString(` (eid, conference_id, open, close) VALUES (?, ?, ?, ?)`)
	result, err := Exec(tx, stmt.String(), c.EID, c.ConferenceID, c.Open, c.Close)
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

func (c ConferenceDate) Update(tx *sql.Tx) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker(`ConferenceDate.Update`).BindError(&err)
		defer g.End()
	}
	if c.OID != 0 {
		if pdebug.Enabled {
			pdebug.Printf(`Using OID (%d) as key`, c.OID)
		}
		const sqltext = `UPDATE conference_dates SET eid = ?, conference_id = ?, open = ?, close = ? WHERE oid = ?`
		if _, err := Exec(tx, sqltext, c.EID, c.ConferenceID, c.Open, c.Close, c.OID); err != nil {
			return errors.Wrap(err, `failed to execute statement`)
		}
		return nil
	}

	if c.EID != "" {
		if pdebug.Enabled {
			pdebug.Printf(`Using EID (%s) as key`, c.EID)
		}
		const sqltext = `UPDATE conference_dates SET eid = ?, conference_id = ?, open = ?, close = ? WHERE eid = ?`
		if _, err := Exec(tx, sqltext, c.EID, c.ConferenceID, c.Open, c.Close, c.EID); err != nil {
			return errors.Wrap(err, `failed to execute statement`)
		}
		return nil
	}
	return errors.New("either OID/EID must be filled")
}

func (c ConferenceDate) Delete(tx *sql.Tx) error {
	if c.OID != 0 {
		const sqltext = `DELETE FROM conference_dates WHERE oid = ?`
		if _, err := Exec(tx, sqltext, c.OID); err != nil {
			return errors.Wrap(err, `failed to execute statement`)
		}
		return nil
	}

	if c.EID != "" {
		const sqltext = `DELETE FROM conference_dates WHERE eid = ?`
		if _, err := Exec(tx, sqltext, c.EID); err != nil {
			return errors.Wrap(err, `failed to get statement`)
		}
		return nil
	}
	return errors.New("either OID/EID must be filled")
}

func (v *ConferenceDateList) FromRows(rows *sql.Rows, capacity int) error {
	var res []ConferenceDate
	if capacity > 0 {
		res = make([]ConferenceDate, 0, capacity)
	} else {
		res = []ConferenceDate{}
	}

	for rows.Next() {
		vdb := ConferenceDate{}
		if err := vdb.Scan(rows); err != nil {
			return err
		}
		res = append(res, vdb)
	}
	*v = res
	return nil
}

func (v *ConferenceDateList) LoadSinceEID(tx *sql.Tx, since string, limit int) error {
	var s int64
	if id := since; id != "" {
		vdb := ConferenceDate{}
		if err := vdb.LoadByEID(tx, id); err != nil {
			return err
		}

		s = vdb.OID
	}
	return v.LoadSince(tx, s, limit)
}

func (v *ConferenceDateList) LoadSince(tx *sql.Tx, since int64, limit int) error {
	rows, err := Query(tx, `SELECT `+ConferenceDateStdSelectColumns+` FROM `+ConferenceDateTable+` WHERE conference_dates.oid > ? ORDER BY oid ASC LIMIT `+strconv.Itoa(limit), since)
	if err != nil {
		return err
	}
	defer rows.Close()

	if err := v.FromRows(rows, limit); err != nil {
		return err
	}
	return nil
}
