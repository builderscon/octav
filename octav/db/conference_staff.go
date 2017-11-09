package db

import (
	"database/sql"

	"github.com/builderscon/octav/octav/tools"
	pdebug "github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

var (
	sqlConferenceStaffLoad   string
	sqlConferenceStaffDelete string
)

func init() {
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
	sqlConferenceStaffLoad = stmt.String()

	stmt.Reset()
	stmt.WriteString(`DELETE FROM `)
	stmt.WriteString(ConferenceStaffTable)
	stmt.WriteString(` WHERE conference_id = ? AND user_id = ?`)
	sqlConferenceStaffDelete = stmt.String()
}

func DeleteConferenceStaff(tx *sql.Tx, cid, uid string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.DeleteConferenceStaff conference %s, user %s", cid, uid).BindError(&err)
		defer g.End()
	}

	_, err = Exec(tx, sqlConferenceStaffDelete, cid, uid)
	return errors.Wrap(err, `failed to execute statements`)
}

func LoadConferenceStaff(tx *sql.Tx, admins *UserList, cid string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.LoadConferenceStaff conference %s", cid).BindError(&err)
		defer g.End()
	}
	rows, err := Query(tx, sqlConferenceStaffLoad, cid)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}
	defer rows.Close()

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
