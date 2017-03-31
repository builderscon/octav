package service

import (
	"context"
	"testing"

	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/stretchr/testify/assert"
)

func TestSessionPopulateRowForUpdate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := Session()

	var vdb db.Session
	var payload model.UpdateSessionRequest

	payload.ID = "abc"
	payload.Status.Set("accepted")
	vdb.EID = payload.ID

	if !assert.NoError(t, s.populateRowForUpdate(ctx, &vdb, &payload), "populateRowForUpdate should succeed") {
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

func TestFormatSessionTweet(t *testing.T) {
	series := model.ConferenceSeries{
		Slug: "builderscon",
	}
	conf := model.Conference{
		Slug: "tokyo/2016",
	}

	session := model.Session{
		ID:    "ff8657cb-a751-4415-ad93-374fb9fda2b6",
		Title: "Highly available and scalable Kubernetes on AWS",
	}

	s, err := formatSessionTweet(&session, &conf, &series)
	if !assert.NoError(t, err, "formatSessionTweet should succeed") {
		return
	}

	if !assert.Equal(t, `New submission "Highly available and scalable Kubernetes on AWS" https://builderscon.io/builderscon/tokyo/2016/session/ff8657cb-a751-4415-ad93-374fb9fda2b6`, s, "tweet should match") {
		return
	}
}
