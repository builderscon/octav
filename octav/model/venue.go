package model

func (v *Venue) unmarshalFromMap(m map[string]interface{}) error {
	if jv, ok := m["id"]; ok {
		switch jv.(type) {
		case string:
			v.ID = jv.(string)
		default:
			return ErrInvalidJSONFieldType{Field: "id"}
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

	if jv, ok := m["address"]; ok {
		switch jv.(type) {
		case string:
			v.Address = jv.(string)
		default:
			return ErrInvalidJSONFieldType{Field: "address"}
		}
	}

	if jv, ok := m["latitude"]; ok {
		switch jv.(type) {
		case float64:
			v.Latitude = jv.(float64)
		default:
			return ErrInvalidJSONFieldType{Field: "latitude"}
		}
	}

	if jv, ok := m["longitude"]; ok {
		switch jv.(type) {
		case float64:
			v.Longitude = jv.(float64)
		default:
			return ErrInvalidJSONFieldType{Field: "longitude"}
		}
	}

	return nil
}


