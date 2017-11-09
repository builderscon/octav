package db

import (
	"database/sql"
	"strconv"

	"github.com/builderscon/octav/octav/tools"
	"github.com/pkg/errors"
)

var (
	sqlSessionTypeIsAcceptingSubmissions string
	sqlSessionTypeLoadByConferenceID     string
)

func init() {
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
	sqlSessionTypeIsAcceptingSubmissions = stmt.String()

	stmt.Reset()
	stmt.WriteString(`SELECT `)
	stmt.WriteString(SessionTypeStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(SessionTypeTable)
	stmt.WriteString(` WHERE `)
	stmt.WriteString(SessionTypeTable)
	stmt.WriteString(`.conference_id = ?`)
	stmt.WriteString(` ORDER BY sort_order ASC, oid ASC`)
	sqlSessionTypeLoadByConferenceID = stmt.String()

}

func IsAcceptingSubmissions(tx *sql.Tx, id string) error {
	var i int
	row, err := QueryRow(tx, sqlSessionTypeIsAcceptingSubmissions)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}

	if err := row.Scan(&i); err != nil {
		return errors.Wrap(err, "failed to select for IsAcceptingSubmissions")
	}

	if i == 1 {
		return nil
	}
	return errors.New("currently not accepting submissions")
}

func (v *SessionTypeList) LoadByConferenceSinceEID(tx *sql.Tx, confID, since string, limit int) error {
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

func (v *SessionTypeList) LoadByConferenceSince(tx *sql.Tx, confID string, since int64, limit int) error {
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

	rows, err := Query(tx, stmt.String(), confID, since)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}
	defer rows.Close()

	if err := v.FromRows(rows, limit); err != nil {
		return err
	}
	return nil
}

func (v *SessionTypeList) LoadByConferenceID(tx *sql.Tx, cid string) error {
	rows, err := Query(tx, sqlSessionTypeLoadByConferenceID, cid)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}
	defer rows.Close()

	var res SessionTypeList
	for rows.Next() {
		var u SessionType
		if err := u.Scan(rows); err != nil {
			return errors.Wrap(err, `failed to scan row`)
		}

		res = append(res, u)
	}

	*v = res
	return nil
}
