package db

import (
	"database/sql"

	"github.com/builderscon/octav/octav/tools"
	"github.com/pkg/errors"
)

var (
	sqlTemporaryEmailLoadByUserIDAndConfirmationKey string
	sqlTemporaryEmailUpsert                         string
)

func init() {
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	stmt.WriteString(`SELECT `)
	stmt.WriteString(TemporaryEmailStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(TemporaryEmailTable)
	stmt.WriteString(` WHERE user_id = ? AND confirmation_key = ?`)
	sqlTemporaryEmailLoadByUserIDAndConfirmationKey = stmt.String()

	stmt.Reset()
	stmt.WriteString(`INSERT INTO `)
	stmt.WriteString(TemporaryEmailTable)
	stmt.WriteString(` (user_id, confirmation_key, email, expires_on) VALUES (?, ?, ?, ?) `)
	stmt.WriteString(` ON DUPLICATE KEY UPDATE confirmation_key = VALUES(confirmation_key), expires_on = VALUES(expires_on)`)
	sqlTemporaryEmailUpsert = stmt.String()
}

func (vdb *TemporaryEmail) LoadByUserIDAndConfirmationKey(tx *sql.Tx, userID, confirmationKey string) error {
	row, err := QueryRow(tx, sqlTemporaryEmailLoadByUserIDAndConfirmationKey, userID, confirmationKey)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}

	if err := vdb.Scan(row); err != nil {
		return errors.Wrap(err, "failed to execute query")
	}

	return nil
}

func (vdb *TemporaryEmail) Upsert(tx *sql.Tx) error {
	result, err := Exec(tx, sqlTemporaryEmailUpsert, vdb.UserID, vdb.ConfirmationKey, vdb.Email, vdb.ExpiresOn)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}

	lii, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, `failed to fetch last insert ID`)
	}

	vdb.OID = lii
	return nil
}
