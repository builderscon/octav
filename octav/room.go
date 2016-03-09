package octav

import (
	"github.com/builderscon/octav/octav/db"
)

func (v Room) ToRow(vdb *db.Room) error {
	vdb.EID = v.ID
	vdb.VenueID = v.VenueID
	vdb.Name = v.Name
	vdb.Capacity = v.Capacity
	return nil
}
