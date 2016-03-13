// Automatically generated by genmodel utility. DO NOT EDIT!
package model

import (
	"encoding/json"
	"time"

	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/tools"
	"github.com/lestrrat/go-pdebug"
)

var _ = time.Time{}

type RoomL10N struct {
	Room
	L10N tools.LocalizedFields `json:"-"`
}

func (v RoomL10N) MarshalJSON() ([]byte, error) {
	buf, err := json.Marshal(v.Room)
	if err != nil {
		return nil, err
	}
	return tools.MarshalJSONWithL10N(buf, v.L10N)
}

func (v *Room) Load(tx *db.Tx, id string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("Room.Load").BindError(&err)
		defer g.End()
	}
	vdb := db.Room{}
	if err := vdb.LoadByEID(tx, id); err != nil {
		return err
	}

	if err := v.FromRow(vdb); err != nil {
		return err
	}
	return nil
}

func (v *Room) FromRow(vdb db.Room) error {
	v.ID = vdb.EID
	v.VenueID = vdb.VenueID
	v.Name = vdb.Name
	v.Capacity = vdb.Capacity
	return nil
}

func (v *Room) ToRow(vdb *db.Room) error {
	vdb.EID = v.ID
	vdb.VenueID = v.VenueID
	vdb.Name = v.Name
	vdb.Capacity = v.Capacity
	return nil
}

func (v RoomL10N) GetPropNames() ([]string, error) {
	l, _ := v.L10N.GetPropNames()
	return append(l, "name"), nil
}

func (v RoomL10N) GetPropValue(s string) (interface{}, error) {
	switch s {
	case "id":
		return v.ID, nil
	case "venue_id":
		return v.VenueID, nil
	case "name":
		return v.Name, nil
	case "capacity":
		return v.Capacity, nil
	default:
		return v.L10N.GetPropValue(s)
	}
}

func (v *RoomL10N) UnmarshalJSON(data []byte) error {
	var s Room
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	v.Room = s
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	if err := tools.ExtractL10NFields(m, &v.L10N, []string{"name"}); err != nil {
		return err
	}

	return nil
}

func (v *RoomL10N) LoadLocalizedFields(tx *db.Tx) error {
	ls, err := db.LoadLocalizedStringsForParent(tx, v.Room.ID, "Room")
	if err != nil {
		return err
	}

	if len(ls) > 0 {
		v.L10N = tools.LocalizedFields{}
		for _, l := range ls {
			v.L10N.Set(l.Language, l.Name, l.Localized)
		}
	}
	return nil
}
