package db

func (v *Venue) LoadByEID(tx *Tx, eid string) error {
	row := tx.QueryRow(`SELECT oid, eid, name, created_on, modified_on WHERE eid = ?`, eid)
	if err := row.Scan(&v.OID, &v.EID, &v.Name, &v.CreatedOn, &v.ModifiedOn); err != nil {
		return err
	}
	return nil
}

