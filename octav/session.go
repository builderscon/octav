package octav

import (
	"database/sql"
	"strings"

	"github.com/builderscon/octav/octav/db"
)

func (v *Session) Load(tx *db.Tx, id string) error {
	vdb := db.Session{}
	if err := vdb.LoadByEID(tx, id); err != nil {
		return err
	}

	return v.FromRow(vdb)
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

func (v *SessionList) Load(tx *db.Tx, since string) error {
	var s int64
	if id := since; id != "" {
		vdb := db.Session{}
		if err := vdb.LoadByEID(tx, id); err != nil {
			return err
		}

		s = vdb.OID
	}

	rows, err := tx.Query(`SELECT eid, name FROM venue WHERE oid > ? ORDER BY oid LIMIT 10`, s)
	if err != nil {
		return err
	}

	if err := v.FromCursor(rows); err != nil {
		return err
	}
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
