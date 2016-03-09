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

