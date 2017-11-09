package db

import (
	"bytes"
	"database/sql"
	"strconv"

	"github.com/pkg/errors"
)

func (v *RoomList) LoadForVenueSinceEID(tx *sql.Tx, venueID string, since string, limit int) error {
	var s int64
	if id := since; id != "" {
		vdb := Room{}
		if err := vdb.LoadByEID(tx, id); err != nil {
			return err
		}

		s = vdb.OID
	}
	return v.LoadForVenueSince(tx, venueID, s, limit)
}

func (v *RoomList) LoadForVenueSince(tx *sql.Tx, venueID string, since int64, limit int) error {
	buf := bytes.Buffer{}
	buf.WriteString("SELECT ")
	buf.WriteString(RoomStdSelectColumns)
	buf.WriteString(" FROM ")
	buf.WriteString(RoomTable)
	buf.WriteString(" WHERE ")
	buf.WriteString(RoomTable)
	buf.WriteString(".oid > ? AND ")
	buf.WriteString(RoomTable)
	buf.WriteString(".venue_id = ? ORDER BY ")
	buf.WriteString(RoomTable)
	buf.WriteString(".oid ASC")
	if limit > 0 {
		buf.WriteString(" LIMIT ")
		buf.WriteString(strconv.Itoa(limit))
	}

	rows, err := Query(tx, buf.String(), since, venueID)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}
	defer rows.Close()

	res := make([]Room, 0, limit)
	for rows.Next() {
		vdb := Room{}
		if err := vdb.Scan(rows); err != nil {
			return errors.Wrap(err, `failed to scan row`)
		}
		res = append(res, vdb)
	}
	*v = res
	return nil
}

func LoadVenueRooms(tx *sql.Tx, rooms *RoomList, vid string) error {
	stmt := bytes.Buffer{}
	stmt.WriteString(`SELECT `)
	stmt.WriteString(RoomStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(RoomTable)
	stmt.WriteString(` WHERE `)
	stmt.WriteString(RoomTable)
	stmt.WriteString(`.venue_id = ?`)

	rows, err := Query(tx, stmt.String(), vid)
	if err != nil {
		return errors.Wrap(err, `failed to execute statement`)
	}
	defer rows.Close()

	var res RoomList
	for rows.Next() {
		var u Room
		if err := u.Scan(rows); err != nil {
			return errors.Wrap(err, `failed to scan row`)
		}

		res = append(res, u)
	}

	*rooms = res
	return nil
}
