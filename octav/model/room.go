package model

import "github.com/builderscon/octav/octav/db"

func (v *Room) unmarshalFromMap(m map[string]interface{}) error {
	if jv, ok := m["id"]; ok {
		switch jv.(type) {
		case string:
			v.ID = jv.(string)
		default:
			return ErrInvalidJSONFieldType{Field: "id"}
		}
	}

	if jv, ok := m["venue_id"]; ok {
		switch jv.(type) {
		case string:
			v.VenueID = jv.(string)
		default:
			return ErrInvalidJSONFieldType{Field: "venue_id"}
		}
	}

	if jv, ok := m["name"]; ok {
		switch jv.(type) {
		case string:
			v.Name = jv.(string)
		default:
			return ErrInvalidJSONFieldType{Field: "name"}
		}
	}

	if jv, ok := m["capacity"]; ok {
		switch jv.(type) {
		case float64:
			v.Capacity = uint(jv.(float64))
		default:
			return ErrInvalidJSONFieldType{Field: "capacity"}
		}
	}

	return nil
}

func (v *RoomList) LoadForVenue(tx *db.Tx, venueID, since string, limit int) error {
	vdbl := db.RoomList{}
	if err := vdbl.LoadForVenueSinceEID(tx, venueID, since, limit); err != nil {
		return err
	}

	res := make([]RoomL10N, len(vdbl))
	for i, vdb := range vdbl {
		v := Room{}
		if err := v.FromRow(vdb); err != nil {
			return err
		}
		vl := RoomL10N{Room: v}
		if err := vl.LoadLocalizedFields(tx); err != nil {
			return err
		}
		res[i] = vl
	}
	*v = res
	return nil
}
