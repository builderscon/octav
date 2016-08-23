package db

import "github.com/pkg/errors"

func DeleteConferenceComponentsByIDAndName(tx *Tx, conferenceID string, names ...string) error {
	stmt := getStmtBuf()
	defer releaseStmtBuf(stmt)

	l := len(names)
	if l == 0 {
		return errors.New("empty list of names")
	}

	stmt.WriteString(`DELETE FROM `)
	stmt.WriteString(ConferenceComponentTable)
	stmt.WriteString(` WHERE conference_id = ? AND name IN (`)
	for i := range names {
		stmt.WriteByte('?')
		if i < l - 1 {
			stmt.WriteByte(',')
		}
	}
	stmt.WriteString(`)`)

	args := make([]interface{}, len(names)+1)
	args[0] = conferenceID
	for i, name := range names {
		args[i+1] = name
	}
	if _, err := tx.Exec(stmt.String(), args...); err != nil {
		return errors.Wrap(err, "failed to execute delete statement")
	}

	return nil
}

func UpsertConferenceComponentsByIDAndName(tx *Tx, conferenceID string, values map[string]string) error {
	stmt := getStmtBuf()
	defer releaseStmtBuf(stmt)

	l := len(values)
	if l == 0 {
		return errors.New("empty value map")
	}

	stmt.WriteString(`INSERT INTO `)
	stmt.WriteString(ConferenceComponentTable)
	stmt.WriteString(` (conference_id, name, value) VALUES `)

	var args []interface{}
	i := 0
	for k, v := range values {
		args = append(args, conferenceID, k, v)
		stmt.WriteString(`(?, ?, ?)`)
		if i < l - 1 {
			stmt.WriteByte(',')
		}
		i++
	}

	stmt.WriteString(` ON DUPLICATE UPDATE value = VALUES(value)`)
	if _, err := tx.Exec(stmt.String(), args...); err != nil {
		return errors.Wrap(err, "failed to execute insert statement")
	}

	return nil
}


func (ccl *ConferenceComponentList) LoadByConferenceID(tx *Tx, cid string) error {
	stmt := getStmtBuf()
	defer releaseStmtBuf(stmt)

	stmt.WriteString(`SELECT `)
	stmt.WriteString(ConferenceComponentStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(ConferenceComponentTable)
	stmt.WriteString(` WHERE conference_id = ?`)

	rows, err := tx.Query(stmt.String(), cid)
	if err != nil {
		return errors.Wrap(err, "failed to execute query for conference component")
	}

	var res ConferenceComponentList
	for rows.Next() {
		var row ConferenceComponent
		if err := row.Scan(rows); err != nil {
			return errors.Wrap(err, "failed to scan conference component")
		}
		res = append(res, row)
	}

	*ccl = res
	return nil
}
