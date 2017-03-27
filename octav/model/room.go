package model

import (
	"database/sql"

	"github.com/builderscon/octav/octav/db"
)

func (v *RoomList) LoadForVenue(tx *sql.Tx, venueID, since string, limit int) error {
	vdbl := db.RoomList{}
	if err := vdbl.LoadForVenueSinceEID(tx, venueID, since, limit); err != nil {
		return err
	}

	res := make([]Room, len(vdbl))
	for i, vdb := range vdbl {
		if err := res[i].FromRow(&vdb); err != nil {
			return err
		}
	}
	*v = res
	return nil
}
