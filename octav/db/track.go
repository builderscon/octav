package db

import (
	"database/sql"

	"github.com/builderscon/octav/octav/tools"
	pdebug "github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

var (
	sqlTrackDeleteTracksByConferenceID string
	sqlTrackLoad                       string
	sqlTrackListLoadByConferenceID     string
)

func init() {
	stmt := tools.GetBuffer()
	defer tools.ReleaseBuffer(stmt)

	stmt.Reset()
	stmt.WriteString(`DELETE FROM `)
	stmt.WriteString(TrackTable)
	stmt.WriteString(` WHERE conference_id = ?`)
	sqlTrackDeleteTracksByConferenceID = stmt.String()

	stmt.Reset()
	stmt.WriteString(`SELECT `)
	stmt.WriteString(TrackStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(TrackTable)
	stmt.WriteString(` WHERE conference_id = ? AND room_id = ? LIMIT 1`)
	sqlTrackLoad = stmt.String()

	stmt.Reset()
	stmt.WriteString(`SELECT `)
	stmt.WriteString(TrackStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(TrackTable)
	stmt.WriteString(` WHERE conference_id = ? ORDER BY sort_order ASC`)
	sqlTrackListLoadByConferenceID = stmt.String()
}

func DeleteTracks(tx *sql.Tx, conferenceID string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker(`DeleteTracks %s`, conferenceID).BindError(&err)
		defer g.End()
	}

	if _, err = Exec(tx, sqlTrackDeleteTracksByConferenceID, conferenceID); err != nil {
		return errors.Wrap(err, `failed execute statement`)
	}

	return nil
}

func (vdb *Track) Load(tx *sql.Tx, conferenceID, roomID string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker(`Track.Load %s, %s`, conferenceID, roomID).BindError(&err)
		defer g.End()
	}

	row, err := QueryRow(tx, sqlTrackLoad, conferenceID, roomID)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}

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

	rows, err := Query(tx, sqlTrackListLoadByConferenceID, conferenceID)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}
	defer rows.Close()

	if err := v.FromRows(rows, 0); err != nil {
		return errors.Wrap(err, `failed select from database`)
	}

	return nil
}
