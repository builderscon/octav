package db

import (
	"database/sql"

	"github.com/builderscon/octav/octav/tools"
	pdebug "github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

func init() {
	hooks = append(hooks, func() {
		stmt := tools.GetBuffer()
		defer tools.ReleaseBuffer(stmt)

		stmt.Reset()
		stmt.WriteString(`DELETE FROM `)
		stmt.WriteString(TrackTable)
		stmt.WriteString(` WHERE conference_id = ?`)
		library.Register("sqlDeleteTracksByConferenceID", stmt.String())

		stmt.Reset()
		stmt.WriteString(`SELECT `)
		stmt.WriteString(TrackStdSelectColumns)
		stmt.WriteString(` FROM `)
		stmt.WriteString(TrackTable)
		stmt.WriteString(` WHERE conference_id = ? AND room_id = ? LIMIT 1`)
		library.Register("sqlLoadTrack", stmt.String())

		stmt.Reset()
		stmt.WriteString(`SELECT `)
		stmt.WriteString(TrackStdSelectColumns)
		stmt.WriteString(` FROM `)
		stmt.WriteString(TrackTable)
		stmt.WriteString(` WHERE conference_id = ? ORDER BY sort_order ASC`)
		library.Register("sqlLoadByConferenceID", stmt.String())
	})
}

func DeleteTracks(tx *sql.Tx, conferenceID string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker(`DeleteTracks %s`, conferenceID).BindError(&err)
		defer g.End()
	}
	stmt, err := library.GetStmt("sqlDeleteTracksByConferenceID")
	if err != nil {
		return errors.Wrap(err, `failed to get statement`)
	}
	if _, err = tx.Stmt(stmt).Exec(conferenceID); err != nil {
		return errors.Wrap(err, `failed execute from database`)
	}

	return nil
}

func (vdb *Track) Load(tx *sql.Tx, conferenceID, roomID string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker(`Track.Load %s, %s`, conferenceID, roomID).BindError(&err)
		defer g.End()
	}

	stmt, err := library.GetStmt("sqlLoadTrack")
	if err != nil {
		return errors.Wrap(err, `failed to get statement`)
	}
	row := tx.Stmt(stmt).QueryRow(conferenceID, roomID)
	if err := vdb.Scan(row); err != nil {
		return errors.Wrap(err, `failed select from database`)
	}

	return nil
}

func (v *TrackList) LoadByConferenceID(tx *sql.Tx, conferenceID string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker(`TrackList.LoadByConferenceID %s`, conferenceID).BindError(&err)
		defer g.End()
	}

	stmt, err := library.GetStmt("sqlLoadByConferenceID")
	if err != nil {
		return errors.Wrap(err, `failed to get statement`)
	}
	rows, err := tx.Stmt(stmt).Query(conferenceID)
	if err := v.FromRows(rows, 0); err != nil {
		return errors.Wrap(err, `failed select from database`)
	}

	return nil
}
