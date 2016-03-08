package octav

import (
	"github.com/builderscon/octav/octav/db"
	"github.com/lestrrat/go-pdebug"
)

func (v *User) Load(tx *db.Tx, id string) error {
	vdb := db.User{}
	if err := vdb.LoadByEID(tx, id); err != nil {
		return err
	}

	v.ID = vdb.EID
	v.FirstName = vdb.FirstName
	v.LastName = vdb.LastName
	v.Nickname = vdb.Nickname
	v.Email = vdb.Email
	v.TshirtSize = vdb.TshirtSize

	ls, err := db.LoadLocalizedStringsForParent(tx, v.ID, "User")
	if err != nil {
		return err
	}

	if len(ls) > 0 {
		v.L10N = LocalizedFields{}
		for _, l := range ls {
			v.L10N.Set(l.Language, l.Name, l.Localized)
		}
	}

	return nil
}

func (v *User) Create(tx *db.Tx) error {
	if v.ID == "" {
		v.ID = UUID()
	}

	vdb := db.User{
		EID:        v.ID,
		FirstName:  v.FirstName,
		LastName:   v.LastName,
		Nickname:   v.Nickname,
		Email:      v.Email,
		TshirtSize: v.TshirtSize,
	}
	if err := vdb.Create(tx); err != nil {
		return err
	}

	if err := v.L10N.CreateLocalizedStrings(tx, "Venue", v.ID); err != nil {
		return err
	}
	return nil
}

func (v *User) Delete(tx *db.Tx) error {
	if pdebug.Enabled {
		g := pdebug.Marker("User.Delete (%s)", v.ID)
		defer g.End()
	}

	vdb := db.User{EID: v.ID}
	if err := vdb.Delete(tx); err != nil {
		return err
	}

	if err := db.DeleteLocalizedStringsForParent(tx, v.ID, "User"); err != nil {
		return err
	}
	return nil
}
