package octav

import (
	"github.com/builderscon/octav/octav/db"
)

func (v *Conference) Load(tx *db.Tx, id string) error {
	vdb := db.Conference{}
	if err := vdb.LoadByEID(tx, id); err != nil {
		return err
	}

	return v.FromRow(vdb)
}

func (v *Conference) FromRow(vdb db.Conference) error {
	v.ID = vdb.EID
	v.Slug = vdb.Slug
	v.Title = vdb.Title
	if vdb.SubTitle.Valid {
		v.SubTitle = vdb.SubTitle.String
	}
	return nil
}

func (v Conference) ToRow(vdb *db.Conference) error {
	vdb.EID = v.ID
	vdb.Slug = v.Slug
	vdb.Title = v.Title
	vdb.SubTitle.Valid = true
	vdb.SubTitle.String = v.SubTitle
	return nil
}

func (v *ConferenceList) Load(tx *db.Tx, since string) error {
	var s int64
	if id := since; id != "" {
		vdb := db.Conference{}
		if err := vdb.LoadByEID(tx, id); err != nil {
			return err
		}

		s = vdb.OID
	}

	rows, err := tx.Query(`SELECT eid, name FROM venue WHERE oid > ? ORDER BY oid LIMIT 10`, s)
	if err != nil {
		return err
	}

	// Not using db.Conference here
	res := make([]Conference, 0, 10)
	for rows.Next() {
		vdb := db.Conference{}
		if err := vdb.Scan(rows); err != nil {
			return err
		}
		v := Conference{}
		if err := v.FromRow(vdb); err != nil {
			return err
		}
		res = append(res, v)
	}
	*v = res
	return nil
}
