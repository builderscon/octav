package db

import "strconv"

func (v *RoomList) LoadForVenueSinceEID(tx *Tx, venueID string, since string, limit int) error {
	var s int64
	if id := since; id != "" {
		vdb := Room{}
		if err := vdb.LoadByEID(tx, id); err != nil {
			return err
		}

		s = vdb.OID
	}
	return v.LoadForVenueSince(tx, venueID, s, limit)
}

func (v *RoomList) LoadForVenueSince(tx *Tx, venueID string, since int64, limit int) error {
	rows, err := tx.Query(`SELECT oid, eid, venue_id, name, capacity, created_on, modified_on FROM `+RoomTable+` WHERE oid > ? AND venue_id = ? ORDER BY oid ASC LIMIT `+strconv.Itoa(limit), since, venueID)
	if err != nil {
		return err
	}
	res := make([]Room, 0, limit)
	for rows.Next() {
		vdb := Room{}
		if err := vdb.Scan(rows); err != nil {
			return err
		}
		res = append(res, vdb)
	}
	*v = res
	return nil
}
