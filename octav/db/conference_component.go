package db

import (
	"database/sql"
	"time"

	"github.com/builderscon/octav/octav/tools"
	pdebug "github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

func DeleteConferenceComponentsByIDAndName(tx *sql.Tx, conferenceID string, names ...string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.DeleteConferenceComponentsByIDAndName conference = %s, names = %#v", conferenceID, names).BindError(&err)
		defer g.End()
	}

	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	l := len(names)
	if l == 0 {
		return errors.New("empty list of names")
	}

	stmt.WriteString(`DELETE FROM `)
	stmt.WriteString(ConferenceComponentTable)
	stmt.WriteString(` WHERE conference_id = ? AND name IN (`)
	for i := range names {
		stmt.WriteByte('?')
		if i < l-1 {
			stmt.WriteByte(',')
		}
	}
	stmt.WriteString(`)`)

	args := make([]interface{}, len(names)+1)
	args[0] = conferenceID
	for i, name := range names {
		args[i+1] = name
	}
	if _, err := Exec(tx, stmt.String(), args...); err != nil {
		return errors.Wrap(err, "failed to execute delete statement")
	}

	return nil
}

func UpsertConferenceComponentsByIDAndName(tx *sql.Tx, conferenceID string, values map[string]string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.UpsertConferenceComponentsByIDAndName conference = %s, values = %#v", conferenceID, values).BindError(&err)
		defer g.End()
	}
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	l := len(values)
	if l == 0 {
		return errors.New("empty value map")
	}

	stmt.WriteString(`INSERT INTO `)
	stmt.WriteString(ConferenceComponentTable)
	stmt.WriteString(` (eid, conference_id, name, value, created_on) VALUES `)

	var args []interface{}
	i := 0
	now := time.Now()
	for k, v := range values {
		args = append(args, tools.UUID(), conferenceID, k, v, now)
		stmt.WriteString(`(?, ?, ?, ?, ?)`)
		if i < l-1 {
			stmt.WriteByte(',')
		}
		i++
	}

	stmt.WriteString(` ON DUPLICATE KEY UPDATE value = VALUES(value)`)
	if _, err := Exec(tx, stmt.String(), args...); err != nil {
		return errors.Wrap(err, "failed to execute insert statement")
	}

	return nil
}

func (ccl *ConferenceComponentList) LoadByConferenceID(tx *sql.Tx, cid string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.ConferenceComponentList.LoadByaConferenceID conference = %s", cid).BindError(&err)
		defer g.End()
	}
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	stmt.WriteString(`SELECT `)
	stmt.WriteString(ConferenceComponentStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(ConferenceComponentTable)
	stmt.WriteString(` WHERE conference_id = ?`)

	rows, err := Query(tx, stmt.String(), cid)
	if err != nil {
		return errors.Wrap(err, "failed to execute query for conference component")
	}
	defer rows.Close()

	var res ConferenceComponentList
	for rows.Next() {
		var row ConferenceComponent
		if err := row.Scan(rows); err != nil {
			return errors.Wrap(err, "failed to scan conference component")
		}
		res = append(res, row)
	}

	*ccl = res
	return nil
}
