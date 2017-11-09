package db

import (
	"database/sql"

	"github.com/builderscon/octav/octav/tools"
	pdebug "github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

var (
	sqlConferenceAdminLoad   string
	sqlConferenceAdminCheck  string
	sqlConferenceAdminDelete string
)

func init() {
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
	sqlConferenceAdminLoad = stmt.String()

	stmt.Reset()
	stmt.WriteString(`SELECT 1 FROM `)
	stmt.WriteString(ConferenceAdministratorTable)
	stmt.WriteString(` WHERE conference_id = ? AND user_id = ?`)
	sqlConferenceAdminCheck = stmt.String()

	stmt.Reset()
	stmt.WriteString(`DELETE FROM `)
	stmt.WriteString(ConferenceAdministratorTable)
	stmt.WriteString(` WHERE conference_id = ? AND user_id = ?`)
	sqlConferenceAdminDelete = stmt.String()
}

func IsConferenceAdministrator(tx *sql.Tx, cid, uid string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.IsConferenceAdministrator conference = %s, user = %s", cid, uid).BindError(&err)
		defer g.End()
	}

	row, err := QueryRow(tx, sqlConferenceAdminCheck, cid, uid)
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

func DeleteConferenceAdministrator(tx *sql.Tx, cid, uid string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.DeleteConferenceAdministrator conference = %s, user = %s", cid, uid).BindError(&err)
		defer g.End()
	}
	if _, err := Exec(tx, sqlConferenceAdminDelete, cid, uid); err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}
	return nil
}

func LoadConferenceAdministrators(tx *sql.Tx, admins *UserList, cid string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.LoadConferenceAdministrators conference = %s", cid).BindError(&err)
		defer g.End()
	}
	rows, err := Query(tx, sqlConferenceAdminLoad, cid)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}
	defer rows.Close()

	var res UserList
	for rows.Next() {
		var u User
		if err := u.Scan(rows); err != nil {
			return errors.Wrap(err, `failed to scan row`)
		}

		res = append(res, u)
	}

	*admins = res
	return nil
}
