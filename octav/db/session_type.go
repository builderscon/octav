package db

import (
	"strconv"

	"github.com/builderscon/octav/octav/tools"
	pdebug "github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

func IsAcceptingSubmissions(tx *Tx, id string) error {
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	stmt.WriteString(`SELECT 1 FROM `)
	stmt.WriteString(ConferenceTable)
	stmt.WriteString(` JOIN `)
	stmt.WriteString(SessionTypeTable)
	stmt.WriteString(` ON `)
	stmt.WriteString(ConferenceTable)
	stmt.WriteString(`.eid = `)
	stmt.WriteString(SessionTypeTable)
	stmt.WriteString(`.conference_id WHERE `)
	stmt.WriteString(ConferenceTable)
	stmt.WriteString(`.status != "private" AND `)
	stmt.WriteString(SessionTypeTable)
	stmt.WriteString(`.submission_start <= NOW() AND `)
	stmt.WriteString(SessionTypeTable)
	stmt.WriteString(`.submission_end >= NOW()`)

	var i int
	row := tx.QueryRow(stmt.String())
	if err := row.Scan(&i); err != nil {
		return errors.Wrap(err, "failed to select for IsAcceptingSubmissions")
	}

	if i == 1 {
		return nil
	}
	return errors.New("currently not accepting submissions")
}

func (v *SessionTypeList) LoadByConferenceSinceEID(tx *Tx, confID, since string, limit int) error {
	var s int64
	if id := since; id != "" {
		vdb := SessionType{}
		if err := vdb.LoadByEID(tx, id); err != nil {
			return err
		}

		s = vdb.OID
	}
	return v.LoadByConferenceSince(tx, confID, s, limit)
}

func (v *SessionTypeList) LoadByConferenceSince(tx *Tx, confID string, since int64, limit int) error {
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	stmt.WriteString(`SELECT `)
	stmt.WriteString(SessionTypeStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(SessionTypeTable)
	stmt.WriteString(` WHERE `)
	stmt.WriteString(SessionTypeTable)
	stmt.WriteString(`.conference_id = ? AND `)
	stmt.WriteString(SessionTypeTable)
	stmt.WriteString(`.oid > ? ORDER BY oid ASC LIMIT `)
	stmt.WriteString(strconv.Itoa(limit))

	pdebug.Printf(stmt.String())

	rows, err := tx.Query(stmt.String(), confID, since)
	if err != nil {
		return err
	}

	if err := v.FromRows(rows, limit); err != nil {
		return err
	}
	return nil
}

func LoadSessionTypes(tx *Tx, list *SessionTypeList, cid string) error {
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	stmt.WriteString(`SELECT `)
	stmt.WriteString(SessionTypeStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(SessionTypeTable)
	stmt.WriteString(` WHERE `)
	stmt.WriteString(SessionTypeTable)
	stmt.WriteString(`.conference_id = ?`)
	stmt.WriteString(` ORDER BY sort_order ASC, oid ASC`)

	rows, err := tx.Query(stmt.String(), cid)
	if err != nil {
		return err
	}

	var res SessionTypeList
	for rows.Next() {
		var u SessionType
		if err := u.Scan(rows); err != nil {
			return err
		}

		res = append(res, u)
	}

	*list = res
	return nil
}
