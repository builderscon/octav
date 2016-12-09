package db

import (
	"github.com/builderscon/octav/octav/tools"
	"github.com/pkg/errors"
)

const (
	sqlConferenceAdminLoadKey   = "sqlConferenceAdminLoad"
	sqlConferenceAdminCheckKey  = "sqlConferenceAdminCheck"
	sqlConferenceAdminDeleteKey = "sqlConferenceAdminDelete"
)

func init() {
	hooks = append(hooks, func() {
		stmt := tools.GetBuffer()
		defer tools.ReleaseBuffer(stmt)

		stmt.Reset()
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
		stmt.WriteString(`.conference_id = ? ORDER BY sort_order ASC`)
		library.Register(sqlConferenceAdminLoadKey, stmt.String())

		stmt.Reset()
		stmt.WriteString(`SELECT 1 FROM `)
		stmt.WriteString(ConferenceAdministratorTable)
		stmt.WriteString(` WHERE conference_id = ? AND user_id = ?`)
		library.Register(sqlConferenceAdminCheckKey, stmt.String())

		stmt.Reset()
		stmt.WriteString(`DELETE FROM `)
		stmt.WriteString(ConferenceAdministratorTable)
		stmt.WriteString(` WHERE conference_id = ? AND user_id = ?`)
		library.Register(sqlConferenceAdminDeleteKey, stmt.String())
	})
}

func IsConferenceAdministrator(tx *Tx, cid, uid string) error {
	stmt, err := library.GetStmt(sqlConferenceAdminCheckKey)
	if err != nil {
		return errors.Wrap(err, `failed to get statement`)
	}
	row := tx.Stmt(stmt).QueryRow(cid, uid)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}

	var v int
	if err := row.Scan(&v); err != nil {
		return errors.Wrap(err, "failed to scan row")
	}
	if v != 1 {
		return errors.New("no matching administrator found")
	}
	return nil
}

func DeleteConferenceAdministrator(tx *Tx, cid, uid string) error {
	stmt, err := library.GetStmt(sqlConferenceAdminDeleteKey)
	if err != nil {
		return errors.Wrap(err, `failed to get statement`)
	}
	_, err = tx.Stmt(stmt).Exec(cid, uid)
	return errors.Wrap(err, `failed to execute statements`)
}

func LoadConferenceAdministrators(tx *Tx, admins *UserList, cid string) error {
	stmt, err := library.GetStmt(sqlConferenceAdminLoadKey)
	if err != nil {
		return errors.Wrap(err, `failed to get statement`)
	}
	rows, err := tx.Stmt(stmt).Query(cid)
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
