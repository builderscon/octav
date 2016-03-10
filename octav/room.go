package octav

import (
	"log"

	"github.com/builderscon/octav/octav/db"
)

func (v Room) ToRow(vdb *db.Room) error {
	vdb.EID = v.ID
	vdb.VenueID = v.VenueID
	vdb.Name = v.Name
	vdb.Capacity = v.Capacity
	return nil
}

func (v *RoomList) LoadForVenue(tx *db.Tx, venueID, since string, limit int) error {
	vdbl := db.RoomList{}
	if err := vdbl.LoadForVenueSinceEID(tx, venueID, since, limit); err != nil {
		return err
	}

	log.Printf("%#v", vdbl)

	res := make([]Room, len(vdbl))
	for i, vdb := range vdbl {
		v := Room{}
		if err := v.FromRow(vdb); err != nil {
			return err
		}
		if err := v.LoadLocalizedFields(tx); err != nil {
			return err
		}
		res[i] = v
	}
	*v = res
	return nil
}
