package db

import "database/sql"

func (v *SessionList) LoadByConference(tx *Tx, cid, date string) error {
	var rows *sql.Rows
	var err error

	if date == "" {
		rows, err = tx.Query(`SELECT oid, eid, conference_id, room_id, speaker_id, title, abstract, memo, starts_on, duration, material_level, tags, category, spoken_language, slide_language, slide_subtitles, slide_url, video_url, photo_permission, video_permission, has_interpretation, status, sort_order, confirmed, created_on, modified_on FROM `+SessionTable+` WHERE conference_id = ?`, cid)
	} else {
		rows, err = tx.Query(`SELECT oid, eid, conference_id, room_id, speaker_id, title, abstract, memo, starts_on, duration, material_level, tags, category, spoken_language, slide_language, slide_subtitles, slide_url, video_url, photo_permission, video_permission, has_interpretation, status, sort_order, confirmed, created_on, modified_on FROM `+SessionTable+` WHERE conference_id = ? AND DATE(starts_on) = ?`, cid, date)
	}
	if err != nil {
		return err
	}

	if err := v.FromRows(rows, 0); err != nil {
		return err
	}
	return nil
}

