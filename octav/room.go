package octav

import (
	"github.com/builderscon/octav/octav/db"
	"github.com/lestrrat/go-pdebug"
)

func (v *Room) Delete(tx *db.Tx) error {
	if pdebug.Enabled {
		g := pdebug.Marker("Room.Delete (%s)", v.ID)
		defer g.End()
	}

	vdb := db.Room{EID: v.ID}
	if err := vdb.Delete(tx); err != nil {
		return err
	}

	if err := db.DeleteLocalizedStringsForParent(tx, v.ID, "Room"); err != nil {
		return err
	}
	return nil
}

func (v Room) ToRow(vdb *db.Room) error {
	vdb.EID = v.ID
	vdb.VenueID = v.VenueID
	vdb.Name = v.Name
	vdb.Capacity = v.Capacity
	return nil
}
