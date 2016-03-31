package service

import (
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
)

func (v *Venue) populateRowForCreate(vdb *db.Venue, payload model.CreateVenueRequest) error {
	vdb.EID = tools.UUID()
	vdb.Name = payload.Name.String
	vdb.Address = payload.Address.String
	vdb.Latitude = payload.Latitude.Float
	vdb.Longitude = payload.Longitude.Float
	return nil
}

func (v *Venue) populateRowForUpdate(vdb *db.Venue, payload model.UpdateVenueRequest) error {
	if payload.Name.Valid() {
	vdb.Name = payload.Name.String
	}

	if payload.Address.Valid() {
	vdb.Address = payload.Address.String
	}

	if payload.Latitude.Valid() {
		vdb.Latitude = payload.Latitude.Float
	}

	if payload.Longitude.Valid() {
		vdb.Longitude = payload.Longitude.Float
	}
	return nil
}

func (v *Venue) LoadRooms(tx *db.Tx, cdl *model.RoomList, vid string) error {
	var vdbl db.RoomList
	if err := db.LoadVenueRooms(tx, &vdbl, vid); err != nil {
		return err
	}

	res := make(model.RoomList, len(vdbl))
	for i, vdb := range vdbl {
		var u model.Room
		if err := u.FromRow(vdb); err != nil {
			return err
		}
		res[i] = u
	}
	*cdl = res
	return nil
}

func (v *Venue) Decorate(tx *db.Tx, venue *model.Venue) error {
	if err := v.LoadRooms(tx, &venue.Rooms, venue.ID); err != nil {
		return err
	}
	return nil
}