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
	lf := octav.LocalizedFields{}
	lf.Set("ja", "name", `東京ビッグサイト`)
	lf.Set("ja", "address", `〒135-0063 東京都江東区有明３丁目１０−１`)
	in := octav.Venue{
		Name:      "Tokyo Bigsight",
		Address:   "Ariake 3-10-1, Koto-ku, Tokyo",
		L10N:      lf,
		Longitude: 35.6320326,
		Latitude:  139.7976891,
	}
	res, err := cl.CreateVenue(&in)
	if !assert.NotEmpty(t, res.ID, "Returned structure has ID") {
		return
	}

	if !assert.NoError(t, err, "CreateVenue should succeed") {
		return
	}
	if !assert.NoError(t, validator.HTTPCreateVenueResponse.Validate(res), "Validation should succeed") {
		return
	}

	if !assert.NoError(t, cl.DeleteVenue(&octav.DeleteVenueRequest{ID: res.ID}), "Cleanup: delete should be successful") {
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

func TestListSessionsByConference(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)
	var in map[string]interface{}
	res, err := cl.ListSessionsByConference(in)
	if !assert.NoError(t, err, "ListSessionsByConference should succeed") {
		return
	}
	if !assert.NoError(t, validator.HTTPListSessionsByConferenceResponse.Validate(res), "Validation should succeed") {
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

func TestLookupSession(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)
	var in map[string]interface{}
	res, err := cl.LookupSession(in)
	if !assert.NoError(t, err, "LookupSession should succeed") {
		return
	}
	if !assert.NoError(t, validator.HTTPLookupSessionResponse.Validate(res), "Validation should succeed") {
		return
	}
}

func TestLookupVenue(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)
	var in map[string]interface{}
	res, err := cl.LookupVenue(in)
	if !assert.NoError(t, err, "LookupVenue should succeed") {
		return
	}
	if !assert.NoError(t, validator.HTTPLookupVenueResponse.Validate(res), "Validation should succeed") {
		return
	}
}
