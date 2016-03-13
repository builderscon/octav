package octav_test

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/builderscon/octav/octav"
	"github.com/builderscon/octav/octav/client"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/service"
	"github.com/builderscon/octav/octav/tools"
	"github.com/builderscon/octav/octav/validator"
	"github.com/stretchr/testify/assert"
)

func bigsight() *service.CreateVenueRequest {
	lf := tools.LocalizedFields{}
	lf.Set("ja", "name", `東京ビッグサイト`)
	lf.Set("ja", "address", `〒135-0063 東京都江東区有明３丁目１０−１`)

	r := service.CreateVenueRequest{}
	r.L10N = lf
	r.Name.Set("Tokyo Bigsight")
	r.Address.Set("Ariake 3-10-1, Koto-ku, Tokyo")
	r.Longitude.Set(35.6320326)
	r.Latitude.Set(139.7976891)

	return &r
}

func intlConferenceRoom(venueID string) *service.CreateRoomRequest {
	lf := tools.LocalizedFields{}
	lf.Set("ja", "name", `国際会議場`)

	r := service.CreateRoomRequest{}
	r.L10N = lf
	r.Capacity.Set(1000)
	r.Name.Set("International Conference Hall")
	r.VenueID.Set(venueID)

	return &r
}

func testCreateRoom(t *testing.T, cl *client.Client, r *service.CreateRoomRequest) (*model.Room, error) {
	res, err := cl.CreateRoom(r)
	if !assert.NoError(t, err, "CreateRoom should succeed") {
		return nil, err
	}
	return res, nil
}

func testCreateVenue(t *testing.T, cl *client.Client, v *service.CreateVenueRequest) (*model.Venue, error) {
	res, err := cl.CreateVenue(v)
	if !assert.NoError(t, err, "CreateVenue should succeed") {
		return nil, err
	}
	return res, nil
}

func testLookupRoom(t *testing.T, cl *client.Client, id, lang string) (*model.Room, error) {
	r := &service.LookupRoomRequest{ID: id}
	if lang != "" {
		r.Lang.Set(lang)
	}
	venue, err := cl.LookupRoom(r)
	if !assert.NoError(t, err, "LookupRoom succeeds") {
		return nil, err
	}
	return venue, nil
}

func testCreateUser(t *testing.T, cl *client.Client, in *service.CreateUserRequest) (*model.User, error) {
	res, err := cl.CreateUser(in)
	if !assert.NoError(t, err, "CreateUser should succeed") {
		return nil, err
	}
	return res, nil
}

func testLookupUser(t *testing.T, cl *client.Client, id string) (*model.User, error) {
	user, err := cl.LookupUser(&service.LookupUserRequest{ID: id})
	if !assert.NoError(t, err, "LookupUser succeeds") {
		return nil, err
	}
	return user, nil
}

func testDeleteUser(t *testing.T, cl *client.Client, id string) error {
	err := cl.DeleteUser(&service.DeleteUserRequest{ID: id})
	if !assert.NoError(t, err, "DeleteUser should succeed") {
		return err
	}
	return nil
}

func testLookupVenue(t *testing.T, cl *client.Client, id string) (*model.Venue, error) {
	venue, err := cl.LookupVenue(&service.LookupVenueRequest{ID: id})
	if !assert.NoError(t, err, "LookupVenue succeeds") {
		return nil, err
	}
	return venue, nil
}

func testUpdateRoom(t *testing.T, cl *client.Client, in *service.UpdateRoomRequest) error {
	err := cl.UpdateRoom(in)
	if !assert.NoError(t, err, "UpdateRoom succeeds") {
		return err
	}
	return nil
}

func testDeleteRoom(t *testing.T, cl *client.Client, id string) error {
	err := cl.DeleteRoom(&service.DeleteRoomRequest{ID: id})
	if !assert.NoError(t, err, "DeleteRoom should be successful") {
		return err
	}
	return err
}

func testDeleteVenue(t *testing.T, cl *client.Client, id string) error {
	err := cl.DeleteVenue(&service.DeleteVenueRequest{ID: id})
	if !assert.NoError(t, err, "DeleteVenue should be successful") {
		return err
	}
	return err
}

func yapcasia() *service.CreateConferenceRequest {
	return &service.CreateConferenceRequest{
		Title: "YAPC::Asia Tokyo",
		Slug:  "yapcasia",
	}
}

func testCreateConference(t *testing.T, cl *client.Client, in *service.CreateConferenceRequest) (*model.Conference, error) {
	res, err := cl.CreateConference(in)
	if !assert.NoError(t, err, "CreateConference should succeed") {
		return nil, err
	}
	return res, nil
}

func testLookupConference(t *testing.T, cl *client.Client, id, lang string) (*model.Conference, error) {
	r := &service.LookupConferenceRequest{ID: id}
	if lang != "" {
		r.Lang.Set(lang)
	}
	conference, err := cl.LookupConference(r)
	if !assert.NoError(t, err, "LookupConference succeeds") {
		return nil, err
	}
	return conference, nil
}

