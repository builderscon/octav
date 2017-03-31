package db

import (
	"database/sql"

	"github.com/builderscon/octav/octav/tools"
	"github.com/pkg/errors"
)

func init() {
	hooks = append(hooks, func() {
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

		library.Register("userIsAdministratorKey", stmt.String())

		stmt.Reset()
		stmt.WriteString(`SELECT `)
		stmt.WriteString(UserStdSelectColumns)
		stmt.WriteString(` FROM `)
		stmt.WriteString(UserTable)
		stmt.WriteString(` WHERE users.auth_via = ? AND users.auth_user_id = ?`)

		library.Register("userLoadByAuthUserIDKey", stmt.String())

		stmt.Reset()
		stmt.WriteString(`SELECT `)
		stmt.WriteString(UserStdSelectColumns)
		stmt.WriteString(` FROM `)
		stmt.WriteString(UserTable)
		stmt.WriteString(` WHERE nickname LIKE ? AND oid > ? ORDER BY nickname ASC LIMIT ?`)
		library.Register("userListLoadFromQuery", stmt.String())
	})
}

func (vdb *User) LoadByAuthUserID(tx *sql.Tx, via, id string) error {
	stmt, err := library.GetStmt("userLoadByAuthUserIDKey")
	if err != nil {
		return errors.Wrap(err, "failed to get statement")
	}

	row := tx.Stmt(stmt).QueryRow(via, id)
	if err := vdb.Scan(row); err != nil {
		return err
	}
	return nil
}

func IsAdministrator(tx *sql.Tx, userID string) error {
	stmt, err := library.GetStmt("userIsAdministratorKey")
	if err != nil {
		return errors.Wrap(err, "failed to get statement")
	}

	var v int
	row := tx.Stmt(stmt).QueryRow(userID, userID, userID)
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

	stmt, err := library.GetStmt("userListLoadFromQuery")
	if err != nil {
		return errors.Wrap(err, "failed to get statement")
	}

	rows, err := tx.Stmt(stmt).Query(patbuf.String(), s, limit)
	if err := vdbl.FromRows(rows, limit); err != nil {
		return errors.Wrap(err, "failed to scan results")
	}
	return nil
}
