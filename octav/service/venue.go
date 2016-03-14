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