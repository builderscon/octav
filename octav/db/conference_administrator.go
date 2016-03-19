package db

import "bytes"

func DeleteConferenceAdministrator(tx *Tx, cid, uid string) error {
	stmt := bytes.Buffer{}
	stmt.WriteString(`DELETE FROM `)
	stmt.WriteString(ConferenceAdministratorTable)
	stmt.WriteString(` WHERE conference_id = ? AND user_id = ?`)

	_, err := tx.Exec(stmt.String(), cid, uid)
	return err
}

func LoadConferenceAdministrators(tx *Tx, admins *UserList, cid string) error {
	stmt := bytes.Buffer{}
	stmt.WriteString(`SELECT `)
	stmt.WriteString(UserStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(ConferenceAdministratorTable)
	stmt.WriteString(` JOIN `)
	stmt.WriteString(UserTable)
	stmt.WriteString(` ON `)
	stmt.WriteString(ConferenceAdministratorTable)
	stmt.WriteString(`.user_id = `)
	stmt.WriteString(UserTable)
	stmt.WriteString(`.eid WHERE `)
	stmt.WriteString(ConferenceAdministratorTable)
	stmt.WriteString(`.conference_id = ?`)

	rows, err := tx.Query(stmt.String(), cid)
	if err != nil {
		return err
	}

	var res UserList
	for rows.Next() {
		var u User
		if err := u.Scan(rows); err != nil {
			return err
		}

		res = append(res, u)
	}

	*admins = res
	return nil
}