package service

import (
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	pdebug "github.com/lestrrat/go-pdebug"
)

func (v *ConferenceVenueSvc) Init() {}

func (v *ConferenceVenueSvc) populateRowForCreate(vdb *db.ConferenceVenue, payload *model.CreateConferenceVenueRequest) error {
	vdb.ConferenceID = payload.ConferenceID
	vdb.VenueID = payload.VenueID
	return nil
}

func (v *ConferenceVenueSvc) populateRowForUpdate(vdb *db.ConferenceVenue, payload *model.UpdateConferenceVenueRequest) error {
	vdb.ConferenceID = payload.ConferenceID
	vdb.VenueID = payload.VenueID
	return nil
}

func invalidateVenueLoadByConferenceID(conferenceID string) error {
	c := Cache()
	key := c.Key("Venue", "LoadByConferenceID", conferenceID)
	c.Delete(key)
	if pdebug.Enabled {
		pdebug.Printf("CACHE DEL: %s", key)
	}
	return nil
}

func (v *ConferenceVenueSvc) PostCreateHook(tx *db.Tx, vdb *db.ConferenceVenue) error {
	invalidateVenueLoadByConferenceID(vdb.ConferenceID)
	invalidateRoomLoadByVenueID(vdb.VenueID)
	return nil
}

func (v *ConferenceVenueSvc) PostUpdateHook(tx *db.Tx, vdb *db.ConferenceVenue) error {
	invalidateVenueLoadByConferenceID(vdb.ConferenceID)
	invalidateRoomLoadByVenueID(vdb.VenueID)
	return nil
}

func (v *ConferenceVenueSvc) PostDeleteHook(tx *db.Tx, vdb *db.ConferenceVenue) error {
	invalidateVenueLoadByConferenceID(vdb.ConferenceID)
	invalidateRoomLoadByVenueID(vdb.VenueID)
	return nil
}
