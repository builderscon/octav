package service

import (
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	"github.com/pkg/errors"
)

func (v *Room) populateRowForCreate(vdb *db.Room, payload model.CreateRoomRequest) error {
	vdb.EID = tools.UUID()

	if payload.VenueID.Valid() {
		vdb.VenueID = payload.VenueID.String
	}

	if payload.Name.Valid() {
		vdb.Name = payload.Name.String
	}

	if payload.Capacity.Valid() {
		vdb.Capacity = uint(payload.Capacity.Uint)
	}

	return nil
}

func (v *Room) populateRowForUpdate(vdb *db.Room, payload model.UpdateRoomRequest) error {
	if payload.VenueID.Valid() {
		vdb.VenueID = payload.VenueID.String
	}

	if payload.Name.Valid() {
		vdb.Name = payload.Name.String
	}

	if payload.Capacity.Valid() {
		vdb.Capacity = uint(payload.Capacity.Uint)
	}

	return nil
}

func (v *Room) CreateFromPayload(tx *db.Tx, result *model.Room, payload model.CreateRoomRequest) error {
	su := User{}
	if err := su.IsAdministrator(tx, payload.UserID); err != nil {
		return errors.Wrap(err, "creating a room requires conference administrator privilege")
	}

	vdb := db.Room{}
	if err := v.Create(tx, &vdb, payload); err != nil {
		return errors.Wrap(err, "failed to store in database")
	}

	var r model.Room
	if err := r.FromRow(vdb); err != nil {
		return errors.Wrap(err, "failed to populate model from database")
	}

	*result = r
	return nil
}

func (v *Room) ListFromPayload(tx *db.Tx, result *model.RoomList, payload model.ListRoomRequest) error {
	var m model.RoomList
	if err := m.LoadForVenue(tx, payload.VenueID, payload.Since.String, int(payload.Limit.Int)); err != nil {
		return errors.Wrap(err, "failed to load from database")
	}

	for i := range m {
		if err := v.Decorate(tx, &m[i], payload.Lang.String); err != nil {
			return errors.Wrap(err, "failed to associate data to model")
		}
	}

	*result = m
	return nil
}

func (v *Room) UpdateFromPayload(tx *db.Tx, payload model.UpdateRoomRequest) error {
	su := User{}
	if err := su.IsAdministrator(tx, payload.UserID); err != nil {
		return errors.Wrap(err, "deleting rooms require administrator privileges")
	}

	var vdb db.Room
	if err := vdb.LoadByEID(tx, payload.ID); err != nil {
		return errors.Wrap(err, "failed to load from database")
	}

	return errors.Wrap(v.Update(tx, &vdb, payload), "failed to update database")
}

func (v *Room) DeleteFromPayload(tx *db.Tx, payload model.DeleteRoomRequest) error {
	su := User{}
	if err := su.IsAdministrator(tx, payload.UserID); err != nil {
		return errors.Wrap(err, "deleting rooms require administrator privileges")
	}

	return errors.Wrap(v.Delete(tx, payload.ID), "failed to delete from ddatabase")
}

func (v *Room) Decorate(tx *db.Tx, room *model.Room, lang string) error {
	if lang != "" {
		if err := v.ReplaceL10NStrings(tx, room, lang); err != nil {
			return errors.Wrap(err, "failed to replace L10N strings")
		}
	}
	return nil
}