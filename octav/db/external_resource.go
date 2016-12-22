package db

import (
	"github.com/builderscon/octav/octav/tools"
	"github.com/lestrrat/go-pdebug"
)

func (v *ExternalResourceList) LoadByConference(tx *Tx, cid string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.ExternalResourceList.LoadByConference %s", cid).BindError(&err)
		defer g.End()
	}

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
	if err := v.FromRows(rows, 0); err != nil {
		return err
	}

	return nil
}
