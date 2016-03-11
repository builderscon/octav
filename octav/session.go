package octav

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/builderscon/octav/octav/db"
	"github.com/lestrrat/go-pdebug"
)

func (v Session) GetPropNames() ([]string, error) {
	l, _ := v.L10N.GetPropNames()
	return append(l, "id", "conference_id", "room_id", "speaker_id", "title", "abstract", "memo", "starts_on", "duration", "material_level", "tags", "category", "spoken_language", "slide_language", "slide_subtitles", "slide_url", "video_url", "photo_permission", "video_permission", "has_interpretation", "status", "sort_order", "confirmed", "conference", "room", "speaker"), nil
}

func (v Session) GetPropValue(s string) (interface{}, error) {
	switch s {
	case "id":
		return v.ID, nil
	case "conference_id":
		return v.ConferenceID, nil
	case "room_id":
		return v.RoomID, nil
	case "speaker_id":
		return v.SpeakerID, nil
	case "title":
		return v.Title, nil
	case "abstract":
		return v.Abstract, nil
	case "memo":
		return v.Memo, nil
	case "starts_on":
		return v.StartsOn, nil
	case "duration":
		return v.Duration, nil
	case "material_level":
		return v.MaterialLevel, nil
	case "tags":
		return v.Tags, nil
	case "category":
		return v.Category, nil
	case "spoken_language":
		return v.SpokenLanguage, nil
	case "slide_language":
		return v.SlideLanguage, nil
	case "slide_subtitles":
		return v.SlideSubtitles, nil
	case "slide_url":
		return v.SlideURL, nil
	case "video_url":
		return v.VideoURL, nil
	case "photo_permission":
		return v.PhotoPermission, nil
	case "video_permission":
		return v.VideoPermission, nil
	case "has_interpretation":
		return v.HasInterpretation, nil
	case "status":
		return v.Status, nil
	case "sort_order":
		return v.SortOrder, nil
	case "confirmed":
		return v.Confirmed, nil
	case "conference":
		return v.Conference, nil
	case "room":
		return v.Room, nil
	case "speaker":
		return v.Speaker, nil
	default:
		return v.L10N.GetPropValue(s)
	}
}

func (v Session) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["id"] = v.ID
	m["conference_id"] = v.ConferenceID
	if x := v.RoomID; x != "" {
		m["room_id"] = x
	}
	m["speaker_id"] = v.SpeakerID
	m["speaker"] = []interface{}{} // TODO
	m["title"] = v.Title

	if x := v.Abstract; x != "" {
		m["abstract"] = x
	}

	if x := v.Memo; x != "" {
		m["memo"] = x
	}

	if !v.StartsOn.IsZero() {
		m["starts_on"] = v.StartsOn.Format(time.RFC3339)
	}

	if x := v.Duration; x != 0 {
		m["duration"] = x
	}

	if x := v.MaterialLevel; x != "" {
		m["material_level"] = x
	}
	m["tags"] = v.Tags
	m["category"] = v.Category
	m["spoken_language"] = v.SpokenLanguage
	m["slide_language"] = v.SlideLanguage
	m["slide_subtitles"] = v.SlideSubtitles
	m["slide_url"] = v.SlideURL
	m["video_url"] = v.VideoURL
	m["photo_permission"] = v.PhotoPermission
	m["video_permission"] = v.VideoPermission
	m["has_interpretation"] = v.HasInterpretation

	if x := v.Status; x != "" {
		m["status"] = x
	}
	m["sort_order"] = v.SortOrder
	m["confirmed"] = v.Confirmed

	// m["conference"] = v.Conference
	// m["room"] = v.Room
	// m["speaker"] = v.Speaker
	buf, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return marshalJSONWithL10N(buf, v.L10N)
}

