package db

import (
	"github.com/builderscon/octav/octav/tools"
	"github.com/pkg/errors"
)

func IsConferenceAdministrator(tx *Tx, cid, uid string) error {
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)
	stmt.WriteString(`SELECT 1 FROM `)
	stmt.WriteString(ConferenceAdministratorTable)
	stmt.WriteString(` WHERE conference_id = ? AND user_id = ?`)

	var v int
	row := tx.QueryRow(stmt.String(), cid, uid)
	if err := row.Scan(&v); err != nil {
		return errors.Wrap(err, "failed to scan row")
	}
	if v != 1 {
		return errors.New("no matching administrator found")
	}
	return nil
}

func DeleteConferenceAdministrator(tx *Tx, cid, uid string) error {
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)
	stmt.WriteString(`DELETE FROM `)
	stmt.WriteString(ConferenceAdministratorTable)
	stmt.WriteString(` WHERE conference_id = ? AND user_id = ?`)

	_, err := tx.Exec(stmt.String(), cid, uid)
	return err
}

func LoadConferenceAdministrators(tx *Tx, admins *UserList, cid string) error {
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)
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
