package db

import (
	"github.com/builderscon/octav/octav/tools"
	"github.com/pkg/errors"
)

func (vdb *TemporaryEmail) LoadByUserIDAndConfirmationKey(tx *Tx, userID, confirmationKey string) error {
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

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

func (vdb *TemporaryEmail) Upsert(tx *Tx) error {
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	stmt.WriteString(`INSERT INTO `)
	stmt.WriteString(TemporaryEmailTable)
	stmt.WriteString(` (user_id, confirmation_key, email, expires_on) VALUES (?, ?, ?, ?) `)
	stmt.WriteString(` ON DUPLICATE KEY UPDATE confirmation_key = VALUES(confirmation_key), expires_on = VALUES(expires_on)`)

	result, err := tx.Exec(stmt.String(), vdb.UserID, vdb.ConfirmationKey, vdb.Email, vdb.ExpiresOn)
	if err != nil {
		return err
	}

	lii, err := result.LastInsertId()
	if err != nil {
		return err
	}

	vdb.OID = lii
	return nil
}
