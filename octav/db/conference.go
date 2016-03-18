package db

import (
	"bytes"
	"fmt"
	"time"
)

func (v *ConferenceList) LoadByRange(tx *Tx, since string, rangeStart, rangeEnd time.Time, limit int) error {
	// Use JOIN later
	var args []interface{}
	where := bytes.Buffer{}
	if since != "" {
		vdb := Conference{}
		if err := vdb.LoadByEID(tx, since); err != nil {
			return err
		}
		where.WriteString("oid > ?")
		args = append(args, vdb.OID)
	}

	hasDate := false
	if !rangeStart.IsZero() {
		if where.Len() > 0 {
			where.WriteString(" AND ")
		}
		where.WriteString("conference_dates.date >= ?")
		args = append(args, rangeStart)
		hasDate = true
	}

	if !rangeEnd.IsZero() {
		if where.Len() > 0 {
			where.WriteString(" AND ")
		}
		where.WriteString("conference_dates.date <= ?")
		args = append(args, rangeEnd)
		hasDate = true
	}

	qbuf := bytes.Buffer{}
	qbuf.WriteString(`SELECT `)
	qbuf.WriteString(ConferenceStdSelectColumns)
	qbuf.WriteString(` FROM `)
	qbuf.WriteString(ConferenceTable)
	if hasDate {
		qbuf.WriteString(` JOIN conference_dates ON conferences.eid = conference_dates.conference_id `)
	}

	if where.Len() > 0 {
		qbuf.WriteString(` WHERE `)
		where.WriteTo(&qbuf)
	}

	if hasDate {
		qbuf.WriteString(` ORDER BY conference_dates.date DESC`)
	} else {
		qbuf.WriteString(` ORDER BY oid DESC`)
	}

	fmt.Fprintf(&qbuf, " LIMIT %d", limit)

	rows, err := tx.Query(qbuf.String(), args...)
	if err != nil {
		return err
	}

	res := make([]Conference, 0, limit)
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
