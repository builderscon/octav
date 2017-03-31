package db

import (
	"database/sql"

	"github.com/builderscon/octav/octav/tools"
	"github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

func init() {
	hooks = append(hooks, func() {
		stmt := tools.GetBuffer()
		defer tools.ReleaseBuffer(stmt)

		stmt.Reset()
		stmt.WriteString(`SELECT `)
		stmt.WriteString(ExternalResourceStdSelectColumns)
		stmt.WriteString(` FROM `)
		stmt.WriteString(ExternalResourceTable)
		stmt.WriteString(` WHERE conference_id = ? ORDER BY sort_order ASC`)
		library.Register("sqlExternalResourceLoadByConferenceID", stmt.String())
	})
}

func (v *ExternalResourceList) LoadByConference(tx *sql.Tx, conferenceID string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker(`ExternalResourceList.LoadByConference %s`, conferenceID).BindError(&err)
		defer g.End()
	}

	stmt, err := library.GetStmt("sqlExternalResourceLoadByConferenceID")
	if err != nil {
		return errors.Wrap(err, `failed to get statement`)
	}
	rows, err := tx.Stmt(stmt).Query(conferenceID)
	if err := v.FromRows(rows, 0); err != nil {
		return errors.Wrap(err, `failed select from database`)
	}

	return nil
}
