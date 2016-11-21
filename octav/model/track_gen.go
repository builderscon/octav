package model

// Automatically generated by genmodel utility. DO NOT EDIT!

import (
	"encoding/json"
	"time"

	"github.com/builderscon/octav/octav/db"
	"github.com/lestrrat/go-pdebug"
)

var _ = pdebug.Enabled
var _ = time.Time{}

type rawTrack struct {
	ID     string `json:"id"`
	RoomID string `json:"room_id"`
	Name   string `json:"name" l10n:"true"`
}

func (v Track) MarshalJSON() ([]byte, error) {
	var raw rawTrack
	raw.ID = v.ID
	raw.RoomID = v.RoomID
	raw.Name = v.Name
	buf, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	return MarshalJSONWithL10N(buf, v.LocalizedFields)
}

func (v *Track) Load(tx *db.Tx, id string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("model.Track.Load %s", id).BindError(&err)
		defer g.End()
	}
	vdb := db.Track{}
	if err := vdb.LoadByEID(tx, id); err != nil {
		return err
	}

	if err := v.FromRow(&vdb); err != nil {
		return err
	}
	return nil
}

func (v *Track) FromRow(vdb *db.Track) error {
	v.ID = vdb.EID
	v.RoomID = vdb.RoomID
	v.Name = vdb.Name
	return nil
}

func (v *Track) ToRow(vdb *db.Track) error {
	vdb.EID = v.ID
	vdb.RoomID = v.RoomID
	vdb.Name = v.Name
	return nil
}