func (v *Session) UnmarshalJSON(data []byte) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("Session.UnmarshalJSON").BindError(&err)
		defer g.End()
	}

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
			return ErrInvalidJSONFieldType{Field: "id"}
		}
	}

	if jv, ok := m["conference_id"]; ok {
		switch jv.(type) {
		case string:
			v.ConferenceID = jv.(string)
			delete(m, "conference_id")
		default:
			return ErrInvalidJSONFieldType{Field: "conference_id"}
		}
	}

	if jv, ok := m["room_id"]; ok {
		switch jv.(type) {
		case string:
			v.RoomID = jv.(string)
			delete(m, "room_id")
		default:
			return ErrInvalidJSONFieldType{Field: "room_id"}
		}
	}

	if jv, ok := m["speaker_id"]; ok {
		switch jv.(type) {
		case string:
			v.SpeakerID = jv.(string)
			delete(m, "speaker_id")
		default:
			return ErrInvalidJSONFieldType{Field: "speaker_id"}
		}
	}

	if jv, ok := m["title"]; ok {
		switch jv.(type) {
		case string:
			v.Title = jv.(string)
			delete(m, "title")
		default:
			return ErrInvalidJSONFieldType{Field: "title"}
		}
	}

	if jv, ok := m["abstract"]; ok {
		switch jv.(type) {
		case string:
			v.Abstract = jv.(string)
			delete(m, "abstract")
		default:
			return ErrInvalidJSONFieldType{Field: "abstract"}
		}
	}

	if jv, ok := m["memo"]; ok {
		switch jv.(type) {
		case string:
			v.Memo = jv.(string)
			delete(m, "memo")
		default:
			return ErrInvalidJSONFieldType{Field: "memo"}
		}
	}

	if jv, ok := m["starts_on"]; ok {
		switch jv.(type) {
		case time.Time:
			v.StartsOn = jv.(time.Time)
			delete(m, "starts_on")
		case string:
			t, err := time.Parse(time.RFC3339, jv.(string))
			if err != nil {
				return err
			}
			v.StartsOn = t
			delete(m, "starts_on")
		default:
			return ErrInvalidJSONFieldType{Field: "starts_on"}
		}
	}

	if jv, ok := m["duration"]; ok {
		switch jv.(type) {
		case float64:
			v.Duration = int(jv.(float64))
			delete(m, "duration")
		default:
			pdebug.Printf("%s", reflect.ValueOf(jv).Type())
			return ErrInvalidJSONFieldType{Field: "duration"}
		}
	}

	if jv, ok := m["material_level"]; ok {
		switch jv.(type) {
		case string:
			v.MaterialLevel = jv.(string)
			delete(m, "material_level")
		default:
			return ErrInvalidJSONFieldType{Field: "material_level"}
		}
	}

	if jv, ok := m["tags"]; ok {
		switch jv.(type) {
		case []string:
			v.Tags = jv.([]string)
			delete(m, "tags")
		case nil:
			// no op
		default:
			return ErrInvalidJSONFieldType{Field: "tags"}
		}
	}

	if jv, ok := m["category"]; ok {
		switch jv.(type) {
		case string:
			v.Category = jv.(string)
			delete(m, "category")
		default:
			return ErrInvalidJSONFieldType{Field: "category"}
		}
	}

	if jv, ok := m["spoken_language"]; ok {
		switch jv.(type) {
		case string:
			v.SpokenLanguage = jv.(string)
			delete(m, "spoken_language")
		default:
			return ErrInvalidJSONFieldType{Field: "spoken_language"}
		}
	}

	if jv, ok := m["slide_language"]; ok {
		switch jv.(type) {
		case string:
			v.SlideLanguage = jv.(string)
			delete(m, "slide_language")
		default:
			return ErrInvalidJSONFieldType{Field: "slide_language"}
		}
	}

	if jv, ok := m["slide_subtitles"]; ok {
		switch jv.(type) {
		case string:
			v.SlideSubtitles = jv.(string)
			delete(m, "slide_subtitles")
		default:
			return ErrInvalidJSONFieldType{Field: "slide_subtitles"}
		}
	}

	if jv, ok := m["slide_url"]; ok {
		switch jv.(type) {
		case string:
			v.SlideURL = jv.(string)
			delete(m, "slide_url")
		default:
			return ErrInvalidJSONFieldType{Field: "slide_url"}
		}
	}

	if jv, ok := m["video_url"]; ok {
		switch jv.(type) {
		case string:
			v.VideoURL = jv.(string)
			delete(m, "video_url")
		default:
			return ErrInvalidJSONFieldType{Field: "video_url"}
		}
	}

	m["photo_permission"] = "allow"
	if jv, ok := m["photo_permission"]; ok {
		switch jv.(type) {
		case string:
			v.PhotoPermission = jv.(string)
			delete(m, "photo_permission")
		default:
			return ErrInvalidJSONFieldType{Field: "photo_permission"}
		}
	}

	m["video_permission"] = "allow"
	if jv, ok := m["video_permission"]; ok {
		switch jv.(type) {
		case string:
			v.VideoPermission = jv.(string)
			delete(m, "video_permission")
		default:
			return ErrInvalidJSONFieldType{Field: "video_permission"}
		}
	}

	if jv, ok := m["has_interpretation"]; ok {
		switch jv.(type) {
		case bool:
			v.HasInterpretation = jv.(bool)
			delete(m, "has_interpretation")
		default:
			return ErrInvalidJSONFieldType{Field: "has_interpretation"}
		}
	}

	if jv, ok := m["status"]; ok {
		switch jv.(type) {
		case string:
			v.Status = jv.(string)
			delete(m, "status")
		default:
			return ErrInvalidJSONFieldType{Field: "status"}
		}
	}

	if jv, ok := m["sort_order"]; ok {
		switch jv.(type) {
		case float64:
			v.SortOrder = int(jv.(float64))
			delete(m, "sort_order")
		default:
			return ErrInvalidJSONFieldType{Field: "sort_order"}
		}
	}

	if jv, ok := m["confirmed"]; ok {
		switch jv.(type) {
		case bool:
			v.Confirmed = jv.(bool)
			delete(m, "confirmed")
		default:
			return ErrInvalidJSONFieldType{Field: "confirmed"}
		}
	}

	return nil
}

