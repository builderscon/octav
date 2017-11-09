package db

// Automatically generated by gendb utility. DO NOT EDIT!

import (
	"bytes"
	"database/sql"

	"github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

const TemporaryEmailStdSelectColumns = "temporary_emails.oid, temporary_emails.user_id, temporary_emails.confirmation_key, temporary_emails.email, temporary_emails.expires_on"
const TemporaryEmailTable = "temporary_emails"

type TemporaryEmailList []TemporaryEmail

func (t *TemporaryEmail) Scan(scanner interface {
	Scan(...interface{}) error
}) error {
	return scanner.Scan(&t.OID, &t.UserID, &t.ConfirmationKey, &t.Email, &t.ExpiresOn)
}

func (t *TemporaryEmail) Create(tx *sql.Tx, opts ...InsertOption) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.TemporaryEmail.Create").BindError(&err)
		defer g.End()
		pdebug.Printf("%#v", t)
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
	stmt.WriteString(TemporaryEmailTable)
	stmt.WriteString(` (user_id, confirmation_key, email, expires_on) VALUES (?, ?, ?, ?)`)
	result, err := Exec(tx, stmt.String(), t.UserID, t.ConfirmationKey, t.Email, t.ExpiresOn)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}

	lii, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, `failed to fetch last insert ID`)
	}

	t.OID = lii
	return nil
}

func (t TemporaryEmail) Update(tx *sql.Tx) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker(`TemporaryEmail.Update`).BindError(&err)
		defer g.End()
	}
	if t.OID != 0 {
		if pdebug.Enabled {
			pdebug.Printf(`Using OID (%d) as key`, t.OID)
		}
		const sqltext = `UPDATE temporary_emails SET user_id = ?, confirmation_key = ?, email = ?, expires_on = ? WHERE oid = ?`
		if _, err := Exec(tx, sqltext, t.UserID, t.ConfirmationKey, t.Email, t.ExpiresOn, t.OID); err != nil {
			return errors.Wrap(err, `failed to execute statement`)
		}
		return nil
	}
	return errors.New("OID must be filled")
}

func (t TemporaryEmail) Delete(tx *sql.Tx) error {
	if t.OID != 0 {
		const sqltext = `DELETE FROM temporary_emails WHERE oid = ?`
		if _, err := Exec(tx, sqltext, t.OID); err != nil {
			return errors.Wrap(err, `failed to execute statement`)
		}
		return nil
	}
	return errors.New("column OID must be filled")
}

func (v *TemporaryEmailList) FromRows(rows *sql.Rows, capacity int) error {
	var res []TemporaryEmail
	if capacity > 0 {
		res = make([]TemporaryEmail, 0, capacity)
	} else {
		res = []TemporaryEmail{}
	}

	for rows.Next() {
		vdb := TemporaryEmail{}
		if err := vdb.Scan(rows); err != nil {
			return err
		}
		res = append(res, vdb)
	}
	*v = res
	return nil
}
