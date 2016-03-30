package db

import "bytes"

func DeleteConferenceVenue(tx *Tx, cid, vid string) error {
	stmt := bytes.Buffer{}
	stmt.WriteString(`DELETE FROM `)
	stmt.WriteString(ConferenceVenueTable)
	stmt.WriteString(` WHERE conference_id = ? AND venue_id = ?`)

	_, err := tx.Exec(stmt.String(), cid, vid)
	return err
}

func LoadConferenceVenues(tx *Tx, venues *VenueList, cid string) error {
	stmt := bytes.Buffer{}
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

	rows, err := tx.Query(stmt.String(), cid)
	if err != nil {
		return err
	}

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