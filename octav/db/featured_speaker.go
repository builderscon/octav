package db

import (
	"database/sql"
	"strconv"

	"github.com/builderscon/octav/octav/tools"
	"github.com/pkg/errors"
)

var (
	sqlFeaturedSpeakerListLoadByConferenceID string
)

func init() {
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	stmt.WriteString(`SELECT `)
	stmt.WriteString(FeaturedSpeakerStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(FeaturedSpeakerTable)
	stmt.WriteString(` WHERE `)
	stmt.WriteString(FeaturedSpeakerTable)
	stmt.WriteString(`.conference_id = ?`)

	sqlFeaturedSpeakerListLoadByConferenceID = stmt.String()
}

func (v *FeaturedSpeakerList) LoadByConferenceSinceEID(tx *sql.Tx, confID, since string, limit int) error {
	var s int64
	if id := since; id != "" {
		var vdb FeaturedSpeaker
		if err := vdb.LoadByEID(tx, id); err != nil {
			return err
		}

		s = vdb.OID
	}
	return v.LoadSince(tx, s, limit)
}

func (v *FeaturedSpeakerList) LoadByConferenceSince(tx *sql.Tx, confID string, since int64, limit int) error {
	rows, err := Query(tx, `SELECT `+FeaturedSpeakerStdSelectColumns+` FROM `+FeaturedSpeakerTable+` WHERE conference_id = ? AND featured_speakers.oid > ? ORDER BY oid ASC LIMIT `+strconv.Itoa(limit), confID, since)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}
	defer rows.Close()

	if err := v.FromRows(rows, limit); err != nil {
		return err
	}
	return nil
}

func (v *FeaturedSpeakerList) LoadByConferenceID(tx *sql.Tx, cid string) error {
	rows, err := Query(tx, sqlFeaturedSpeakerListLoadByConferenceID, cid)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}
	defer rows.Close()

	var res FeaturedSpeakerList
	for rows.Next() {
		var u FeaturedSpeaker
		if err := u.Scan(rows); err != nil {
			return errors.Wrap(err, `failed to scan row`)
		}

		res = append(res, u)
	}

	*v = res
	return nil
}
