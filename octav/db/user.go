package db

import (
	"database/sql"

	"github.com/builderscon/octav/octav/tools"
	"github.com/pkg/errors"
)

var (
	sqlUserIsAdministrator   string
	sqlUserListLoadFromQuery string
	sqlUserLoadByAuthUserID  string
)

func init() {
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	stmt.WriteString(`SELECT 1 FROM `)
	stmt.WriteString(UserTable)
	stmt.WriteString(` WHERE `)
	stmt.WriteString(UserTable)
	stmt.WriteString(`.is_admin = 1 AND `)
	stmt.WriteString(UserTable)
	stmt.WriteString(`.eid = ? UNION SELECT 1 FROM `)
	stmt.WriteString(ConferenceAdministratorTable)
	stmt.WriteString(` WHERE `)
	stmt.WriteString(ConferenceAdministratorTable)
	stmt.WriteString(`.user_id = ? UNION SELECT 1 FROM `)
	stmt.WriteString(ConferenceSeriesAdministratorTable)
	stmt.WriteString(` WHERE `)
	stmt.WriteString(ConferenceSeriesAdministratorTable)
	stmt.WriteString(`.user_id = ?`)
	sqlUserIsAdministrator = stmt.String()

	stmt.Reset()
	stmt.WriteString(`SELECT `)
	stmt.WriteString(UserStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(UserTable)
	stmt.WriteString(` WHERE users.auth_via = ? AND users.auth_user_id = ?`)
	sqlUserLoadByAuthUserID = stmt.String()

	stmt.Reset()
	stmt.WriteString(`SELECT `)
	stmt.WriteString(UserStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(UserTable)
	stmt.WriteString(` WHERE nickname LIKE ? AND oid > ? ORDER BY nickname ASC LIMIT ?`)
	sqlUserListLoadFromQuery = stmt.String()
}

func (vdb *User) LoadByAuthUserID(tx *sql.Tx, via, id string) error {
	row, err := QueryRow(tx, sqlUserLoadByAuthUserID, via, id)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}

	if err := vdb.Scan(row); err != nil {
		return errors.Wrap(err, `failed to scan row`)
	}
	return nil
}

func IsAdministrator(tx *sql.Tx, userID string) error {
	var v int
	row, err := QueryRow(tx, sqlUserIsAdministrator, userID, userID, userID)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}

	if err := row.Scan(&v); err != nil {
		return errors.Wrap(err, "failed to scan row")
	}

	if v == 0 {
		return errors.Errorf("user %s is not an administrator", userID)
	}
	return nil
}

func (vdbl *UserList) LoadFromQuery(tx *sql.Tx, pattern, since string, limit int) error {
	if pattern == "" {
		return vdbl.LoadSinceEID(tx, since, limit)
	}

	var s int64
	if id := since; id != "" {
		vdb := User{}
		if err := vdb.LoadByEID(tx, id); err != nil {
			return err
		}

		s = vdb.OID
	}

	patbuf := tools.GetBuffer()
	defer tools.ReleaseBuffer(patbuf)

	for _, r := range pattern {
		if r == '%' {
			patbuf.WriteByte('\\')
		}
		patbuf.WriteRune(r)
	}
	patbuf.WriteByte('%')

	rows, err := Query(tx, sqlUserListLoadFromQuery, patbuf.String(), s, limit)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}
	defer rows.Close()

	if err := vdbl.FromRows(rows, limit); err != nil {
		return errors.Wrap(err, "failed to scan results")
	}
	return nil
}
