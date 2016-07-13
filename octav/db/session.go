package db

import (
	"database/sql"

	"github.com/pkg/errors"
)

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


func IsSessionOwner(tx *Tx, sessionID, userID string) error {
	stmt := getStmtBuf()
	defer releaseStmtBuf(stmt)

	stmt.WriteString(`SELECT 1 FROM `)
	stmt.WriteString(SessionTable)
	stmt.WriteString(` WHERE id = ? AND speaker_id = ?`)
	row := tx.QueryRow(stmt.String())
	var r int

	if err := row.Scan(&r); err != nil {
		return errors.Wrap(err, "failed to scan from database")
	}

	if r == 0 {
		return errors.Errorf("user %s is not an owner of session %s", userID, sessionID)
	}
	return nil
}

