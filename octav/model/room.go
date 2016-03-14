package model

import "github.com/builderscon/octav/octav/db"

func (v *RoomList) LoadForVenue(tx *db.Tx, venueID, since string, limit int) error {
	vdbl := db.RoomList{}
	if err := vdbl.LoadForVenueSinceEID(tx, venueID, since, limit); err != nil {
		return err
	}

	res := make([]RoomL10N, len(vdbl))
	for i, vdb := range vdbl {
		v := Room{}
		if err := v.FromRow(vdb); err != nil {
			return err
		}
		vl := RoomL10N{Room: v}
		if err := vl.LoadLocalizedFields(tx); err != nil {
			return err
		}
		res[i] = vl
	}
	*v = res
	return nil
}
