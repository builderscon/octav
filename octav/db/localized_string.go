package db

func (l *LocalizedString) LoadByLangKey(tx *Tx, name, language string) error {
	row := tx.QueryRow(`SELECT oid, parent_id, parent_type, name, language, localized FROM `+LocalizedStringTable+` WHERE name = ? AND language = ?`, name, language)
	if err := l.Scan(row); err != nil {
		return err
	}
	return nil
}

func LoadLocalizedStringsForParent(tx *Tx, parentID, parentType string) (LocalizedStringList, error) {
	ret := []LocalizedString{}

	rows, err := tx.Query(`SELECT oid, parent_id, parent_type, name, language, localized FROM `+LocalizedStringTable+` WHERE parent_id = ? AND parent_type = ?`, parentID, parentType)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		l := LocalizedString{}
		if err := l.Scan(rows); err != nil {
			return nil, err
		}
		ret = append(ret, l)
	}
	return ret, nil
}

func DeleteLocalizedStringsForParent(tx *Tx, parentID, parentType string) error {
	_, err := tx.Exec(`DELETE FROM `+LocalizedStringTable+` WHERE parent_id = ? AND parent_type = ?`, parentID, parentType)
	if err != nil {
		return err
	}
	return nil
}

