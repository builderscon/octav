package db

import (
	"database/sql"

	"github.com/builderscon/octav/octav/tools"
	"github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

var (
	sqlExternalResourceLoadByConference string
)

func init() {
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	stmt.Reset()
	stmt.WriteString(`SELECT `)
	stmt.WriteString(ExternalResourceStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(ExternalResourceTable)
	stmt.WriteString(` WHERE conference_id = ? ORDER BY sort_order ASC`)
	sqlExternalResourceLoadByConference = stmt.String()
}

func (v *ExternalResourceList) LoadByConference(tx *sql.Tx, conferenceID string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker(`ExternalResourceList.LoadByConference %s`, conferenceID).BindError(&err)
		defer g.End()
	}

	rows, err := Query(tx, sqlExternalResourceLoadByConference, conferenceID)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}
	defer rows.Close()

	if err := v.FromRows(rows, 0); err != nil {
		return errors.Wrap(err, `failed select from database`)
	}

	return nil
}
