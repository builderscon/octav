package octav_test

import (
	"net/http/httptest"
	"testing"

	"github.com/builderscon/octav/octav"
	"github.com/builderscon/octav/octav/client"
	"github.com/builderscon/octav/octav/validator"
	"github.com/stretchr/testify/assert"
)

func TestCreateConference(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)
	var in *octav.Conference
	res, err := cl.CreateConference(in)
	if !assert.NoError(t, err, "CreateConference should succeed") {
		return
	}
	if !assert.NoError(t, validator.HTTPCreateConferenceResponse.Validate(res), "Validation should succeed") {
		return
	}
}

func TestCreateRoom(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)
	var in *octav.Room
	res, err := cl.CreateRoom(in)
	if !assert.NoError(t, err, "CreateRoom should succeed") {
		return
	}
	if !assert.NoError(t, validator.HTTPCreateRoomResponse.Validate(res), "Validation should succeed") {
		return
	}
}

func TestCreateSession(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)
	var in interface{}
	res, err := cl.CreateSession(in)
	if !assert.NoError(t, err, "CreateSession should succeed") {
		return
	}
	if !assert.NoError(t, validator.HTTPCreateSessionResponse.Validate(res), "Validation should succeed") {
		return
	}
}

func TestCreateUser(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)
	var in interface{}
	res, err := cl.CreateUser(in)
	if !assert.NoError(t, err, "CreateUser should succeed") {
		return
	}
	if !assert.NoError(t, validator.HTTPCreateUserResponse.Validate(res), "Validation should succeed") {
		return
	}
}

func TestCreateVenue(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)
	var in *octav.Venue
	res, err := cl.CreateVenue(in)
	if !assert.NoError(t, err, "CreateVenue should succeed") {
		return
	}
	if !assert.NoError(t, validator.HTTPCreateVenueResponse.Validate(res), "Validation should succeed") {
		return
	}
}

func TestListRooms(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)
	var in map[string]interface{}
	res, err := cl.ListRooms(in)
	if !assert.NoError(t, err, "ListRooms should succeed") {
		return
	}
	if !assert.NoError(t, validator.HTTPListRoomsResponse.Validate(res), "Validation should succeed") {
		return
	}
}

func TestListVenues(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)
	var in map[string]interface{}
	res, err := cl.ListVenues(in)
	if !assert.NoError(t, err, "ListVenues should succeed") {
		return
	}
	if !assert.NoError(t, validator.HTTPListVenuesResponse.Validate(res), "Validation should succeed") {
		return
	}
}
