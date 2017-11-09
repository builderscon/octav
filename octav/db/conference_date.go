package db

import (
	"database/sql"

	"github.com/builderscon/octav/octav/tools"
	pdebug "github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

func (cd *ConferenceDate) DeleteDate(tx *sql.Tx, cid, eid string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.DeleteDate conference %s, date_id %s", cid, eid).BindError(&err)
		defer g.End()
	}

	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	stmt.WriteString("DELETE FROM ")
	stmt.WriteString(ConferenceDateTable)
	stmt.WriteString(" WHERE conference_id = ? AND eid = ?")
	_, err = Exec(tx, stmt.String(), cid, eid)
	return err
}

func (cdl *ConferenceDateList) LoadByConferenceID(tx *sql.Tx, cid string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.ConferenceDateList.LoadByConferenceID conference %s", cid).BindError(&err)
		defer g.End()
	}

	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	stmt.WriteString("SELECT ")
	stmt.WriteString(ConferenceDateStdSelectColumns)
	stmt.WriteString(" FROM ")
	stmt.WriteString(ConferenceDateTable)
	stmt.WriteString(" WHERE conference_id = ? ORDER BY open ASC")
	rows, err := Query(tx, stmt.String(), cid)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}
	defer rows.Close()

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
