package db

import (
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
)

func compileRangeWhere(dst io.Writer, args *[]interface{}, since int64, rangeStart, rangeEnd time.Time) error {
	where := getStmtBuf()
	defer releaseStmtBuf(where)

	where.WriteString(ConferenceTable)
	where.WriteString(".oid > ?")
	*args = append(*args, since)

	if !rangeStart.IsZero() {
		if where.Len() > 0 {
			where.WriteString(" AND ")
		}
		where.WriteString(ConferenceDateTable)
		where.WriteString(".date >= ?")
		*args = append(*args, rangeStart)
	}

	if !rangeEnd.IsZero() {
		if where.Len() > 0 {
			where.WriteString(" AND ")
		}
		where.WriteString(ConferenceDateTable)
		where.WriteString(".date <= ?")
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

func (v *ConferenceList) LoadByStatusAndRange(tx *Tx, status string, since string, rangeStart, rangeEnd time.Time, limit int) error {
	// We need the oid of "since"
	var sinceOID int64
	if since != "" {
		var vdb Conference
		if err := vdb.LoadByEID(tx, since); err != nil {
			return errors.Wrap(err, "failed to load reference row")
		}
		sinceOID = vdb.OID
	}

	qbuf := getStmtBuf()
	defer releaseStmtBuf(qbuf)

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

	var args []interface{}
	if err := compileRangeWhere(qbuf, &args, sinceOID, rangeStart, rangeEnd); err != nil {
		return errors.Wrap(err, "failed to compile range where clause")
	}
	qbuf.WriteString(` AND `)
	qbuf.WriteString(ConferenceTable)
	qbuf.WriteString(`.status = ?`)
	args = append(args, status)

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
	qbuf := getStmtBuf()
	defer releaseStmtBuf(qbuf)

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

func ListConferencesByOrganizer(tx *Tx, l *ConferenceList, orgID, since string, limit int) error {
	stmt := getStmtBuf()
	defer releaseStmtBuf(stmt)

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
	stmt.WriteString(`.status != "private" AND `)
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