func testUpdateConference(t *testing.T, cl *client.Client, in *service.UpdateConferenceRequest) error {
	err := cl.UpdateConference(in)
	if !assert.NoError(t, err, "UpdateConference succeeds") {
		return err
	}
	return nil
}

func testDeleteConference(t *testing.T, cl *client.Client, id string) error {
	err := cl.DeleteConference(&service.DeleteConferenceRequest{ID: id})
	if !assert.NoError(t, err, "DeleteConference should be successful") {
		return err
	}
	return err
}

func TestConferenceCRUD(t *testing.T) {
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

	res2, err := testLookupConference(t, cl, res.ID, "")
	if err != nil {
		return
	}

	if !assert.Equal(t, res2, res, "LookupConference is the same as the conference created") {
		return
	}

	in := service.UpdateConferenceRequest{ID: res.ID}
	in.SubTitle.Set("Big Bang!")
	in.L10N.Set("ja", "title", "ヤップシー エイジア")
	if err := testUpdateConference(t, cl, &in); err != nil {
		return
	}

	res3, err := testLookupConference(t, cl, res.ID, "ja")
	if err != nil {
		return
	}

	if !assert.Equal(t, res3.SubTitle, "Big Bang!", "Conference.SubTitle is the same as the conference updated") {
		return
	}

	if !assert.Equal(t, "ヤップシー エイジア", res3.Title, "Conference.title#ja is the same as the conference updated") {
		return
	}

	if err := testDeleteConference(t, cl, res.ID); err != nil {
		return
	}
}

func TestRoomCRUD(t *testing.T) {
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

	res2, err := testLookupRoom(t, cl, res.ID, "")
	if err != nil {
		return
	}

	if !assert.Equal(t, res2, res, "LookupRoom is the same as the room created") {
		return
	}

	in := service.UpdateRoomRequest{ID: res.ID}
	in.L10N.Set("ja", "name", "国際会議場")
	if err := testUpdateRoom(t, cl, &in); err != nil {
		return
	}

	res3, err := testLookupRoom(t, cl, res.ID, "ja")
	if err != nil {
		return
	}

	if !assert.Equal(t, "国際会議場", res3.Name, "Room.name#ja is the same as the conference updated") {
		return
	}

	if err := testDeleteRoom(t, cl, res.ID); err != nil {
		return
	}

	if err := testDeleteVenue(t, cl, venue.ID); err != nil {
		return
	}
}

func testCreateSession(t *testing.T, cl *client.Client, in *service.CreateSessionRequest) (*model.Session, error) {
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

	in := service.CreateSessionRequest{}
	in.ConferenceID.Set(conference.ID)
	in.SpeakerID.Set(user.ID)
	in.Title.Set("How To Write A Conference Backend")
	in.Duration.Set(60)
	in.Abstract.Set("Use lots of reflection and generate lots of code")
	res, err := testCreateSession(t, cl, &in)
	if err != nil {
		return
	}

	if !assert.NoError(t, validator.HTTPCreateSessionResponse.Validate(res), "Validation should succeed") {
		return
	}
}

func johndoe() *service.CreateUserRequest {
	lf := tools.LocalizedFields{}
	lf.Set("ja", "first_name", "ジョン")
	lf.Set("ja", "last_name", "ドー")
	return &service.CreateUserRequest{
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

	in := service.ListRoomRequest{
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
	conference, err := testCreateConference(t, cl, yapcasia())
	if err != nil {
		return
	}

	user, err := testCreateUser(t, cl, johndoe())
	if err != nil {
		return
	}

	for i := 0; i < 10; i++ {
		sin := service.CreateSessionRequest{}
		sin.ConferenceID.Set(conference.ID)
		sin.SpeakerID.Set(user.ID)
		sin.Title.Set(fmt.Sprintf("Title %d", i))
		sin.Duration.Set(60)
		sin.Abstract.Set("Use lots of reflection and generate lots of code")
		_, err := testCreateSession(t, cl, &sin)
		if err != nil {
			return
		}
	}

	in := service.ListSessionsByConferenceRequest{
		ConferenceID: conference.ID,
	}
	res, err := cl.ListSessionsByConference(&in)
	if !assert.NoError(t, err, "ListSessionsByConference should succeed") {
		return
	}
	if !assert.NoError(t, validator.HTTPListSessionsByConferenceResponse.Validate(res), "Validation should succeed") {
		return
	}

	if !assert.Len(t, res, 10, "There should be 10 sessions") {
		return
	}
}

func TestListVenues(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)
	in := service.ListVenueRequest{}
	res, err := cl.ListVenues(&in)
	if !assert.NoError(t, err, "ListVenues should succeed") {
		return
	}
	if !assert.NoError(t, validator.HTTPListVenuesResponse.Validate(res), "Validation should succeed") {
		return
	}
}
