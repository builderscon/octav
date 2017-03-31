package model

import (
	"database/sql"

	"github.com/builderscon/octav/octav/db"
	pdebug "github.com/lestrrat/go-pdebug"
)

func (v *Track) LoadByConferenceRoom(tx *sql.Tx, conferenceID, roomID string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("model.Track.LoadByConferenceRoom %s, %s", conferenceID, roomID).BindError(&err)
		defer g.End()
	}
	var vdb db.Track

	if err := vdb.Load(tx, conferenceID, roomID); err != nil {
		return err
	}

	if err := v.FromRow(&vdb); err != nil {
		return err
	}
	return nil
}
