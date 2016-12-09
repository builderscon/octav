package db

import (
	"github.com/builderscon/octav/octav/tools"
	"github.com/pkg/errors"
)

const (
	sqlConferenceStaffLoadKey   = "sqlConferenceStaffLoad"
	sqlConferenceStaffDeleteKey = "sqlConferenceStaffDelete"
)

func init() {
	hooks = append(hooks, func() {
		stmt := tools.GetBuffer()
		defer tools.ReleaseBuffer(stmt)

		stmt.Reset()
		stmt.WriteString(`SELECT `)
		stmt.WriteString(UserStdSelectColumns)
		stmt.WriteString(` FROM `)
		stmt.WriteString(ConferenceStaffTable)
		stmt.WriteString(` JOIN `)
		stmt.WriteString(UserTable)
		stmt.WriteString(` ON `)
		stmt.WriteString(ConferenceStaffTable)
		stmt.WriteString(`.user_id = `)
		stmt.WriteString(UserTable)
		stmt.WriteString(`.eid WHERE `)
		stmt.WriteString(ConferenceStaffTable)
		stmt.WriteString(`.conference_id = ? ORDER BY sort_order ASC`)
		library.Register(sqlConferenceStaffLoadKey, stmt.String())

		stmt.Reset()
		stmt.WriteString(`DELETE FROM `)
		stmt.WriteString(ConferenceStaffTable)
		stmt.WriteString(` WHERE conference_id = ? AND user_id = ?`)
		library.Register(sqlConferenceStaffDeleteKey, stmt.String())
	})
}

func DeleteConferenceStaff(tx *Tx, cid, uid string) error {
	stmt, err := library.GetStmt(sqlConferenceStaffDeleteKey)
	if err != nil {
		return errors.Wrap(err, `failed to get statement`)
	}
	_, err = tx.Stmt(stmt).Exec(cid, uid)
	return errors.Wrap(err, `failed to execute statements`)
}

func LoadConferenceStaff(tx *Tx, admins *UserList, cid string) error {
	stmt, err := library.GetStmt(sqlConferenceStaffLoadKey)
	if err != nil {
		return errors.Wrap(err, `failed to get statement`)
	}
	rows, err := tx.Stmt(stmt).Query(cid)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}

	var res UserList
	for rows.Next() {
		var u User
		if err := u.Scan(rows); err != nil {
			return errors.Wrap(err, `failed to scan rows`)
		}

		res = append(res, u)
	}

	*admins = res
	return nil
}
