package db

import (
	"database/sql"

	pdebug "github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

func IsConferenceSeriesAdministrator(tx *sql.Tx, sid, uid string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.IsConferenceSeriesAdministrator series %s, user %s", sid, uid).BindError(&err)
		defer g.End()
	}
	sqltext := `SELECT 1 FROM ` + ConferenceSeriesAdministratorTable + ` WHERE series_id = ? AND user_id = ?`

	var v int
	row, err := QueryRow(tx, sqltext, sid, uid)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}

	if err := row.Scan(&v); err != nil {
		return errors.Wrap(err, "failed to scan row")
	}
	if v != 1 {
		return errors.New("no matching administrator found")
	}
	return nil
}
