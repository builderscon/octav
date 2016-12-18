package db

import (
	"github.com/builderscon/octav/octav/tools"
)

func LoadExternalResources(tx *Tx, externalResources *ExternalResourceList, cid string) error {
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	stmt.WriteString(`SELECT `)
	stmt.WriteString(ExternalResourceStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(ExternalResourceTable)
	stmt.WriteString(` WHERE `)
	stmt.WriteString(ExternalResourceTable)
	stmt.WriteString(`.conference_id = ?`)

	rows, err := tx.Query(stmt.String(), cid)
	if err != nil {
		return err
	}

	var res ExternalResourceList
	for rows.Next() {
		var u ExternalResource
		if err := u.Scan(rows); err != nil {
			return err
		}

		res = append(res, u)
	}

	*externalResources = res
	return nil
}
