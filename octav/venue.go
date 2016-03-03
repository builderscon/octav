package octav

import (
	"github.com/builderscon/octav/octav/db"
)

func (v *VenueList) Load(tx *db.Tx, since string) error {
	var s uint64
	if id := since; id != "" {
		v := db.Venue{}
		if err := v.LoadByEID(tx, id); err != nil {
			return err
		}

		s = v.OID
	}

	rows, err := tx.Query(`SELECT eid, name FROM venue WHERE oid > ? ORDER BY oid LIMIT 10`, s)
	if err != nil {
		return err
	}

	// Not using db.Venue here
	res := make([]Venue, 0, 10)
	for rows.Next() {
		v := Venue{}
		if err := rows.Scan(&v.ID, &v.Name); err != nil {
			return err
		}
		res = append(res, v)
	}
	*v = res
	return nil
}
