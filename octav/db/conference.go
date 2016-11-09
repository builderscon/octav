package db

import (
	"fmt"
	"io"
	"time"

	"github.com/builderscon/octav/octav/tools"
	"github.com/pkg/errors"
)

func compileRangeWhere(dst io.Writer, args *[]interface{}, since int64, rangeStart, rangeEnd time.Time) error {
	where := tools.GetBuffer()
	defer tools.ReleaseBuffer(where)

	where.WriteString(ConferenceTable)
	where.WriteString(".oid > ?")
	*args = append(*args, since)

	if !rangeStart.IsZero() {
		if where.Len() > 0 {
			where.WriteString(" AND ")
		}
		where.WriteString(ConferenceDateTable)
		where.WriteString(".open >= ?")
		*args = append(*args, rangeStart)
	}

	if !rangeEnd.IsZero() {
		if where.Len() > 0 {
			where.WriteString(" AND ")
		}
		where.WriteString(ConferenceDateTable)
		where.WriteString(".open <= ?")
		*args = append(*args, rangeEnd)
	}

	toWrite := where.Len()
	n, err := where.WriteTo(dst)
	if n != int64(toWrite) {
		if err != nil {
			return errors.Wrap(err, "failed to write where clause to destination")
		}
		return errors.New("failed to write entire where clause to destination")
	}
	return nil
}

func (v *ConferenceList) LoadFromQuery(tx *Tx, status, organizerID []string, rangeStart, rangeEnd time.Time, since string, limit int) error {
	// We need the oid of "since"
	var sinceOID int64
	if since != "" {
		var vdb Conference
		if err := vdb.LoadByEID(tx, since); err != nil {
			return errors.Wrap(err, "failed to load reference row")
		}
		sinceOID = vdb.OID
	}

	qbuf := tools.GetBuffer()
	defer tools.ReleaseBuffer(qbuf)

	var hasDate bool
	if !rangeStart.IsZero() || !rangeEnd.IsZero() {
		hasDate = true
	}

	qbuf.WriteString(`SELECT `)
	qbuf.WriteString(ConferenceStdSelectColumns)
	qbuf.WriteString(` FROM `)
	qbuf.WriteString(ConferenceTable)
	if hasDate {
		qbuf.WriteString(` JOIN `)
		qbuf.WriteString(ConferenceDateTable)
		qbuf.WriteString(` ON `)
		qbuf.WriteString(ConferenceTable)
		qbuf.WriteString(`.eid = `)
		qbuf.WriteString(ConferenceDateTable)
		qbuf.WriteString(`.conference_id `)
	}
	if len(organizerID) > 0 {
		qbuf.WriteString(` JOIN `)
		qbuf.WriteString(ConferenceAdministratorTable)
		qbuf.WriteString(` ON `)
		qbuf.WriteString(ConferenceTable)
		qbuf.WriteString(`.eid = `)
		qbuf.WriteString(ConferenceAdministratorTable)
		qbuf.WriteString(`.conference_id `)
	}

	wherebuf := tools.GetBuffer()
	defer tools.ReleaseBuffer(wherebuf)

	var args []interface{}

	if !rangeStart.IsZero() || !rangeEnd.IsZero() {
		if err := compileRangeWhere(wherebuf, &args, sinceOID, rangeStart, rangeEnd); err != nil {
			return errors.Wrap(err, "failed to compile range where clause")
		}
	}

	if len(organizerID) > 0 {
		if wherebuf.Len() > 0 {
			wherebuf.WriteString(` AND `)
		}
		wherebuf.WriteString(ConferenceAdministratorTable)
		wherebuf.WriteString(`.user_id IN (`)
		for i, id := range organizerID {
			wherebuf.WriteByte('?')
			if i < len(organizerID) - 1 {
				wherebuf.WriteByte(',')
			}
			args = append(args, id)
		}
		wherebuf.WriteByte(')')
	}

	if len(status) > 0 {
		if wherebuf.Len() > 0 {
			wherebuf.WriteString(` AND `)
		}
		wherebuf.WriteString(ConferenceTable)
		wherebuf.WriteString(`.status IN (`)
		for i, st := range status {
			wherebuf.WriteByte('?')
			if i < len(status) - 1 {
				wherebuf.WriteByte(',')
			}
			args = append(args, st)
		}
		wherebuf.WriteByte(')')
	}

	if wherebuf.Len() > 0 {
		qbuf.WriteString(` WHERE `)
		wherebuf.WriteTo(qbuf)
	}

	qbuf.WriteString(` ORDER BY oid DESC`)
	fmt.Fprintf(qbuf, " LIMIT %d", limit)

	return v.execSQLAndExtract(tx, qbuf.String(), limit, args...)
}

