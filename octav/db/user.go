package db

import (
	"github.com/builderscon/octav/octav/tools"
	"github.com/pkg/errors"
)

var userLoadByAuthUserIDKey StmtKey
var userIsAdministratorKey StmtKey

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

	userIsAdministratorKey = makeStmtKey(stmt.Bytes())
	stmtPool.Register(userIsAdministratorKey, stmt.String())

	stmt.Reset()
	stmt.WriteString(`SELECT `)
	stmt.WriteString(UserStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(UserTable)
	stmt.WriteString(` WHERE users.auth_via = ? AND users.auth_user_id = ?`)

	userLoadByAuthUserIDKey = makeStmtKey(stmt.Bytes())
	stmtPool.Register(userLoadByAuthUserIDKey, stmt.String())
}

func (vdb *User) LoadByAuthUserID(tx *Tx, via, id string) error {
	stmt, err := stmtPool.Get(userLoadByAuthUserIDKey)
	if err != nil {
		return errors.Wrap(err, "failed to get statement")
	}

	row := tx.Stmt(stmt).QueryRow(via, id)
	if err := vdb.Scan(row); err != nil {
		return err
	}
	return nil
}

func IsAdministrator(tx *Tx, userID string) error {
	stmt, err := stmtPool.Get(userIsAdministratorKey)
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
