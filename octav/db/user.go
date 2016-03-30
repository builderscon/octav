package db

func (vdb *User) LoadByAuthUserID(tx *Tx, via, id string) error {
	row := tx.QueryRow(`SELECT `+UserStdSelectColumns+` FROM `+UserTable+` WHERE users.auth_via = ? AND users.auth_user_id = ?`, via, id)
	if err := vdb.Scan(row); err != nil {
		return err
	}
	return nil
}
