package db

import (
	"database/sql"

	"github.com/builderscon/octav/octav/tools"
	pdebug "github.com/lestrrat/go-pdebug"
)

func (v *BlogEntryList) LoadByConference(tx *sql.Tx, confID string, status []string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("BlogEntryList.LoadByConference %s,%s", confID, status).BindError(&err)
		defer g.End()
	}

	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	stmt.WriteString(`SELECT `)
	stmt.WriteString(BlogEntryStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(BlogEntryTable)

	var args []interface{}
	args = append(args, confID)

	where := tools.GetBuffer()
	defer tools.ReleaseBuffer(where)

	where.WriteString(` WHERE `)
	where.WriteString(BlogEntryTable)
	where.WriteString(`.conference_id = ? `)
	if l := len(status); l > 0 {
		if where.Len() > 0 {
			where.WriteString(` AND `)
		}
		where.WriteString(BlogEntryTable)
		where.WriteString(`.status IN (`)
		for i, st := range status {
			where.WriteString(`?`)
			if i < l-1 {
				where.WriteString(`, `)
			}
			args = append(args, st)
		}
		where.WriteString(`)`)
	}

	if where.Len() > 0 {
		where.WriteString(` ORDER BY created_on ASC`)
	}

	where.WriteTo(stmt)

	rows, err := tx.Query(stmt.String(), args...)
	if err != nil {
		return err
	}

	if err := v.FromRows(rows, 0); err != nil {
		return err
	}
	return nil
}
