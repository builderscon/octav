package octav_test

import (
	"net/http/httptest"
	"testing"

	"github.com/builderscon/octav/octav"
	"github.com/builderscon/octav/octav/client"
	"github.com/builderscon/octav/octav/validator"
	"github.com/stretchr/testify/assert"
)

func bigsight() *octav.Venue {
	lf := octav.LocalizedFields{}
	lf.Set("ja", "name", `東京ビッグサイト`)
	lf.Set("ja", "address", `〒135-0063 東京都江東区有明３丁目１０−１`)
	return &octav.Venue{
		Name:      "Tokyo Bigsight",
		Address:   "Ariake 3-10-1, Koto-ku, Tokyo",
		L10N:      lf,
		Longitude: 35.6320326,
		Latitude:  139.7976891,
	}
}

func intlConferenceRoom(venueID string) *octav.Room {
	lf := octav.LocalizedFields{}
	lf.Set("ja", "name", `国際会議場`)
	return &octav.Room{
		Capacity: 1000,
		L10N:     lf,
		Name:     "International Conference Hall",
		VenueID:  venueID,
	}
}

func testCreateRoom(t *testing.T, cl *client.Client, r *octav.Room) (*octav.Room, error) {
	res, err := cl.CreateRoom(r)
	if !assert.NoError(t, err, "CreateRoom should succeed") {
		return nil, err
	}
	return res, nil
}

func testCreateVenue(t *testing.T, cl *client.Client, v *octav.Venue) (*octav.Venue, error) {
	res, err := cl.CreateVenue(v)
	if !assert.NoError(t, err, "CreateVenue should succeed") {
		return nil, err
	}
	return res, nil
}

func testLookupRoom(t *testing.T, cl *client.Client, id string) (*octav.Room, error) {
	venue, err := cl.LookupRoom(&octav.LookupRoomRequest{ID: id})
	if !assert.NoError(t, err, "LookupRoom succeeds") {
		return nil, err
	}
	return venue, nil
}

func testCreateUser(t *testing.T, cl *client.Client, in *octav.CreateUserRequest) (*octav.User, error) {
	res, err := cl.CreateUser(in)
	if !assert.NoError(t, err, "CreateUser should succeed") {
		return nil, err
	}
	return res, nil
}

func testLookupUser(t *testing.T, cl *client.Client, id string) (*octav.User, error) {
	user, err := cl.LookupUser(&octav.LookupUserRequest{ID: id})
	if !assert.NoError(t, err, "LookupUser succeeds") {
		return nil, err
	}
	return user, nil
}

func testDeleteUser(t *testing.T, cl *client.Client, id string) error {
	err := cl.DeleteUser(&octav.DeleteUserRequest{ID: id})
	if !assert.NoError(t, err, "DeleteUser should succeed") {
		return err
	}
	return nil
}

func testLookupVenue(t *testing.T, cl *client.Client, id string) (*octav.Venue, error) {
	venue, err := cl.LookupVenue(&octav.LookupVenueRequest{ID: id})
	if !assert.NoError(t, err, "LookupVenue succeeds") {
		return nil, err
	}
	return venue, nil
}

func testDeleteRoom(t *testing.T, cl *client.Client, id string) error {
	err := cl.DeleteRoom(&octav.DeleteRoomRequest{ID: id})
	if !assert.NoError(t, err, "DeleteRoom should be successful") {
		return err
	}
	return err
}

func testDeleteVenue(t *testing.T, cl *client.Client, id string) error {
	err := cl.DeleteVenue(&octav.DeleteVenueRequest{ID: id})
	if !assert.NoError(t, err, "DeleteVenue should be successful") {
		return err
	}
	return err
}

func yapcasia() *octav.CreateConferenceRequest {
	return &octav.CreateConferenceRequest{
		Title: "YAPC::Asia Tokyo",
		Slug: "yapcasia",
	}
}

func testCreateConference(t *testing.T, cl *client.Client, in *octav.CreateConferenceRequest) (*octav.Conference, error) {
	res, err := cl.CreateConference(in)
	if !assert.NoError(t, err, "CreateConference should succeed") {
		return nil, err
	}
	return res, nil
}

func testLookupConference(t *testing.T, cl *client.Client, id string) (*octav.Conference, error) {
	venue, err := cl.LookupConference(&octav.LookupConferenceRequest{ID: id})
	if !assert.NoError(t, err, "LookupConference succeeds") {
		return nil, err
	}
	return venue, nil
}

func testDeleteConference(t *testing.T, cl *client.Client, id string) error {
	err := cl.DeleteConference(&octav.DeleteConferenceRequest{ID: id})
	if !assert.NoError(t, err, "DeleteConference should be successful") {
		return err
	}
	return err
}

