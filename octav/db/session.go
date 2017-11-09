package db

import (
	"database/sql"
	"time"

	"github.com/builderscon/octav/octav/tools"
	pdebug "github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

var (
	sqlSessionIsSessionOwner string
)

func init() {
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	stmt.WriteString(`SELECT 1 FROM `)
	stmt.WriteString(SessionTable)
	stmt.WriteString(` WHERE eid = ? AND speaker_id = ?`)
	sqlSessionIsSessionOwner = stmt.String()
}

func (v *SessionList) LoadByConference(tx *sql.Tx, conferenceID, speakerID string, rangeStart, rangeEnd time.Time, status []string, confirmed []bool) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.SessionList.LoadByConference %s,%s,%s,%s,%s,%s", conferenceID, speakerID, rangeStart, rangeEnd, status, confirmed).BindError(&err)
		defer g.End()
	}

	// The caller of this method should ensure that query fields are
	// present and that we don't accidentally run an empty query
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	stmt.WriteString(`SELECT `)
	stmt.WriteString(SessionStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(SessionTable)
	stmt.WriteString(` WHERE `)

	where := tools.GetBuffer()
	defer tools.ReleaseBuffer(where)

	var args []interface{}
	if conferenceID != "" {
		where.WriteString(` conference_id = ?`)
		args = append(args, conferenceID)
	}

	if speakerID != "" {
		if where.Len() > 0 {
			where.WriteString(` AND`)
		}
		where.WriteString(` speaker_id = ?`)
		args = append(args, speakerID)
	}

	if !rangeStart.IsZero() && !rangeEnd.IsZero() {
		if where.Len() > 0 {
			where.WriteString(` AND`)
		}
		where.WriteString(` starts_on BETWEEN ? AND ?`)
		args = append(args, rangeStart, rangeEnd)
	}

	if l := len(status); l > 0 {
		if where.Len() > 0 {
			where.WriteString(` AND`)
		}
		where.WriteString(` status IN (`)
		for i, st := range status {
			where.WriteString(`?`)
			if i < l-1 {
				where.WriteString(`, `)
			}
			args = append(args, st)
		}
		where.WriteString(`)`)
	}

	if l := len(confirmed); l > 0 {
		if where.Len() > 0 {
			where.WriteString(` AND`)
		}
		where.WriteString(` confirmed IN (`)
		for i, c := range confirmed {
			where.WriteString(`?`)
			if i < l-1 {
				where.WriteString(`, `)
			}
			args = append(args, c)
		}
		where.WriteString(`)`)
	}

	if where.Len() > 0 {
		where.WriteString(` ORDER BY starts_on ASC`)
	}

	where.WriteTo(stmt)

	rows, err := Query(tx, stmt.String(), args...)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}
	defer rows.Close()

	if err := v.FromRows(rows, 0); err != nil {
		return err
	}
	return nil
}

func IsSessionOwner(tx *sql.Tx, sessionID, userID string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.IsSessionOwner %s %s", sessionID, userID).BindError(&err)
		defer g.End()
	}

	row, err := QueryRow(tx, sqlSessionIsSessionOwner, sessionID, userID)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}

	var r int

	if err := row.Scan(&r); err != nil {
		return errors.Wrap(err, "failed to scan from database")
	}

	if r == 0 {
		return errors.Errorf("user %s is not an owner of session %s", userID, sessionID)
	}
	return nil
}
