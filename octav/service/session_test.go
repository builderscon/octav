package service

import (
	"testing"

	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/stretchr/testify/assert"
)

func TestSessionPopulateRowForUpdate(t *testing.T) {
	s := Session()

	var vdb db.Session
	var payload model.UpdateSessionRequest

	payload.ID = "abc"
	payload.Status.Set("accepted")
	vdb.EID = payload.ID

	if !assert.NoError(t, s.populateRowForUpdate(&vdb, payload), "populateRowForUpdate should succeed") {
		return
	}

	if !assert.Equal(t, payload.ID, vdb.EID, "ID should match") {
		return
	}

	if !assert.Equal(t, payload.Status.String, vdb.Status, "Status should match") {
		return
	}

	if !assert.False(t, vdb.StartsOn.Valid, "StartsOn should be invalid") {
		return
	}
}
