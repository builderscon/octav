package db

import (
	"database/sql"

	"github.com/builderscon/octav/octav/tools"
	pdebug "github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

var (
	sqlConferenceVenueDelete string
	sqlConferenceVenueLoad   string
)

func init() {
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	stmt.WriteString(`DELETE FROM `)
	stmt.WriteString(ConferenceVenueTable)
	stmt.WriteString(` WHERE conference_id = ? AND venue_id = ?`)
	sqlConferenceVenueDelete = stmt.String()

	stmt.Reset()
	stmt.WriteString(`SELECT `)
	stmt.WriteString(VenueStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(ConferenceVenueTable)
	stmt.WriteString(` JOIN `)
	stmt.WriteString(VenueTable)
	stmt.WriteString(` ON `)
	stmt.WriteString(ConferenceVenueTable)
	stmt.WriteString(`.venue_id = `)
	stmt.WriteString(VenueTable)
	stmt.WriteString(`.eid WHERE `)
	stmt.WriteString(ConferenceVenueTable)
	stmt.WriteString(`.conference_id = ?`)
	sqlConferenceVenueLoad = stmt.String()
}

func DeleteConferenceVenue(tx *sql.Tx, cid, vid string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.DeleteConferenceVenues conference %s, venue %s", cid, vid).BindError(&err)
		defer g.End()
	}
	_, err = Exec(tx, sqlConferenceVenueDelete, cid, vid)
	return err
}

func LoadConferenceVenues(tx *sql.Tx, venues *VenueList, cid string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.LoadConferenceVenues %s", cid).BindError(&err)
		defer g.End()
	}

	rows, err := Query(tx, sqlConferenceVenueLoad, cid)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}
	defer rows.Close()

	var res VenueList
	for rows.Next() {
		var u Venue
		if err := u.Scan(rows); err != nil {
			return err
		}

		res = append(res, u)
	}

	*venues = res
	return nil
}
