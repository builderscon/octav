package db

import "bytes"

func (cd *ConferenceDate) DeleteDates(tx *Tx, cid string, dates ...string) error {
	stmt := bytes.Buffer{}
	stmt.WriteString("DELETE FROM ")
	stmt.WriteString(ConferenceDateTable)
	stmt.WriteString(" WHERE conference_id = ?")
	var args []interface{}
	switch l := len(dates); l {
	case 0:
		args = make([]interface{}, 1)
	default:
		args = make([]interface{}, l+1)
		stmt.WriteString(" AND date IN (")
		for i := 0; i < l; i++ {
			stmt.WriteByte('?')
			if i < l-1 {
				stmt.WriteString(", ")
			}
			args[i+1] = dates[i]
		}
		stmt.WriteByte(')')
	}
	args[0] = cid

	_, err := tx.Exec(stmt.String(), args...)
	return err
}

func (cdl *ConferenceDateList) LoadByConferenceID(tx *Tx, cid string) error {
	stmt := bytes.Buffer{}
	stmt.WriteString("SELECT ")
	stmt.WriteString(ConferenceDateStdSelectColumns)
	stmt.WriteString(" FROM ")
	stmt.WriteString(ConferenceDateTable)
	stmt.WriteString(" WHERE conference_id = ? ORDER BY date,open ASC")
	rows, err := tx.Query(stmt.String(), cid)
	if err != nil {
		return err
	}

	res := ConferenceDateList{}
	for rows.Next() {
		var cd ConferenceDate
		if err := cd.Scan(rows); err != nil {
			return err
		}
		res = append(res, cd)
	}
	*cdl = res
	return nil
}
