package db

import "github.com/pkg/errors"

func (ccl *ConferenceComponentList) LoadByConferenceID(tx *Tx, cid string) error {
	stmt := getStmtBuf()
	defer releaseStmtBuf(stmt)

	stmt.WriteString(`SELECT `)
	stmt.WriteString(ConferenceComponentStdSelectColumns)
	stmt.WriteString(` FROM `)
	stmt.WriteString(ConferenceComponentTable)
	stmt.WriteString(` WHERE conference_id = ?`)

	rows, err := tx.Query(stmt.String(), cid)
	if err != nil {
		return errors.Wrap(err, "failed to execute query for conference component")
	}

	var res ConferenceComponentList
	for rows.Next() {
		var row ConferenceComponent
		if err := row.Scan(rows); err != nil {
			return errors.Wrap(err, "failed to scan conference component")
		}
		res = append(res, row)
	}

	*ccl = res
	return nil
}
