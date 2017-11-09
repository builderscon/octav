package db

// Automatically generated by gendb utility. DO NOT EDIT!

import (
	"bytes"
	"database/sql"
	"strconv"

	"github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

const ExternalResourceStdSelectColumns = "external_resources.oid, external_resources.eid, external_resources.conference_id, external_resources.description, external_resources.image_url, external_resources.title, external_resources.url, external_resources.sort_order"
const ExternalResourceTable = "external_resources"

type ExternalResourceList []ExternalResource

func (e *ExternalResource) Scan(scanner interface {
	Scan(...interface{}) error
}) error {
	return scanner.Scan(&e.OID, &e.EID, &e.ConferenceID, &e.Description, &e.ImageURL, &e.Title, &e.URL, &e.SortOrder)
}

func (e *ExternalResource) LoadByEID(tx *sql.Tx, eid string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker(`ExternalResource.LoadByEID %s`, eid).BindError(&err)
		defer g.End()
	}
	const sqltext = `SELECT external_resources.oid, external_resources.eid, external_resources.conference_id, external_resources.description, external_resources.image_url, external_resources.title, external_resources.url, external_resources.sort_order FROM external_resources WHERE external_resources.eid = ?`
	row, err := QueryRow(tx, sqltext, eid)
	if err != nil {
		return errors.Wrap(err, `failed to query row`)
	}
	if err := e.Scan(row); err != nil {
		return err
	}
	return nil
}

func (e *ExternalResource) Create(tx *sql.Tx, opts ...InsertOption) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.ExternalResource.Create").BindError(&err)
		defer g.End()
		pdebug.Printf("%#v", e)
	}
	if e.EID == "" {
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
	stmt.WriteString(ExternalResourceTable)
	stmt.WriteString(` (eid, conference_id, description, image_url, title, url, sort_order) VALUES (?, ?, ?, ?, ?, ?, ?)`)
	result, err := Exec(tx, stmt.String(), e.EID, e.ConferenceID, e.Description, e.ImageURL, e.Title, e.URL, e.SortOrder)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}

	lii, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, `failed to fetch last insert ID`)
	}

	e.OID = lii
	return nil
}

func (e ExternalResource) Update(tx *sql.Tx) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker(`ExternalResource.Update`).BindError(&err)
		defer g.End()
	}
	if e.OID != 0 {
		if pdebug.Enabled {
			pdebug.Printf(`Using OID (%d) as key`, e.OID)
		}
		const sqltext = `UPDATE external_resources SET eid = ?, conference_id = ?, description = ?, image_url = ?, title = ?, url = ?, sort_order = ? WHERE oid = ?`
		if _, err := Exec(tx, sqltext, e.EID, e.ConferenceID, e.Description, e.ImageURL, e.Title, e.URL, e.SortOrder, e.OID); err != nil {
			return errors.Wrap(err, `failed to execute statement`)
		}
		return nil
	}

	if e.EID != "" {
		if pdebug.Enabled {
			pdebug.Printf(`Using EID (%s) as key`, e.EID)
		}
		const sqltext = `UPDATE external_resources SET eid = ?, conference_id = ?, description = ?, image_url = ?, title = ?, url = ?, sort_order = ? WHERE eid = ?`
		if _, err := Exec(tx, sqltext, e.EID, e.ConferenceID, e.Description, e.ImageURL, e.Title, e.URL, e.SortOrder, e.EID); err != nil {
			return errors.Wrap(err, `failed to execute statement`)
		}
		return nil
	}
	return errors.New("either OID/EID must be filled")
}

func (e ExternalResource) Delete(tx *sql.Tx) error {
	if e.OID != 0 {
		const sqltext = `DELETE FROM external_resources WHERE oid = ?`
		if _, err := Exec(tx, sqltext, e.OID); err != nil {
			return errors.Wrap(err, `failed to execute statement`)
		}
		return nil
	}

	if e.EID != "" {
		const sqltext = `DELETE FROM external_resources WHERE eid = ?`
		if _, err := Exec(tx, sqltext, e.EID); err != nil {
			return errors.Wrap(err, `failed to get statement`)
		}
		return nil
	}
	return errors.New("either OID/EID must be filled")
}

func (v *ExternalResourceList) FromRows(rows *sql.Rows, capacity int) error {
	var res []ExternalResource
	if capacity > 0 {
		res = make([]ExternalResource, 0, capacity)
	} else {
		res = []ExternalResource{}
	}

	for rows.Next() {
		vdb := ExternalResource{}
		if err := vdb.Scan(rows); err != nil {
			return err
		}
		res = append(res, vdb)
	}
	*v = res
	return nil
}

func (v *ExternalResourceList) LoadSinceEID(tx *sql.Tx, since string, limit int) error {
	var s int64
	if id := since; id != "" {
		vdb := ExternalResource{}
		if err := vdb.LoadByEID(tx, id); err != nil {
			return err
		}

		s = vdb.OID
	}
	return v.LoadSince(tx, s, limit)
}

func (v *ExternalResourceList) LoadSince(tx *sql.Tx, since int64, limit int) error {
	rows, err := Query(tx, `SELECT `+ExternalResourceStdSelectColumns+` FROM `+ExternalResourceTable+` WHERE external_resources.oid > ? ORDER BY oid ASC LIMIT `+strconv.Itoa(limit), since)
	if err != nil {
		return err
	}
	defer rows.Close()

	if err := v.FromRows(rows, limit); err != nil {
		return err
	}
	return nil
}
