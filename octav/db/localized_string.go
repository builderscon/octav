package db

import (
	"database/sql"

	"github.com/builderscon/octav/octav/tools"
	"github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

var (
	sqlLocalizedStringLoadByLangKey                   string
	sqlLocalizedStringUpsert                          string
	sqlLocalizedStringLoadLocalizedStringsForParent   string
	sqlLocalizedStringDeleteLocalizedStringsForParent string
)

func init() {
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	stmt.WriteString(`SELECT oid, parent_id, parent_type, name, language, localized FROM `)
	stmt.WriteString(LocalizedStringTable)
	stmt.WriteString(` WHERE parent_type = ? AND parent_id = ? AND name = ? AND language = ?`)
	sqlLocalizedStringLoadByLangKey = stmt.String()

	stmt.Reset()
	stmt.WriteString(`INSERT INTO `)
	stmt.WriteString(LocalizedStringTable)
	stmt.WriteString(`(parent_id, parent_type, name, language, localized) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE localized = VALUES(localized)`)
	sqlLocalizedStringUpsert = stmt.String()

	stmt.Reset()
	stmt.WriteString(`SELECT oid, parent_id, parent_type, name, language, localized FROM `)
	stmt.WriteString(LocalizedStringTable)
	stmt.WriteString(` WHERE parent_id = ? AND parent_type = ?`)
	sqlLocalizedStringLoadLocalizedStringsForParent = stmt.String()

	stmt.Reset()
	stmt.WriteString(`DELETE FROM `)
	stmt.WriteString(LocalizedStringTable)
	stmt.WriteString(` WHERE parent_id = ? AND parent_type = ?`)
	sqlLocalizedStringDeleteLocalizedStringsForParent = stmt.String()
}

func (l *LocalizedString) LoadByLangKey(tx *sql.Tx, language, name, parentType, parentID string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("LocalizedString.LoadByLangKey %s %s %s %s", language, name, parentType, parentID).BindError(&err)
		defer g.End()
	}
	if err != nil {
		return errors.Wrap(err, "failed to get statement")
	}

	row, err := QueryRow(tx, sqlLocalizedStringLoadByLangKey, parentType, parentID, name, language)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}

	if err := l.Scan(row); err != nil {
		return errors.Wrap(err, `failed to scan row`)
	}
	return nil
}

func (l *LocalizedString) Upsert(tx *sql.Tx) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("LocalizedString.Upsert (%s#%s)", l.Language, l.Name).BindError(&err)
		defer g.End()
	}

	result, err := Exec(tx, sqlLocalizedStringUpsert, l.ParentID, l.ParentType, l.Name, l.Language, l.Localized)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}

	lii, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, `failed to get last insert ID`)
	}

	l.OID = lii
	return nil
}

func (v *LocalizedStringList) LoadForParent(tx *sql.Tx, parentID, parentType string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.LocalizedStringList.LoadForParent %s %s", parentID, parentType).BindError(&err)
		defer g.End()
	}

	rows, err := Query(tx, sqlLocalizedStringLoadLocalizedStringsForParent, parentID, parentType)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}
	defer rows.Close()

	var ret []LocalizedString
	for rows.Next() {
		var l LocalizedString
		if err := l.Scan(rows); err != nil {
			return err
		}
		ret = append(ret, l)
	}
	*v = ret
	return nil
}

func DeleteLocalizedStringsForParent(tx *sql.Tx, parentID, parentType string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.DeleteLocalizedStringsForParent %s %s", parentID, parentType).BindError(&err)
		defer g.End()
	}
	if _, err = Exec(tx, sqlLocalizedStringDeleteLocalizedStringsForParent, parentID, parentType); err != nil {
		return errors.Wrap(err, "failed to execute query")
	}
	return nil
}
