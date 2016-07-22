package db

import "strconv"

func (v *SponsorList) LoadByConferenceSinceEID(tx *Tx, confID, since string, limit int) error {
	var s int64
	if id := since; id != "" {
		vdb := Sponsor{}
		if err := vdb.LoadByEID(tx, id); err != nil {
			return err
		}

		s = vdb.OID
	}
	return v.LoadSince(tx, s, limit)
}

func (v *SponsorList) LoadByConferenceSince(tx *Tx, confID string, since int64, limit int) error {
	rows, err := tx.Query(`SELECT `+SponsorStdSelectColumns+` FROM `+SponsorTable+` WHERE conference_id = ? AND featured_speakers.oid > ? ORDER BY oid ASC LIMIT `+strconv.Itoa(limit), confID, since)
	if err != nil {
		return err
	}

	if err := v.FromRows(rows, limit); err != nil {
		return err
	}
	return nil
}

func LoadSponsors(tx *Tx, venues *SponsorList, cid string) error {
	stmt := getStmtBuf()
	defer releaseStmtBuf(stmt)

	stmt.WriteString(`SELECT `)
	stmt.WriteString(SponsorStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(SponsorTable)
	stmt.WriteString(` WHERE `)
	stmt.WriteString(SponsorTable)
	stmt.WriteString(`.conference_id = ?`)

	rows, err := tx.Query(stmt.String(), cid)
	if err != nil {
		return err
	}

	var res SponsorList
	for rows.Next() {
		var u Sponsor
		if err := u.Scan(rows); err != nil {
			return err
		}

		res = append(res, u)
	}

	*venues = res
	return nil
}
