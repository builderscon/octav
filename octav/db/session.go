package db

import "github.com/pkg/errors"

func (v *SessionList) LoadByConference(tx *Tx, conferenceID, speakerID, date, status string) error {
	// The caller of this method should ensure that query fields are
	// present and that we don't accidentally run an empty query
	stmt := getStmtBuf()
	defer releaseStmtBuf(stmt)

	stmt.WriteString(`SELECT `)
	stmt.WriteString(SessionStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(SessionTable)
	stmt.WriteString(` WHERE `)

	where := getStmtBuf()
	defer releaseStmtBuf(where)

	var args []interface{}
	if conferenceID != "" {
		where.WriteString(` conference_id = ?`)
		args = append(args, conferenceID)
	}

	if speakerID != "" {
		if where.Len() > 0 {
			where.WriteString(` AND`)
		}
		where.WriteString(` speaker_id = ?`)
		args = append(args, speakerID)
	}

	if date != "" {
		if where.Len() > 0 {
			where.WriteString(` AND`)
		}
		where.WriteString(` DATE(date) = ?`)
		args = append(args, date)
	}

	where.WriteTo(stmt)

	rows, err := tx.Query(stmt.String(), args...)
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
