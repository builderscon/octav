package db

import "github.com/lestrrat/go-pdebug"

func (l *LocalizedString) LoadByLangKey(tx *Tx, language, name, parentType, parentID string) error {
	row := tx.QueryRow(`SELECT oid, parent_id, parent_type, name, language, localized FROM `+LocalizedStringTable+` WHERE parent_type = ? AND parent_id = ? AND name = ? AND language = ?`, parentType, parentID, name, language)
	if err := l.Scan(row); err != nil {
		return err
	}
	return nil
}

func (l *LocalizedString) Upsert(tx *Tx) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("LocalizedString.Upsert (%s#%s)", l.Language, l.Name).BindError(&err)
		defer g.End()
	}

	result, err := tx.Exec(`INSERT INTO `+LocalizedStringTable+` (parent_id, parent_type, name, language, localized) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE localized = VALUES(localized)`, l.ParentID, l.ParentType, l.Name, l.Language, l.Localized)
	if err != nil {
		return err
	}

	lii, err := result.LastInsertId()
	if err != nil {
		return err
	}

	l.OID = lii
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
