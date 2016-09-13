package db

import (
	"github.com/builderscon/octav/octav/tools"
	"github.com/lestrrat/go-pdebug"
	sqllib "github.com/lestrrat/go-sqllib"
	"github.com/pkg/errors"
)

var sqlLocalizedStringLoadByLangKeyKey sqllib.Key
var sqlLocalizedStringUpsertKey sqllib.Key
var sqlLocalizedStringLoadLocalizedStringsForParentKey sqllib.Key
var sqlLocalizedStringDeleteLocalizedStringsForParentKey sqllib.Key

func init() {
	hooks = append(hooks, func() {
		buf := tools.GetBuffer()
		defer tools.ReleaseBuffer(buf)

		buf.WriteString(`SELECT oid, parent_id, parent_type, name, language, localized FROM `)
		buf.WriteString(LocalizedStringTable)
		buf.WriteString(` WHERE parent_type = ? AND parent_id = ? AND name = ? AND language = ?`)

		sqlLocalizedStringLoadByLangKeyKey = library.Register(buf.String())

		buf.Reset()
		buf.WriteString(`INSERT INTO `)
		buf.WriteString(LocalizedStringTable)
		buf.WriteString(`(parent_id, parent_type, name, language, localized) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE localized = VALUES(localized)`)
		sqlLocalizedStringUpsertKey = library.Register(buf.String())

		buf.Reset()
		buf.WriteString(`SELECT oid, parent_id, parent_type, name, language, localized FROM `)
		buf.WriteString(LocalizedStringTable)
		buf.WriteString(` WHERE parent_id = ? AND parent_type = ?`)
		sqlLocalizedStringLoadLocalizedStringsForParentKey = library.Register(buf.String())

		buf.Reset()
		buf.WriteString(`DELETE FROM `)
		buf.WriteString(LocalizedStringTable)
		buf.WriteString(` WHERE parent_id = ? AND parent_type = ?`)
		sqlLocalizedStringDeleteLocalizedStringsForParentKey = library.Register(buf.String())
	})
}

func (l *LocalizedString) LoadByLangKey(tx *Tx, language, name, parentType, parentID string) error {
	stmt, err := library.GetStmt(sqlLocalizedStringLoadByLangKeyKey)
	if err != nil {
		return errors.Wrap(err, "failed to get statement")
	}

	row := tx.Stmt(stmt).QueryRow(parentType, parentID, name, language)
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

	stmt, err := library.GetStmt(sqlLocalizedStringUpsertKey)
	if err != nil {
		return errors.Wrap(err, "failed to get statement")
	}

	result, err := tx.Stmt(stmt).Exec(l.ParentID, l.ParentType, l.Name, l.Language, l.Localized)
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
	stmt, err := library.GetStmt(sqlLocalizedStringLoadLocalizedStringsForParentKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get statement")
	}

	var ret []LocalizedString

	rows, err := tx.Stmt(stmt).Query(parentID, parentType)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var l LocalizedString
		if err := l.Scan(rows); err != nil {
			return nil, err
		}
		ret = append(ret, l)
	}
	return ret, nil
}

func DeleteLocalizedStringsForParent(tx *Tx, parentID, parentType string) error {
	stmt, err := library.GetStmt(sqlLocalizedStringDeleteLocalizedStringsForParentKey)
	if err != nil {
		return errors.Wrap(err, "failed to get statement")
	}

	if _, err = tx.Stmt(stmt).Exec(parentID, parentType); err != nil {
		return errors.Wrap(err, "failed to execute query")
	}
	return nil
}
