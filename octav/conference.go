package octav

import (
	"database/sql"
	"encoding/json"

	"github.com/builderscon/octav/octav/db"
	"github.com/lestrrat/go-pdebug"
)

func (v Conference) GetPropNames() ([]string, error) {
	l, _ := v.L10N.GetPropNames()
	return append(l, "id", "title", "subtitle", "slug", "dates"), nil
}

func (v Conference) GetPropValue(s string) (interface{}, error) {
	switch s {
	case "id":
		return v.ID, nil
	case "title":
		return v.Title, nil
	case "subtitle":
		return v.SubTitle, nil
	case "slug":
		return v.Slug, nil
	case "dates":
		return v.Dates, nil
	default:
		return v.L10N.GetPropValue(s)
	}
}

func (v Conference) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["id"] = v.ID
	m["title"] = v.Title
	m["subtitle"] = v.SubTitle
	m["slug"] = v.Slug
	buf, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return marshalJSONWithL10N(buf, v.L10N)
}
func (v *Conference) UnmarshalJSON(data []byte) error {
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
			return ErrInvalidJSONFieldType{Field:"id"}
		}
	}
	if jv, ok := m["title"]; ok {
		switch jv.(type) {
		case string:
			v.Title = jv.(string)
			delete(m, "title")
		default:
			return ErrInvalidJSONFieldType{Field:"title"}
		}
	}
	if jv, ok := m["subtitle"]; ok {
		switch jv.(type) {
		case string:
			v.SubTitle = jv.(string)
			delete(m, "subtitle")
		default:
			return ErrInvalidJSONFieldType{Field:"subtitle"}
		}
	}
	if jv, ok := m["slug"]; ok {
		switch jv.(type) {
		case string:
			v.Slug = jv.(string)
			delete(m, "slug")
		default:
			return ErrInvalidJSONFieldType{Field:"slug"}
		}
	}
	return nil
}

func (v *Conference) Load(tx *db.Tx, id string) error {
	vdb := db.Conference{}
	if err := vdb.LoadByEID(tx, id); err != nil {
		return err
	}

	v.ID = vdb.EID
	v.Title = vdb.Title
	if vdb.SubTitle.Valid {
		v.SubTitle = vdb.SubTitle.String
	}
	v.Slug = vdb.Slug

	ls, err := db.LoadLocalizedStringsForParent(tx, v.ID, "Conference")
	if err != nil {
		return err
	}

	if len(ls) > 0 {
		v.L10N = LocalizedFields{}
		for _, l := range ls {
			v.L10N.Set(l.Language, l.Name, l.Localized)
		}
	}
	return nil
}

func (v *Conference) Create(tx *db.Tx) error {
	if v.ID == "" {
		v.ID = UUID()
	}

	vdb := db.Conference{
		EID:      v.ID,
		Title:    v.Title,
		SubTitle: sql.NullString{String: v.SubTitle, Valid: true},
		Slug:     v.Slug,
	}
	if err := vdb.Create(tx); err != nil {
		return err
	}

	if err := v.L10N.CreateLocalizedStrings(tx, "Conference", v.ID); err != nil {
		return err
	}
	return nil
}
func (v *Conference) FromRow(vdb db.Conference) error {
	v.ID = vdb.EID
	v.Slug = vdb.Slug
	v.Title = vdb.Title
	if vdb.SubTitle.Valid {
		v.SubTitle = vdb.SubTitle.String
	}
	return nil
}

func (v Conference) ToRow(vdb *db.Conference) error {
	vdb.EID = v.ID
	vdb.Slug = v.Slug
	vdb.Title = v.Title
	vdb.SubTitle.Valid = true
	vdb.SubTitle.String = v.SubTitle
	return nil
}

func (v *Conference) Delete(tx *db.Tx) error {
	if pdebug.Enabled {
		g := pdebug.Marker("Conference.Delete (%s)", v.ID)
		defer g.End()
	}

	vdb := db.Conference{EID: v.ID}
	if err := vdb.Delete(tx); err != nil {
		return err
	}
	if err := db.DeleteLocalizedStringsForParent(tx, v.ID, "Conference"); err != nil {
		return err
	}
	return nil
}

func (v *ConferenceList) Load(tx *db.Tx, since string) error {
	var s int64
	if id := since; id != "" {
		vdb := db.Conference{}
		if err := vdb.LoadByEID(tx, id); err != nil {
			return err
		}

		s = vdb.OID
	}

	rows, err := tx.Query(`SELECT eid, name FROM venue WHERE oid > ? ORDER BY oid LIMIT 10`, s)
	if err != nil {
		return err
	}

	// Not using db.Conference here
	res := make([]Conference, 0, 10)
	for rows.Next() {
		vdb := db.Conference{}
		if err := vdb.Scan(rows); err != nil {
			return err
		}
		v := Conference{}
		if err := v.FromRow(vdb); err != nil {
			return err
		}
		res = append(res, v)
	}
	*v = res
	return nil
}
