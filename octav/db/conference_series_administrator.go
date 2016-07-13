package db

import (
	"github.com/pkg/errors"
)

func IsConferenceSeriesAdministrator(tx *Tx, sid, uid string) error {
	stmt := getStmtBuf()
	defer releaseStmtBuf(stmt)
	stmt.WriteString(`SELECT 1 FROM `)
	stmt.WriteString(ConferenceSeriesAdministratorTable)
	stmt.WriteString(` WHERE series_id = ? AND user_id = ?`)

	var v int
	row := tx.QueryRow(stmt.String(), sid, uid)
	if err := row.Scan(&v); err != nil {
		return errors.Wrap(err, "failed to scan row")
	}
	if v != 1 {
		return errors.New("no matching administrator found")
	}
	return nil
}