func TestCreateConference(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)
	res, err := testCreateConference(t, cl, yapcasia())
	if err != nil {
		return
	}

	if !assert.NoError(t, validator.HTTPCreateConferenceResponse.Validate(res), "Validation should succeed") {
		return
	}

	res2, err := testLookupConference(t, cl, res.ID)
	if err != nil {
		return
	}

	if !assert.Equal(t, res2, res, "LookupConference is the same as the conference created") {
		return
	}

	if err := testDeleteConference(t, cl, res.ID); err != nil {
		return
	}
}

func TestCreateRoom(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)

	venue, err := testCreateVenue(t, cl, bigsight())
	if err != nil {
		return
	}

	res, err := testCreateRoom(t, cl, intlConferenceRoom(venue.ID))
	if err != nil {
		return
	}

	if !assert.NotEmpty(t, res.ID, "Returned structure has ID") {
		return
	}

	if !assert.NoError(t, validator.HTTPCreateRoomResponse.Validate(res), "Validation should succeed") {
		return
	}

	res2, err := testLookupRoom(t, cl, res.ID)
	if err != nil {
		return
	}

	if !assert.Equal(t, res2, res, "LookupRoom is the same as the room created") {
		return
	}

	if err := testDeleteRoom(t, cl, res.ID); err != nil {
		return
	}

	if err := testDeleteVenue(t, cl, venue.ID); err != nil {
		return
	}
}

func testCreateSession(t *testing.T, cl *client.Client, in *octav.CreateSessionRequest) (*octav.Session, error) {
	res, err := cl.CreateSession(in)
	if !assert.NoError(t, err, "CreateSession should succeed") {
		return nil, err
	}
	return res, nil
}

func TestCreateSession(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)

	conference, err := testCreateConference(t, cl, yapcasia())
	if err != nil {
		return
	}

	user, err := testCreateUser(t, cl, johndoe())
	if err != nil {
		return
	}

	in := octav.CreateSessionRequest{
		ConferenceID: conference.ID,
		SpeakerID: user.ID,
		Title: "How To Write A Conference Backend",
	}
	res, err := testCreateSession(t, cl, &in)
	if err != nil {
		return
	}

t.Logf("%#v", res)
	if !assert.NoError(t, validator.HTTPCreateSessionResponse.Validate(res), "Validation should succeed") {
		return
	}
}

func johndoe() *octav.CreateUserRequest {
	lf := octav.LocalizedFields{}
	lf.Set("ja", "first_name", "ジョン")
	lf.Set("ja", "last_name", "ドー")
	return &octav.CreateUserRequest{
		FirstName:  "John",
		LastName:   "Doe",
		Nickname:   "enigma621",
		Email:      "john.doe@example.com",
		TshirtSize: "XL",
		L10N:       lf,
	}
}

func TestCreateUser(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)
	res, err := testCreateUser(t, cl, johndoe())
	if err != nil {
		return
	}

	if !assert.NotEmpty(t, res.ID, "Returned structure has ID") {
		return
	}

	if !assert.NoError(t, validator.HTTPCreateUserResponse.Validate(res), "Validation should succeed") {
		return
	}

	res2, err := testLookupUser(t, cl, res.ID)
	if err != nil {
		return
	}

	if !assert.Equal(t, res2, res, "LookupUser is the same as the user created") {
		return
	}

	if err := testDeleteUser(t, cl, res.ID); err != nil {
		return
	}
}

func TestCreateVenue(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)
	res, err := testCreateVenue(t, cl, bigsight())
	if err != nil {
		return
	}

	if !assert.NotEmpty(t, res.ID, "Returned structure has ID") {
		return
	}

	if !assert.NoError(t, validator.HTTPCreateVenueResponse.Validate(res), "Validation should succeed") {
		return
	}

	res2, err := testLookupVenue(t, cl, res.ID)
	if err != nil {
		return
	}

	if !assert.Equal(t, res2, res, "LookupVenue is the same as the venue created") {
		return
	}

	if err := testDeleteVenue(t, cl, res.ID); err != nil {
		return
	}
}

type setPropValuer interface {
  SetPropValue(string, interface{}) error
}

func TestListRooms(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)
	venue, err := testCreateVenue(t, cl, bigsight())
	if err != nil {
		return
	}

	_, err = testCreateRoom(t, cl, intlConferenceRoom(venue.ID))
	if err != nil {
		return
	}

	in := octav.ListRoomRequest{
		VenueID: venue.ID,
	}
	res, err := cl.ListRooms(&in)
	if !assert.NoError(t, err, "ListRooms should succeed") {
		return
	}

	if !assert.NoError(t, validator.HTTPListRoomsResponse.Validate(res), "Validation should succeed") {
		return
	}

	if !assert.Len(t, res, 1, "ListRooms returns 1 rooms") {
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
	in := octav.ListVenueRequest{}
	res, err := cl.ListVenues(&in)
	if !assert.NoError(t, err, "ListVenues should succeed") {
		return
	}
	if !assert.NoError(t, validator.HTTPListVenuesResponse.Validate(res), "Validation should succeed") {
		return
	}
}
