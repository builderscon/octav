package db

// Automatically generated by gendb utility. DO NOT EDIT!

import (
	"bytes"
	"database/sql"
	"strconv"
	"time"

	"github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

const QuestionStdSelectColumns = "questions.oid, questions.eid, questions.session_id, questions.user_id, questions.body, questions.created_on, questions.modified_on"
const QuestionTable = "questions"

type QuestionList []Question

func (q *Question) Scan(scanner interface {
	Scan(...interface{}) error
}) error {
	return scanner.Scan(&q.OID, &q.EID, &q.SessionID, &q.UserID, &q.Body, &q.CreatedOn, &q.ModifiedOn)
}

func (q *Question) LoadByEID(tx *Tx, eid string) error {
	row := tx.QueryRow(`SELECT `+QuestionStdSelectColumns+` FROM `+QuestionTable+` WHERE questions.eid = ?`, eid)
	if err := q.Scan(row); err != nil {
		return err
	}
	return nil
}

func (q *Question) Create(tx *Tx, opts ...InsertOption) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.Question.Create").BindError(&err)
		defer g.End()
		pdebug.Printf("%#v", q)
	}
	if q.EID == "" {
		return errors.New("create: non-empty EID required")
	}

	q.CreatedOn = time.Now()
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
	stmt.WriteString(QuestionTable)
	stmt.WriteString(` (eid, session_id, user_id, body, created_on, modified_on) VALUES (?, ?, ?, ?, ?, ?)`)
	result, err := tx.Exec(stmt.String(), q.EID, q.SessionID, q.UserID, q.Body, q.CreatedOn, q.ModifiedOn)
	if err != nil {
		return err
	}

	lii, err := result.LastInsertId()
	if err != nil {
		return err
	}

	q.OID = lii
	return nil
}

func (q Question) Update(tx *Tx) error {
	if q.OID != 0 {
		_, err := tx.Exec(`UPDATE `+QuestionTable+` SET eid = ?, session_id = ?, user_id = ?, body = ? WHERE oid = ?`, q.EID, q.SessionID, q.UserID, q.Body, q.OID)
		return err
	}
	if q.EID != "" {
		_, err := tx.Exec(`UPDATE `+QuestionTable+` SET session_id = ?, user_id = ?, body = ? WHERE eid = ?`, q.SessionID, q.UserID, q.Body, q.EID)
		return err
	}
	return errors.New("either OID/EID must be filled")
}

func (q Question) Delete(tx *Tx) error {
	if q.OID != 0 {
		_, err := tx.Exec(`DELETE FROM `+QuestionTable+` WHERE oid = ?`, q.OID)
		return err
	}

	if q.EID != "" {
		_, err := tx.Exec(`DELETE FROM `+QuestionTable+` WHERE eid = ?`, q.EID)
		return err
	}

	return errors.New("either OID/EID must be filled")
}

func (v *QuestionList) FromRows(rows *sql.Rows, capacity int) error {
	var res []Question
	if capacity > 0 {
		res = make([]Question, 0, capacity)
	} else {
		res = []Question{}
	}

	for rows.Next() {
		vdb := Question{}
		if err := vdb.Scan(rows); err != nil {
			return err
		}
		res = append(res, vdb)
	}
	*v = res
	return nil
}

func (v *QuestionList) LoadSinceEID(tx *Tx, since string, limit int) error {
	var s int64
	if id := since; id != "" {
		vdb := Question{}
		if err := vdb.LoadByEID(tx, id); err != nil {
			return err
		}

		s = vdb.OID
	}
	return v.LoadSince(tx, s, limit)
}

func (v *QuestionList) LoadSince(tx *Tx, since int64, limit int) error {
	rows, err := tx.Query(`SELECT `+QuestionStdSelectColumns+` FROM `+QuestionTable+` WHERE questions.oid > ? ORDER BY oid ASC LIMIT `+strconv.Itoa(limit), since)
	if err != nil {
		return err
	}

	if err := v.FromRows(rows, limit); err != nil {
		return err
	}
	return nil
}
