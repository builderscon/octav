package db

// Automatically generated by gendb utility. DO NOT EDIT!

import (
	"bytes"
	"database/sql"
	"strconv"
	"time"

	"github.com/builderscon/octav/octav/tools"
	"github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

const FeaturedSpeakerStdSelectColumns = "featured_speakers.oid, featured_speakers.eid, featured_speakers.conference_id, featured_speakers.speaker_id, featured_speakers.avatar_url, featured_speakers.display_name, featured_speakers.description, featured_speakers.created_on, featured_speakers.modified_on"
const FeaturedSpeakerTable = "featured_speakers"

type FeaturedSpeakerList []FeaturedSpeaker

func (f *FeaturedSpeaker) Scan(scanner interface {
	Scan(...interface{}) error
}) error {
	return scanner.Scan(&f.OID, &f.EID, &f.ConferenceID, &f.SpeakerID, &f.AvatarURL, &f.DisplayName, &f.Description, &f.CreatedOn, &f.ModifiedOn)
}

func init() {
	hooks = append(hooks, func() {
		stmt := tools.GetBuffer()
		defer tools.ReleaseBuffer(stmt)

		stmt.Reset()
		stmt.WriteString(`DELETE FROM `)
		stmt.WriteString(FeaturedSpeakerTable)
		stmt.WriteString(` WHERE oid = ?`)
		library.Register("sqlFeaturedSpeakerDeleteByOIDKey", stmt.String())

		stmt.Reset()
		stmt.WriteString(`UPDATE `)
		stmt.WriteString(FeaturedSpeakerTable)
		stmt.WriteString(` SET eid = ?, conference_id = ?, speaker_id = ?, avatar_url = ?, display_name = ?, description = ? WHERE oid = ?`)
		library.Register("sqlFeaturedSpeakerUpdateByOIDKey", stmt.String())

		stmt.Reset()
		stmt.WriteString(`SELECT `)
		stmt.WriteString(FeaturedSpeakerStdSelectColumns)
		stmt.WriteString(` FROM `)
		stmt.WriteString(FeaturedSpeakerTable)
		stmt.WriteString(` WHERE `)
		stmt.WriteString(FeaturedSpeakerTable)
		stmt.WriteString(`.eid = ?`)
		library.Register("sqlFeaturedSpeakerLoadByEIDKey", stmt.String())

		stmt.Reset()
		stmt.WriteString(`DELETE FROM `)
		stmt.WriteString(FeaturedSpeakerTable)
		stmt.WriteString(` WHERE eid = ?`)
		library.Register("sqlFeaturedSpeakerDeleteByEIDKey", stmt.String())

		stmt.Reset()
		stmt.WriteString(`UPDATE `)
		stmt.WriteString(FeaturedSpeakerTable)
		stmt.WriteString(` SET eid = ?, conference_id = ?, speaker_id = ?, avatar_url = ?, display_name = ?, description = ? WHERE eid = ?`)
		library.Register("sqlFeaturedSpeakerUpdateByEIDKey", stmt.String())
	})
}

func (f *FeaturedSpeaker) LoadByEID(tx *Tx, eid string) error {
	stmt, err := library.GetStmt("sqlFeaturedSpeakerLoadByEIDKey")
	if err != nil {
		return errors.Wrap(err, `failed to get statement`)
	}
	row := tx.Stmt(stmt).QueryRow(eid)
	if err := f.Scan(row); err != nil {
		return err
	}
	return nil
}

func (f *FeaturedSpeaker) Create(tx *Tx, opts ...InsertOption) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.FeaturedSpeaker.Create").BindError(&err)
		defer g.End()
		pdebug.Printf("%#v", f)
	}
	if f.EID == "" {
		return errors.New("create: non-empty EID required")
	}

	f.CreatedOn = time.Now()
	doIgnore := false
	for _, opt := range opts {
		switch opt.(type) {
		case insertIgnoreOption:
			doIgnore = true
		}
	}

	stmt := bytes.Buffer{}
	stmt.WriteString("INSERT ")
	if doIgnore {
		stmt.WriteString("IGNORE ")
	}
	stmt.WriteString("INTO ")
	stmt.WriteString(FeaturedSpeakerTable)
	stmt.WriteString(` (eid, conference_id, speaker_id, avatar_url, display_name, description, created_on, modified_on) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	result, err := tx.Exec(stmt.String(), f.EID, f.ConferenceID, f.SpeakerID, f.AvatarURL, f.DisplayName, f.Description, f.CreatedOn, f.ModifiedOn)
	if err != nil {
		return err
	}

	lii, err := result.LastInsertId()
	if err != nil {
		return err
	}

	f.OID = lii
	return nil
}

func (f FeaturedSpeaker) Update(tx *Tx) error {
	if f.OID != 0 {
		stmt, err := library.GetStmt("sqlFeaturedSpeakerUpdateByOIDKey")
		if err != nil {
			return errors.Wrap(err, `failed to get statement`)
		}
		_, err = tx.Stmt(stmt).Exec(f.EID, f.ConferenceID, f.SpeakerID, f.AvatarURL, f.DisplayName, f.Description, f.OID)
		return err
	}
	if f.EID != "" {
		stmt, err := library.GetStmt("sqlFeaturedSpeakerUpdateByEIDKey")
		if err != nil {
			return errors.Wrap(err, `failed to get statement`)
		}
		_, err = tx.Stmt(stmt).Exec(f.EID, f.ConferenceID, f.SpeakerID, f.AvatarURL, f.DisplayName, f.Description, f.EID)
		return err
	}
	return errors.New("either OID/EID must be filled")
}

func (f FeaturedSpeaker) Delete(tx *Tx) error {
	if f.OID != 0 {
		stmt, err := library.GetStmt("sqlFeaturedSpeakerDeleteByOIDKey")
		if err != nil {
			return errors.Wrap(err, `failed to get statement`)
		}
		_, err = tx.Stmt(stmt).Exec(f.OID)
		return err
	}

	if f.EID != "" {
		stmt, err := library.GetStmt("sqlFeaturedSpeakerDeleteByEIDKey")
		if err != nil {
			return errors.Wrap(err, `failed to get statement`)
		}
		_, err = tx.Stmt(stmt).Exec(f.EID)
		return err
	}

	return errors.New("either OID/EID must be filled")
}

func (v *FeaturedSpeakerList) FromRows(rows *sql.Rows, capacity int) error {
	var res []FeaturedSpeaker
	if capacity > 0 {
		res = make([]FeaturedSpeaker, 0, capacity)
	} else {
		res = []FeaturedSpeaker{}
	}

	for rows.Next() {
		vdb := FeaturedSpeaker{}
		if err := vdb.Scan(rows); err != nil {
			return err
		}
		res = append(res, vdb)
	}
	*v = res
	return nil
}

func (v *FeaturedSpeakerList) LoadSinceEID(tx *Tx, since string, limit int) error {
	var s int64
	if id := since; id != "" {
		vdb := FeaturedSpeaker{}
		if err := vdb.LoadByEID(tx, id); err != nil {
			return err
		}

		s = vdb.OID
	}
	return v.LoadSince(tx, s, limit)
}

func (v *FeaturedSpeakerList) LoadSince(tx *Tx, since int64, limit int) error {
	rows, err := tx.Query(`SELECT `+FeaturedSpeakerStdSelectColumns+` FROM `+FeaturedSpeakerTable+` WHERE featured_speakers.oid > ? ORDER BY oid ASC LIMIT `+strconv.Itoa(limit), since)
	if err != nil {
		return err
	}

	if err := v.FromRows(rows, limit); err != nil {
		return err
	}
	return nil
}
