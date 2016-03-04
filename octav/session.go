package octav

import (
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