func (v *ConferenceList) LoadByRange(tx *Tx, since string, rangeStart, rangeEnd time.Time, limit int) error {
	// We need the oid of "since"
	var sinceOID int64
	if since != "" {
		var vdb Conference
		if err := vdb.LoadByEID(tx, since); err != nil {
			return errors.Wrap(err, "failed to load reference row")
		}
		sinceOID = vdb.OID
	}

	// Use JOIN later
	var args []interface{}
	qbuf := tools.GetBuffer()
	defer tools.ReleaseBuffer(qbuf)

	var hasDate bool
	if !rangeStart.IsZero() || !rangeEnd.IsZero() {
		hasDate = true
	}

	qbuf.WriteString(`SELECT `)
	qbuf.WriteString(ConferenceStdSelectColumns)
	qbuf.WriteString(` FROM `)
	qbuf.WriteString(ConferenceTable)
	if hasDate {
		qbuf.WriteString(` JOIN `)
		qbuf.WriteString(ConferenceDateTable)
		qbuf.WriteString(` ON `)
		qbuf.WriteString(ConferenceTable)
		qbuf.WriteString(`.eid = `)
		qbuf.WriteString(ConferenceDateTable)
		qbuf.WriteString(`.conference_id `)
	}

	qbuf.WriteString(` WHERE `)
	if err := compileRangeWhere(qbuf, &args, sinceOID, rangeStart, rangeEnd); err != nil {
		return errors.Wrap(err, "failed to compile range where clause")
	}

	if hasDate {
		qbuf.WriteString(` ORDER BY `)
		qbuf.WriteString(ConferenceDateTable)
		qbuf.WriteString(`.date DESC`)
	} else {
		qbuf.WriteString(` ORDER BY oid DESC`)
	}

	fmt.Fprintf(qbuf, " LIMIT %d", limit)
	return v.execSQLAndExtract(tx, qbuf.String(), limit, args...)
}

func (v *ConferenceList) execSQLAndExtract(tx *Tx, sql string, limit int, args ...interface{}) error {
	rows, err := tx.Query(sql, args...)
	if err != nil {
		return err
	}

	res := make(ConferenceList, 0, limit)
	for rows.Next() {
		row := Conference{}
		if err := row.Scan(rows); err != nil {
			return err
		}
		res = append(res, row)
	}
	*v = res
	return nil
}

func ListConferencesByOrganizer(tx *Tx, l *ConferenceList, orgID string, statuses []string, since string, limit int) error {
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	stmt.WriteString(`SELECT `)
	stmt.WriteString(ConferenceStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(ConferenceTable)
	stmt.WriteString(` JOIN `)
	stmt.WriteString(ConferenceAdministratorTable)
	stmt.WriteString(` ON `)
	stmt.WriteString(ConferenceTable)
	stmt.WriteString(`.eid = `)
	stmt.WriteString(ConferenceAdministratorTable)
	stmt.WriteString(`.conference_id WHERE `)
	stmt.WriteString(ConferenceTable)

	if len(statuses) == 0 {
		statuses = []string{"public"}
	}
	stmt.WriteString(`.status IN (`)
	for i := range statuses {
		stmt.WriteByte('?')
		if i != len(statuses) - 1 {
			stmt.WriteByte(',')
		}
	}
	stmt.WriteString(`) AND `)
	stmt.WriteString(ConferenceAdministratorTable)
	stmt.WriteString(`.user_id = ? `)
	if since != "" {
		// Unimplemented
	}
	if limit > 0 {
		// Unimplemented
	}

	rows, err := tx.Query(stmt.String(), orgID)
	if err != nil {
		return errors.Wrap(err, "failed to execute query")
	}
	res := make(ConferenceList, 0, limit)
	for rows.Next() {
		row := Conference{}
		if err := row.Scan(rows); err != nil {
			return err
		}
		res = append(res, row)
	}
	*l = res
	return nil
}
