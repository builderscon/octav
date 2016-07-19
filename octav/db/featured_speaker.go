package db

import "strconv"

func (v *FeaturedSpeakerList) LoadByConferenceSinceEID(tx *Tx, confID, since string, limit int) error {
	var s int64
	if id := since; id != "" {
		vdb := FeaturedSpeaker{}
		if err := vdb.LoadByEID(tx, id); err != nil {
			return err
		}

		s = vdb.OID
	}
	return v.LoadSince(tx, s, limit)
}

func (v *FeaturedSpeakerList) LoadByConferenceSince(tx *Tx, confID string, since int64, limit int) error {
	rows, err := tx.Query(`SELECT `+FeaturedSpeakerStdSelectColumns+` FROM `+FeaturedSpeakerTable+` WHERE conference_id = ? AND featured_speakers.oid > ? ORDER BY oid ASC LIMIT `+strconv.Itoa(limit), confID, since)
	if err != nil {
		return err
	}

	if err := v.FromRows(rows, limit); err != nil {
		return err
	}
	return nil
}
