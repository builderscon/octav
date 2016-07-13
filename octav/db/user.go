package db

import (
	"fmt"

	"github.com/pkg/errors"
)

func (vdb *User) LoadByAuthUserID(tx *Tx, via, id string) error {
	row := tx.QueryRow(`SELECT `+UserStdSelectColumns+` FROM `+UserTable+` WHERE users.auth_via = ? AND users.auth_user_id = ?`, via, id)
	if err := vdb.Scan(row); err != nil {
		return err
	}
	return nil
}

func IsAdministrator(tx *Tx, userID string) error {
	stmt := getStmtBuf()
	defer releaseStmtBuf(stmt)
	fmt.Fprintf(stmt, `SELECT 1 FROM %s WHERE %s.is_admin = 1 and %s.eid = ?`, UserTable, UserTable, UserTable)
	fmt.Fprintf(stmt, ` UNION SELECT 1 FROM %s WHERE %s.user_id = ?`, ConferenceAdministratorTable, ConferenceAdministratorTable)
	fmt.Fprintf(stmt, ` UNION SELECT 1 FROM %s WHERE %s.user_id = ?`, ConferenceSeriesAdministratorTable, ConferenceSeriesAdministratorTable)

	var v int
	row := tx.QueryRow(stmt.String(), userID, userID, userID)
	if err := row.Scan(&v); err != nil {
		return errors.Wrap(err, "failed to scan row")
	}

	if v == 0 {
		return errors.Errorf("user %s is not an administrator", userID)
	}
	return nil
}
