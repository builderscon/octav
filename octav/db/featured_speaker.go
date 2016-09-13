package db

import (
	"strconv"

	"github.com/builderscon/octav/octav/tools"
	sqllib "github.com/lestrrat/go-sqllib"
	"github.com/pkg/errors"
)

var sqlFeaturedSpeakerLoadFeaturedSpeakersKey sqllib.Key

func init() {
	hooks = append(hooks, func() {
		buf := tools.GetBuffer()
		defer tools.ReleaseBuffer(buf)

		buf.WriteString(`SELECT `)
		buf.WriteString(FeaturedSpeakerStdSelectColumns)
		buf.WriteString(` FROM `)
		buf.WriteString(FeaturedSpeakerTable)
		buf.WriteString(` WHERE `)
		buf.WriteString(FeaturedSpeakerTable)
		buf.WriteString(`.conference_id = ?`)

		sqlFeaturedSpeakerLoadFeaturedSpeakersKey = library.Register(buf.String())
	})
}

func (v *FeaturedSpeakerList) LoadByConferenceSinceEID(tx *Tx, confID, since string, limit int) error {
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
	stmt, err := library.GetStmt(sqlFeaturedSpeakerLoadFeaturedSpeakersKey)
	if err != nil {
		return errors.Wrap(err, "failed to get statement")
	}

	rows, err := tx.Stmt(stmt).Query(cid)
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
