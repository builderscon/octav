package db

import (
	"database/sql"

	"github.com/builderscon/octav/octav/tools"
)

func (cd *ConferenceDate) DeleteDate(tx *sql.Tx, cid, eid string) error {
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	stmt.WriteString("DELETE FROM ")
	stmt.WriteString(ConferenceDateTable)
	stmt.WriteString(" WHERE conference_id = ? AND eid = ?")
	_, err := tx.Exec(stmt.String(), cid, eid)
	return err
}

func (cdl *ConferenceDateList) LoadByConferenceID(tx *sql.Tx, cid string) error {
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	stmt.WriteString("SELECT ")
	stmt.WriteString(ConferenceDateStdSelectColumns)
	stmt.WriteString(" FROM ")
	stmt.WriteString(ConferenceDateTable)
	stmt.WriteString(" WHERE conference_id = ? ORDER BY open ASC")
	rows, err := tx.Query(stmt.String(), cid)
	if err != nil {
		return err
	}

	res := ConferenceDateList{}
	for rows.Next() {
		var cd ConferenceDate
		if err := cd.Scan(rows); err != nil {
			return err
		}
		res = append(res, cd)
	}
	*cdl = res
	return nil
}
