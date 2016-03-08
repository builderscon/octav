package octav

import (
	"encoding/json"

	"github.com/builderscon/octav/octav/db"
	"github.com/lestrrat/go-pdebug"
)

func (v Room) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["id"] = v.ID
	m["venue_id"] = v.VenueID
	m["name"] = v.Name
	m["capacity"] = v.Capacity

	buf, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return marshalJSONWithL10N(buf, v.L10N)
}

func (v *Room) UnmarshalJSON(data []byte) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if jv, ok := m["id"]; ok {
		switch jv.(type) {
		case string:
			v.ID = jv.(string)
			delete(m, "id")
		default:
			return ErrInvalidFieldType
		}
	}
	if jv, ok := m["venue_id"]; ok {
		switch jv.(type) {
		case string:
			v.VenueID = jv.(string)
			delete(m, "venue_id")
		default:
			return ErrInvalidFieldType
		}
	}
	if jv, ok := m["name"]; ok {
		switch jv.(type) {
		case string:
			v.Name = jv.(string)
			delete(m, "name")
		default:
			return ErrInvalidFieldType
		}
	}
	if jv, ok := m["capacity"]; ok {
		switch jv.(type) {
		case float64:
			v.Capacity = uint(jv.(float64))
			delete(m, "capacity")
		default:
			return ErrInvalidFieldType
		}
	}

	if err := ExtractL10NFields(m, &v.L10N, []string{"name"}); err != nil {
		return err
	}

	return nil
}

func (v *Room) Load(tx *db.Tx, id string) error {
	vdb := db.Room{}
	if err := vdb.LoadByEID(tx, id); err != nil {
		return err
	}

	v.FromRow(vdb)
	ls, err := db.LoadLocalizedStringsForParent(tx, v.ID, "Room")
	if err != nil {
		return err
	}

	if len(ls) > 0 {
		v.L10N = LocalizedFields{}
		for _, l := range ls {
			v.L10N.Set(l.Language, l.Name, l.Localized)
		}
	}

	return v.FromRow(vdb)
}

func (v *Room) Create(tx *db.Tx) error {
	if v.ID == "" {
		v.ID = UUID()
	}

	vdb := db.Room{
		EID:     v.ID,
		Name:    v.Name,
		VenueID: v.VenueID,
		Capacity: v.Capacity,
	}
	if err := vdb.Create(tx); err != nil {
		return err
	}

	if err := v.L10N.CreateLocalizedStrings(tx, "Room", v.ID); err != nil {
		return err
	}
	return nil
}

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

func (v *Room) FromRow(vdb db.Room) error {
	v.ID = vdb.EID
	v.VenueID = vdb.VenueID
	v.Name = vdb.Name
	v.Capacity = vdb.Capacity
	return nil
}

func (v Room) ToRow(vdb *db.Room) error {
	vdb.EID = v.ID
	vdb.VenueID = v.VenueID
	vdb.Name = v.Name
	vdb.Capacity = v.Capacity
	return nil
}

func (v *RoomList) Load(tx *db.Tx, since string) error {
	var s int64
	if id := since; id != "" {
		vdb := db.Room{}
		if err := vdb.LoadByEID(tx, id); err != nil {
			return err
		}

		s = vdb.OID
	}

	rows, err := tx.Query(`SELECT eid, name FROM venue WHERE oid > ? ORDER BY oid LIMIT 10`, s)
	if err != nil {
		return err
	}

	// Not using db.Room here
	res := make([]Room, 0, 10)
	for rows.Next() {
		vdb := db.Room{}
		if err := vdb.Scan(rows); err != nil {
			return err
		}
		v := Room{}
		if err := v.FromRow(vdb); err != nil {
			return err
		}
		res = append(res, v)
	}
	*v = res
	return nil
}