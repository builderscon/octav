package octav

import "github.com/builderscon/octav/octav/db"

func (v *SessionList) LoadByConference(tx *db.Tx, cid string, date string) error {
	vdbl := db.SessionList{}
	if err := vdbl.LoadByConference(tx, cid, date); err != nil {
		return err
	}

	res := make([]Session, len(vdbl))
	for i, vdb := range vdbl {
		v := Session{}
		if err := v.FromRow(vdb); err != nil {
			return err
		}
		if err := v.LoadLocalizedFields(tx); err != nil {
			return err
		}
		res[i] = v
	}
	*v = res
	return nil
}
