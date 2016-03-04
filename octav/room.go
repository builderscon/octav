package octav

import (
	"github.com/builderscon/octav/octav/db"
)

func (v *Room) Load(tx *db.Tx, id string) error {
	vdb := db.Room{}
	if err := vdb.LoadByEID(tx, id); err != nil {
		return err
	}

	return v.FromRow(vdb)
}

func (v *Room) FromRow(vdb db.Room) error {
	v.ID = vdb.EID
	v.VenueID = vdb.VenueID
	v.Name = vdb.Name
	v.Capacity = vdb.Capacity
	return nil
}

func (v Room) ToRow(vdb *db.Room) error {
	vdb.EID = v.ID
	vdb.VenueID = v.VenueID
	vdb.Name = v.Name
	vdb.Capacity = v.Capacity
	return nil
}

func (v *RoomList) Load(tx *db.Tx, since string) error {
	var s int64
	if id := since; id != "" {
		vdb := db.Room{}
		if err := vdb.LoadByEID(tx, id); err != nil {
			return err
		}

		s = vdb.OID
	}

	rows, err := tx.Query(`SELECT eid, name FROM venue WHERE oid > ? ORDER BY oid LIMIT 10`, s)
	if err != nil {
		return err
	}

	// Not using db.Room here
	res := make([]Room, 0, 10)
	for rows.Next() {
		vdb := db.Room{}
		if err := vdb.Scan(rows); err != nil {
			return err
		}
		v := Room{}
		if err := v.FromRow(vdb); err != nil {
			return err
		}
		res = append(res, v)
	}
	*v = res
	return nil
}