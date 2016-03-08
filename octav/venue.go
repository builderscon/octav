package octav

import (
	"github.com/builderscon/octav/octav/db"
	"github.com/lestrrat/go-pdebug"
)

func (v *Venue) Delete(tx *db.Tx) error {
	if pdebug.Enabled {
		g := pdebug.Marker("Venue.Delete (%s)", v.ID)
		defer g.End()
	}

	vdb := db.Venue{EID: v.ID}
	if err := vdb.Delete(tx); err != nil {
		return err
	}

	if err := db.DeleteLocalizedStringsForParent(tx, v.ID, "Venue"); err != nil {
		return err
	}
	return nil
}

func (v *VenueList) Load(tx *db.Tx, since string) error {
	var s int64
	if id := since; id != "" {
		vdb := db.Venue{}
		if err := vdb.LoadByEID(tx, id); err != nil {
			return err
		}

		s = vdb.OID
	}

	rows, err := tx.Query(`SELECT eid, name FROM venues WHERE oid > ? ORDER BY oid LIMIT 10`, s)
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
