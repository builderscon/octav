package model

func (v *Conference) unmarshalFromMap(m map[string]interface{}) error {
	if mv, ok := m["id"]; ok {
		switch mv.(type) {
		case string:
		default:
			return ErrInvalidJSONFieldType{Field: "id"}
		}
		v.ID = mv.(string)
	}

	if mv, ok := m["title"]; ok {
		switch mv.(type) {
		case string:
		default:
			return ErrInvalidJSONFieldType{Field: "title"}
		}
		v.Title = mv.(string)
	}

	if mv, ok := m["sub_title"]; ok {
		switch mv.(type) {
		case string:
		default:
			return ErrInvalidJSONFieldType{Field: "sub_title"}
		}
		v.SubTitle = mv.(string)
	}

	if mv, ok := m["slug"]; ok {
		switch mv.(type) {
		case string:
		default:
			return ErrInvalidJSONFieldType{Field: "slug"}
		}
		v.Slug = mv.(string)
	}

	return nil
}


