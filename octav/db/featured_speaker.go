package db

import (
	"bytes"
	"strconv"
)

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

func LoadFeaturedSpeakers(tx *Tx, venues *FeaturedSpeakerList, cid string) error {
	stmt := bytes.Buffer{}
	stmt.WriteString(`SELECT `)
	stmt.WriteString(FeaturedSpeakerStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(FeaturedSpeakerTable)
	stmt.WriteString(` WHERE `)
	stmt.WriteString(FeaturedSpeakerTable)
	stmt.WriteString(`.conference_id = ?`)

	rows, err := tx.Query(stmt.String(), cid)
	if err != nil {
		return err
	}

	var res FeaturedSpeakerList
	for rows.Next() {
		var u FeaturedSpeaker
		if err := u.Scan(rows); err != nil {
			return err
		}

		res = append(res, u)
	}

	*venues = res
	return nil
}
