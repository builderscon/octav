package db

// Automatically generated by gendb utility. DO NOT EDIT!

import (
	"bytes"
	"database/sql"
	"strconv"
	"time"

	"github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

const BlogEntryStdSelectColumns = "blog_entries.oid, blog_entries.eid, blog_entries.conference_id, blog_entries.title, blog_entries.url, blog_entries.url_hash, blog_entries.status, blog_entries.created_on, blog_entries.modified_on"
const BlogEntryTable = "blog_entries"

type BlogEntryList []BlogEntry

func (b *BlogEntry) Scan(scanner interface {
	Scan(...interface{}) error
}) error {
	return scanner.Scan(&b.OID, &b.EID, &b.ConferenceID, &b.Title, &b.URL, &b.URLHash, &b.Status, &b.CreatedOn, &b.ModifiedOn)
}

func (b *BlogEntry) LoadByEID(tx *sql.Tx, eid string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker(`BlogEntry.LoadByEID %s`, eid).BindError(&err)
		defer g.End()
	}
	const sqltext = `SELECT blog_entries.oid, blog_entries.eid, blog_entries.conference_id, blog_entries.title, blog_entries.url, blog_entries.url_hash, blog_entries.status, blog_entries.created_on, blog_entries.modified_on FROM blog_entries WHERE blog_entries.eid = ?`
	row, err := QueryRow(tx, sqltext, eid)
	if err != nil {
		return errors.Wrap(err, `failed to query row`)
	}
	if err := b.Scan(row); err != nil {
		return err
	}
	return nil
}

func (b *BlogEntry) Create(tx *sql.Tx, opts ...InsertOption) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.BlogEntry.Create").BindError(&err)
		defer g.End()
		pdebug.Printf("%#v", b)
	}
	if b.EID == "" {
		return errors.New("create: non-empty EID required")
	}

	b.CreatedOn = time.Now()
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
	stmt.WriteString(BlogEntryTable)
	stmt.WriteString(` (eid, conference_id, title, url, url_hash, status, created_on, modified_on) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	result, err := Exec(tx, stmt.String(), b.EID, b.ConferenceID, b.Title, b.URL, b.URLHash, b.Status, b.CreatedOn, b.ModifiedOn)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}

	lii, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, `failed to fetch last insert ID`)
	}

	b.OID = lii
	return nil
}

func (b BlogEntry) Update(tx *sql.Tx) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker(`BlogEntry.Update`).BindError(&err)
		defer g.End()
	}
	if b.OID != 0 {
		if pdebug.Enabled {
			pdebug.Printf(`Using OID (%d) as key`, b.OID)
		}
		const sqltext = `UPDATE blog_entries SET eid = ?, conference_id = ?, title = ?, url = ?, url_hash = ?, status = ? WHERE oid = ?`
		if _, err := Exec(tx, sqltext, b.EID, b.ConferenceID, b.Title, b.URL, b.URLHash, b.Status, b.OID); err != nil {
			return errors.Wrap(err, `failed to execute statement`)
		}
		return nil
	}

	if b.EID != "" {
		if pdebug.Enabled {
			pdebug.Printf(`Using EID (%s) as key`, b.EID)
		}
		const sqltext = `UPDATE blog_entries SET eid = ?, conference_id = ?, title = ?, url = ?, url_hash = ?, status = ? WHERE eid = ?`
		if _, err := Exec(tx, sqltext, b.EID, b.ConferenceID, b.Title, b.URL, b.URLHash, b.Status, b.EID); err != nil {
			return errors.Wrap(err, `failed to execute statement`)
		}
		return nil
	}
	return errors.New("either OID/EID must be filled")
}

func (b BlogEntry) Delete(tx *sql.Tx) error {
	if b.OID != 0 {
		const sqltext = `DELETE FROM blog_entries WHERE oid = ?`
		if _, err := Exec(tx, sqltext, b.OID); err != nil {
			return errors.Wrap(err, `failed to execute statement`)
		}
		return nil
	}

	if b.EID != "" {
		const sqltext = `DELETE FROM blog_entries WHERE eid = ?`
		if _, err := Exec(tx, sqltext, b.EID); err != nil {
			return errors.Wrap(err, `failed to get statement`)
		}
		return nil
	}
	return errors.New("either OID/EID must be filled")
}

func (v *BlogEntryList) FromRows(rows *sql.Rows, capacity int) error {
	var res []BlogEntry
	if capacity > 0 {
		res = make([]BlogEntry, 0, capacity)
	} else {
		res = []BlogEntry{}
	}

	for rows.Next() {
		vdb := BlogEntry{}
		if err := vdb.Scan(rows); err != nil {
			return err
		}
		res = append(res, vdb)
	}
	*v = res
	return nil
}

func (v *BlogEntryList) LoadSinceEID(tx *sql.Tx, since string, limit int) error {
	var s int64
	if id := since; id != "" {
		vdb := BlogEntry{}
		if err := vdb.LoadByEID(tx, id); err != nil {
			return err
		}

		s = vdb.OID
	}
	return v.LoadSince(tx, s, limit)
}

func (v *BlogEntryList) LoadSince(tx *sql.Tx, since int64, limit int) error {
	rows, err := Query(tx, `SELECT `+BlogEntryStdSelectColumns+` FROM `+BlogEntryTable+` WHERE blog_entries.oid > ? ORDER BY oid ASC LIMIT `+strconv.Itoa(limit), since)
	if err != nil {
		return err
	}
	defer rows.Close()

	if err := v.FromRows(rows, limit); err != nil {
		return err
	}
	return nil
}
