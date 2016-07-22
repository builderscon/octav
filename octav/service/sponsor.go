package service

import (
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"

	"github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

func (v *Sponsor) populateRowForCreate(vdb *db.Sponsor, payload model.CreateSponsorRequest) error {
	vdb.EID = tools.UUID()

	vdb.ConferenceID = payload.ConferenceID
	vdb.Name = payload.Name
	vdb.LogoURL1 = payload.LogoURL1
	vdb.URL = payload.URL
	vdb.GroupName = payload.GroupName
	vdb.SortOrder = payload.SortOrder

	if payload.LogoURL2.Valid() {
		vdb.LogoURL2.Valid = true
		vdb.LogoURL2.String = payload.LogoURL2.String
	}

	if payload.LogoURL3.Valid() {
		vdb.LogoURL3.Valid = true
		vdb.LogoURL3.String = payload.LogoURL3.String
	}

	return nil
}

func (v *Sponsor) populateRowForUpdate(vdb *db.Sponsor, payload model.UpdateSponsorRequest) error {
	if payload.Name.Valid() {
		vdb.Name = payload.Name.String
	}

	if payload.LogoURL1.Valid() {
		vdb.LogoURL1 = payload.LogoURL1.String
	}

	if payload.URL.Valid() {
		vdb.URL = payload.URL.String
	}

	if payload.GroupName.Valid() {
		vdb.GroupName = payload.GroupName.String
	}

	if payload.SortOrder.Valid() {
		vdb.SortOrder = int(payload.SortOrder.Int)
	}

	if payload.LogoURL2.Valid() {
		vdb.LogoURL2.Valid = true
		vdb.LogoURL2.String = payload.LogoURL2.String
	}

	if payload.LogoURL3.Valid() {
		vdb.LogoURL3.Valid = true
		vdb.LogoURL3.String = payload.LogoURL3.String
	}

	return nil
}

func (v *Sponsor) CreateFromPayload(tx *db.Tx, payload model.AddSponsorRequest, result *model.Sponsor) error {
	su := User{}
	if err := su.IsConferenceAdministrator(tx, payload.ConferenceID, payload.UserID); err != nil {
		return errors.Wrap(err, "creating a featured speaker requires conference administrator privilege")
	}

	vdb := db.Sponsor{}
	if err := v.Create(tx, &vdb, model.CreateSponsorRequest{payload}); err != nil {
		return errors.Wrap(err, "failed to store in database")
	}

	c := model.Sponsor{}
	if err := c.FromRow(vdb); err != nil {
		return errors.Wrap(err, "failed to populate model from database")
	}

	*result = c
	return nil
}

func (v *Sponsor) UpdateFromPayload(tx *db.Tx, payload model.UpdateSponsorRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Sponsor.UpdateFromPayload").BindError(&err)
		defer g.End()
	}

	vdb := db.Sponsor{}
	if err := vdb.LoadByEID(tx, payload.ID); err != nil {
		return errors.Wrap(err, "failed to load featured speaker from database")
	}

	su := User{}
	if err := su.IsConferenceAdministrator(tx, vdb.ConferenceID, payload.UserID); err != nil {
		return errors.Wrap(err, "updating a featured speaker requires conference administrator privilege")
	}

	return errors.Wrap(v.Update(tx, &vdb, payload), "failed to load featured speaker from database")
}

func (v *Sponsor) DeleteFromPayload(tx *db.Tx, payload model.DeleteSponsorRequest) error {
	var m db.Sponsor
	if err := m.LoadByEID(tx, payload.ID); err != nil {
		return errors.Wrap(err, "failed to load featured speaker from database")
	}

	su := User{}
	if err := su.IsConferenceAdministrator(tx, m.ConferenceID, payload.UserID); err != nil {
		return errors.Wrap(err, "deleting venues require administrator privileges")
	}

	return errors.Wrap(v.Delete(tx, m.EID), "failed to delete from database")
}

func (v *Sponsor) ListFromPayload(tx *db.Tx, result *model.SponsorList, payload model.ListSponsorsRequest) error {
	var vdbl db.SponsorList
	if err := vdbl.LoadByConferenceSinceEID(tx, payload.ConferenceID, payload.Since.String, int(payload.Limit.Int)); err != nil {
		return errors.Wrap(err, "failed to load featured speakers from database")
	}

	l := make(model.SponsorList, len(vdbl))
	for i, vdb := range vdbl {
		if err := (l[i]).FromRow(vdb); err != nil {
			return errors.Wrap(err, "failed to populate model from database")
		}

		if err := v.Decorate(tx, &l[i], payload.Lang.String); err != nil {
			return errors.Wrap(err, "failed to decorate venue with associated data")
		}
	}

	*result = l
	return nil
}

func (v *Sponsor) Decorate(tx *db.Tx, speaker *model.Sponsor, lang string) error {
	if lang == "" {
		return nil
	}

	if err := v.ReplaceL10NStrings(tx, speaker, lang); err != nil {
		return errors.Wrap(err, "failed to replace L10N strings")
	}

	return nil
}
