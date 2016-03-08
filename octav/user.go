package octav

import (
	"github.com/builderscon/octav/octav/db"
	"github.com/lestrrat/go-pdebug"
)

func (v *User) Delete(tx *db.Tx) error {
	if pdebug.Enabled {
		g := pdebug.Marker("User.Delete (%s)", v.ID)
		defer g.End()
	}

	vdb := db.User{EID: v.ID}
	if err := vdb.Delete(tx); err != nil {
		return err
	}

	if err := db.DeleteLocalizedStringsForParent(tx, v.ID, "User"); err != nil {
		return err
	}
	return nil
}