func (v *Session) Load(tx *db.Tx, id string) error {
	vdb := db.Session{}
	if err := vdb.LoadByEID(tx, id); err != nil {
		return err
	}

	if err := v.FromRow(vdb); err != nil {
		return err
	}
	if err := v.LoadLocalizedFields(tx); err != nil {
		return err
	}
	return nil
}

func (v *Session) LoadLocalizedFields(tx *db.Tx) error {
	ls, err := db.LoadLocalizedStringsForParent(tx, v.ID, "Session")
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

func (v *Session) Create(tx *db.Tx) error {
	if v.ID == "" {
		v.ID = UUID()
	}

	vdb := db.Session{}
	if err := v.ToRow(&vdb); err != nil {
		return err
	}

	if err := vdb.Create(tx); err != nil {
		return err
	}

	if err := v.L10N.CreateLocalizedStrings(tx, "Session", v.ID); err != nil {
		return err
	}
	return nil
}

func (v *Session) Delete(tx *db.Tx) error {
	if pdebug.Enabled {
		g := pdebug.Marker("Session.Delete (%s)", v.ID)
		defer g.End()
	}

	vdb := db.Session{EID: v.ID}
	if err := vdb.Delete(tx); err != nil {
		return err
	}
	if err := db.DeleteLocalizedStringsForParent(tx, v.ID, "Session"); err != nil {
		return err
	}
	return nil
}

func (v *SessionList) Load(tx *db.Tx, since string, limit int) error {
	vdbl := db.SessionList{}
	if err := vdbl.LoadSinceEID(tx, since, limit); err != nil {
		return err
	}
	res := make([]Session, len(vdbl))
	for i, vdb := range vdbl {
		v := Session{}
		if err := v.FromRow(vdb); err != nil {
			return err
		}
		if err := v.LoadLocalizedFields(tx); err != nil {
			return err
		}
		res[i] = v
	}
	*v = res
	return nil
}

func (v *Session) FromRow(vdb db.Session) error {
	v.ID = vdb.EID
	v.ConferenceID = vdb.ConferenceID
	v.SpeakerID = vdb.SpeakerID
	if vdb.RoomID.Valid {
		v.RoomID = vdb.RoomID.String
	}
	if vdb.Title.Valid {
		v.Title = vdb.Title.String
	}
	if vdb.Abstract.Valid {
		v.Abstract = vdb.Abstract.String
	}
	if vdb.Abstract.Valid {
		v.Memo = vdb.Memo.String
	}
	if vdb.StartsOn.Valid {
		v.StartsOn = vdb.StartsOn.Time
	}
	v.Duration = vdb.Duration
	v.MaterialLevel = vdb.MaterialLevel
	if vdb.Tags.Valid {
		v.Tags = strings.Split(vdb.Tags.String, ",")
	}
	v.Category = vdb.Category
	v.SpokenLanguage = vdb.SpokenLanguage
	v.SlideLanguage = vdb.SlideLanguage
	v.SlideSubtitles = vdb.SlideSubtitles
	if vdb.SlideURL.Valid {
		v.SlideURL = vdb.SlideURL.String
	}
	if vdb.VideoURL.Valid {
		v.VideoURL = vdb.VideoURL.String
	}
	v.PhotoPermission = vdb.PhotoPermission
	v.VideoPermission = vdb.VideoPermission
	v.HasInterpretation = vdb.HasInterpretation
	v.Status = vdb.Status
	v.SortOrder = vdb.SortOrder
	v.Confirmed = vdb.Confirmed
	return nil
}

func (v *Session) ToRow(vdb *db.Session) error {
	vdb.EID = v.ID
	vdb.ConferenceID = v.ConferenceID
	vdb.SpeakerID = v.SpeakerID
	vdb.RoomID.Valid = true
	vdb.RoomID.String = v.RoomID
	vdb.Title.Valid = true
	vdb.Title.String = v.Title
	vdb.Abstract.Valid = true
	vdb.Abstract.String = v.Abstract
	vdb.Abstract.Valid = true
	vdb.Memo.String = v.Memo

	if !v.StartsOn.IsZero() {
		vdb.StartsOn.Valid = true
		vdb.StartsOn.Time = v.StartsOn
	}
	vdb.Duration = v.Duration
	vdb.MaterialLevel = v.MaterialLevel
	vdb.Tags.Valid = true
	vdb.Tags.String = strings.Join(v.Tags, ",")
	vdb.Category = v.Category
	vdb.SpokenLanguage = v.SpokenLanguage
	vdb.SlideLanguage = v.SlideLanguage
	vdb.SlideSubtitles = v.SlideSubtitles
	vdb.SlideURL.Valid = true
	vdb.SlideURL.String = v.SlideURL
	vdb.VideoURL.Valid = true
	vdb.VideoURL.String = v.VideoURL

	vdb.PhotoPermission = v.PhotoPermission
	vdb.VideoPermission = v.PhotoPermission
	vdb.HasInterpretation = v.HasInterpretation
	vdb.Status = v.Status
	vdb.SortOrder = v.SortOrder
	vdb.Confirmed = v.Confirmed
	return nil
}

func (v *SessionList) FromCursor(rows *sql.Rows) error {
	// Not using db.Session here
	res := make([]Session, 0, 10)
	for rows.Next() {
		vdb := db.Session{}
		if err := vdb.Scan(rows); err != nil {
			return err
		}
		v := Session{}
		if err := v.FromRow(vdb); err != nil {
			return err
		}
		res = append(res, v)
	}
	*v = res
	return nil
}

func (v *SessionList) LoadByConference(tx *db.Tx, cid, date string) error {
	var rows *sql.Rows
	var err error

	if date == "" {
		rows, err = tx.Query(`SELECT oid, eid, conference_id, room_id, speaker_id, title, abstract, memo, starts_on, duration, material_level, tags, category, spoken_language, slide_language, slide_subtitles, slide_url, video_url, photo_permission, video_permission, has_interpretation, status, sort_order, confirmed, created_on, modified_on FROM `+db.SessionTable+` WHERE conference_id = ?`, cid)
	} else {
		rows, err = tx.Query(`SELECT oid, eid, conference_id, room_id, speaker_id, title, abstract, memo, starts_on, duration, material_level, tags, category, spoken_language, slide_language, slide_subtitles, slide_url, video_url, photo_permission, video_permission, has_interpretation, status, sort_order, confirmed, created_on, modified_on FROM `+db.SessionTable+` WHERE conference_id = ? AND DATE(starts_on) = ?`, cid, date)
	}
	if err != nil {
		return err
	}

	if err := v.FromCursor(rows); err != nil {
		return err
	}
	return nil
}
