package db

import "github.com/pkg/errors"

func (vdb *TemporaryEmail) LoadByUserIDAndConfirmationKey(tx *Tx, userID, confirmationKey string) error {
	stmt := getStmtBuf()
	defer releaseStmtBuf(stmt)

	stmt.WriteString(`SELECT `)
	stmt.WriteString(TemporaryEmailStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(TemporaryEmailTable)
	stmt.WriteString(` WHERE user_id = ? AND confirmation_key = ?`)

	row := tx.QueryRow(stmt.String(), userID, confirmationKey)
	if err := vdb.Scan(row); err != nil {
		return errors.Wrap(err, "failed to execute query")
	}

	return nil
}
