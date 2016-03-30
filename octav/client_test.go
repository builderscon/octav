package octav_test

import (
	"fmt"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"

	"github.com/builderscon/octav/octav"
	"github.com/builderscon/octav/octav/client"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	"github.com/builderscon/octav/octav/validator"
	"github.com/stretchr/testify/assert"
)

func bigsight() *model.CreateVenueRequest {
	lf := tools.LocalizedFields{}
	lf.Set("ja", "name", `東京ビッグサイト`)
	lf.Set("ja", "address", `〒135-0063 東京都江東区有明３丁目１０−１`)

	r := model.CreateVenueRequest{}
	r.L10N = lf
	r.Name.Set("Tokyo Bigsight")
	r.Address.Set("Ariake 3-10-1, Koto-ku, Tokyo")
	r.Longitude.Set(35.6320326)
	r.Latitude.Set(139.7976891)

	return &r
}

func intlConferenceRoom(venueID string) *model.CreateRoomRequest {
	lf := tools.LocalizedFields{}
	lf.Set("ja", "name", `国際会議場`)

	r := model.CreateRoomRequest{}
	r.L10N = lf
	r.Capacity.Set(1000)
	r.Name.Set("International Conference Hall")
	r.VenueID.Set(venueID)

	return &r
}

func testCreateRoom(t *testing.T, cl *client.Client, r *model.CreateRoomRequest) (*model.Room, error) {
	res, err := cl.CreateRoom(r)
	if !assert.NoError(t, err, "CreateRoom should succeed") {
		return nil, err
	}
	return res, nil
}

func testCreateVenue(t *testing.T, cl *client.Client, v *model.CreateVenueRequest) (*model.Venue, error) {
	res, err := cl.CreateVenue(v)
	if !assert.NoError(t, err, "CreateVenue should succeed") {
		return nil, err
	}
	return res, nil
}

func testLookupRoom(t *testing.T, cl *client.Client, id, lang string) (*model.Room, error) {
	r := &model.LookupRoomRequest{ID: id}
	if lang != "" {
		r.Lang.Set(lang)
	}
	venue, err := cl.LookupRoom(r)
	if !assert.NoError(t, err, "LookupRoom succeeds") {
		return nil, err
	}
	return venue, nil
}

func testCreateUser(t *testing.T, cl *client.Client, in *model.CreateUserRequest) (*model.User, error) {
	res, err := cl.CreateUser(in)
	if !assert.NoError(t, err, "CreateUser should succeed") {
		return nil, err
	}
	return res, nil
}

func testLookupUser(t *testing.T, cl *client.Client, id, lang string) (*model.User, error) {
	r := &model.LookupUserRequest{ID: id}
	if lang != "" {
		r.Lang.Set(lang)
	}
	user, err := cl.LookupUser(r)
	if !assert.NoError(t, err, "LookupUser succeeds") {
		return nil, err
	}
	return user, nil
}

func testDeleteUser(t *testing.T, cl *client.Client, id string) error {
	err := cl.DeleteUser(&model.DeleteUserRequest{ID: id})
	if !assert.NoError(t, err, "DeleteUser should succeed") {
		return err
	}
	return nil
}

func testLookupVenue(t *testing.T, cl *client.Client, id, lang string) (*model.Venue, error) {
	r := &model.LookupVenueRequest{ID: id}
	if lang != "" {
		r.Lang.Set(lang)
	}
	venue, err := cl.LookupVenue(r)
	if !assert.NoError(t, err, "LookupVenue succeeds") {
		return nil, err
	}
	return venue, nil
}

func testUpdateRoom(t *testing.T, cl *client.Client, in *model.UpdateRoomRequest) error {
	err := cl.UpdateRoom(in)
	if !assert.NoError(t, err, "UpdateRoom succeeds") {
		return err
	}
	return nil
}

func testDeleteRoom(t *testing.T, cl *client.Client, id string) error {
	err := cl.DeleteRoom(&model.DeleteRoomRequest{ID: id})
	if !assert.NoError(t, err, "DeleteRoom should be successful") {
		return err
	}
	return err
}

func testUpdateVenue(t *testing.T, cl *client.Client, in *model.UpdateVenueRequest) error {
	err := cl.UpdateVenue(in)
	if !assert.NoError(t, err, "UpdateVenue succeeds") {
		return err
	}
	return nil
}

func testDeleteVenue(t *testing.T, cl *client.Client, id string) error {
	err := cl.DeleteVenue(&model.DeleteVenueRequest{ID: id})
	if !assert.NoError(t, err, "DeleteVenue should be successful") {
		return err
	}
	return err
}

func yapcasia(userID string) *model.CreateConferenceRequest {
	return &model.CreateConferenceRequest{
		Title:  "YAPC::Asia Tokyo",
		Slug:   "yapcasia",
		UserID: userID,
	}
}

func testCreateConference(t *testing.T, cl *client.Client, in *model.CreateConferenceRequest) (*model.Conference, error) {
	res, err := cl.CreateConference(in)
	if !assert.NoError(t, err, "CreateConference should succeed") {
		return nil, err
	}
	return res, nil
}

func testLookupConference(t *testing.T, cl *client.Client, id, lang string) (*model.Conference, error) {
	r := &model.LookupConferenceRequest{ID: id}
	if lang != "" {
		r.Lang.Set(lang)
	}
	conference, err := cl.LookupConference(r)
	if !assert.NoError(t, err, "LookupConference succeeds") {
		return nil, err
	}
	return conference, nil
}

func testUpdateConference(t *testing.T, cl *client.Client, in *model.UpdateConferenceRequest) error {
	err := cl.UpdateConference(in)
	if !assert.NoError(t, err, "UpdateConference succeeds") {
		return err
	}
	return nil
}

func testDeleteConference(t *testing.T, cl *client.Client, id string) error {
	err := cl.DeleteConference(&model.DeleteConferenceRequest{ID: id})
	if !assert.NoError(t, err, "DeleteConference should be successful") {
		return err
	}
	return err
}

func TestConferenceCRUD(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)

	user, err := testCreateUser(t, cl, johndoe())
	if err != nil {
		return
	}
	defer testDeleteUser(t, cl, user.ID)

	res, err := testCreateConference(t, cl, yapcasia(user.ID))
	if err != nil {
		return
	}
	defer testDeleteConference(t, cl, res.ID)

	if !assert.NoError(t, validator.HTTPCreateConferenceResponse.Validate(res), "Validation should succeed") {
		return
	}

	res2, err := testLookupConference(t, cl, res.ID, "")
	if err != nil {
		return
	}

	// The result from LookupConference contains the administrator field
	// Remove that (and make sure it's populated), then do the comparison
	admins := res2.Administrators
	res2.Administrators = model.UserList(nil)
	if !assert.Equal(t, res2, res, "LookupConference is the same as the conference created") {
		return
	}

	if !assert.Len(t, admins, 1, "There should be 1 administrator") {
		return
	}

	in := model.UpdateConferenceRequest{ID: res.ID}
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

	in := model.UpdateRoomRequest{ID: res.ID}
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

func testCreateSession(t *testing.T, cl *client.Client, in *model.CreateSessionRequest) (*model.Session, error) {
	res, err := cl.CreateSession(in)
	if !assert.NoError(t, err, "CreateSession should succeed") {
		return nil, err
	}
	return res, nil
}

func testLookupSession(t *testing.T, cl *client.Client, id, lang string) (*model.Session, error) {
	r := &model.LookupSessionRequest{ID: id}
	if lang != "" {
		r.Lang.Set(lang)
	}
	v, err := cl.LookupSession(r)
	if !assert.NoError(t, err, "LookupSession succeeds") {
		return nil, err
	}
	return v, nil
}

func testUpdateSession(t *testing.T, cl *client.Client, in *model.UpdateSessionRequest) error {
	err := cl.UpdateSession(in)
	if !assert.NoError(t, err, "UpdateSession succeeds") {
		return err
	}
	return nil
}

func testDeleteSession(t *testing.T, cl *client.Client, id string) error {
	err := cl.DeleteSession(&model.DeleteSessionRequest{ID: id})
	if !assert.NoError(t, err, "DeleteSession should be successful") {
		return err
	}
	return err
}

func bconsession(cid, uid string) *model.CreateSessionRequest {
	in := model.CreateSessionRequest{}
	in.ConferenceID.Set(cid)
	in.SpeakerID.Set(uid)
	in.Title.Set("How To Write A Conference Backend")
	in.Duration.Set(60)
	in.Abstract.Set("Use lots of reflection and generate lots of code")
	return &in
}

func TestSessionCRUD(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)

	user, err := testCreateUser(t, cl, johndoe())
	if err != nil {
		return
	}
	defer testDeleteUser(t, cl, user.ID)

	conference, err := testCreateConference(t, cl, yapcasia(user.ID))
	if err != nil {
		return
	}

	res, err := testCreateSession(t, cl, bconsession(conference.ID, user.ID))
	if err != nil {
		return
	}

	if !assert.NoError(t, validator.HTTPCreateSessionResponse.Validate(res), "Validation should succeed") {
		return
	}

	res2, err := testLookupSession(t, cl, res.ID, "")
	if err != nil {
		return
	}

	if !assert.Equal(t, res2, res, "LookupSession is the same as the room created") {
		return
	}

	in := model.UpdateSessionRequest{ID: res.ID}
	in.L10N.Set("ja", "title", "カンファレンス用ソフトウェアの作り方")
	if err := testUpdateSession(t, cl, &in); err != nil {
		return
	}

	res3, err := testLookupSession(t, cl, res.ID, "ja")
	if err != nil {
		return
	}

	if !assert.Equal(t, "カンファレンス用ソフトウェアの作り方", res3.Title, "Session.title#ja is the same as the conference updated") {
		return
	}

	if err := testDeleteSession(t, cl, res.ID); err != nil {
		return
	}

	if err := testDeleteConference(t, cl, conference.ID); err != nil {
		return
	}
}

var ghidL = sync.Mutex{}
var ghid = 0

func johndoe() *model.CreateUserRequest {
	ghidL.Lock()
	defer ghidL.Unlock()

	ghid++
	lf := tools.LocalizedFields{}
	lf.Set("ja", "first_name", "ジョン")
	lf.Set("ja", "last_name", "ドー")

	r := model.CreateUserRequest{
		Nickname: tools.UUID(),
	}
	r.AuthVia.Set("github")
	r.AuthUserID.Set(strconv.Itoa(ghid))
	r.FirstName.Set("John")
	r.LastName.Set("Doe")
	r.Email.Set("john.doe@example.com")
	r.TshirtSize.Set("XL")
	r.L10N = lf
	return &r
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

	res2, err := testLookupUser(t, cl, res.ID, "")
	if err != nil {
		return
	}

	if !assert.Equal(t, res2, res, "LookupUser is the same as the user created") {
		return
	}

	res3, err := testLookupUser(t, cl, res.ID, "ja")
	if err != nil {
		return
	}

	if !assert.Equal(t, "ジョン", res3.FirstName, "User.first_name#ja is localized") {
		return
	}

	if !assert.Equal(t, "ドー", res3.LastName, "User.last_name#ja is localized") {
		return
	}

	if err := testDeleteUser(t, cl, res.ID); err != nil {
		return
	}
}

func TestVenueCRUD(t *testing.T) {
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

	res2, err := testLookupVenue(t, cl, res.ID, "")
	if err != nil {
		return
	}

	if !assert.Equal(t, res2, res, "LookupVenue is the same as the venue created") {
		return
	}

	in := model.UpdateVenueRequest{ID: res.ID}
	in.L10N.Set("ja", "name", "東京ビッグサイト")
	if err := testUpdateVenue(t, cl, &in); err != nil {
		return
	}

	res3, err := testLookupVenue(t, cl, res.ID, "ja")
	if err != nil {
		return
	}

	if !assert.Equal(t, "東京ビッグサイト", res3.Name, "Venue.name#ja is the same as the object updated") {
		return
	}

	if err := testDeleteVenue(t, cl, res.ID); err != nil {
		return
	}
}

func TestDeleteConferenceDates(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)

	user, err := testCreateUser(t, cl, johndoe())
	if err != nil {
		return
	}
	defer testDeleteUser(t, cl, user.ID)

	conf, err := testCreateConference(t, cl, &model.CreateConferenceRequest{
		UserID: user.ID,
	})
	if err != nil {
		return
	}
	defer testDeleteConference(t, cl, conf.ID)

	err = cl.AddConferenceDates(&model.AddConferenceDatesRequest{
		ConferenceID: conf.ID,
		Dates: []model.ConferenceDate{
			model.ConferenceDate{
				Date:  model.NewDate(2016, 3, 22),
				Open:  model.NewWallClock(10, 0),
				Close: model.NewWallClock(19, 0),
			},
		},
	})
	if !assert.NoError(t, err, "AddConferenceDates works") {
		return
	}

	err = cl.DeleteConferenceDates(&model.DeleteConferenceDatesRequest{
		ConferenceID: conf.ID,
		Dates:        []model.Date{model.NewDate(2016, 3, 22)},
	})
	if !assert.NoError(t, err, "DeleteConferenceDates works") {
		return
	}

	conf2, err := testLookupConference(t, cl, conf.ID, "")
	if err != nil {
		return
	}

	if !assert.Len(t, conf2.Dates, 0, "There should be no dates set") {
		return
	}
}

func TestConferenceAdmins(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)

	user, err := testCreateUser(t, cl, johndoe())
	if err != nil {
		return
	}
	defer testDeleteUser(t, cl, user.ID)

	conf, err := testCreateConference(t, cl, &model.CreateConferenceRequest{
		UserID: user.ID,
	})
	if err != nil {
		return
	}
	defer testDeleteConference(t, cl, conf.ID)

	for i := 0; i < 9; i++ {
		extraAdmin, err := testCreateUser(t, cl, johndoe())
		if err != nil {
			return
		}
		defer testDeleteUser(t, cl, extraAdmin.ID)

		err = cl.AddConferenceAdmin(&model.AddConferenceAdminRequest{
			ConferenceID: conf.ID,
			UserID:       extraAdmin.ID,
		})
		if !assert.NoError(t, err, "AddConferenceAdmin should succeed") {
			return
		}
	}

	conf2, err := testLookupConference(t, cl, conf.ID, "")
	if err != nil {
		return
	}

	if !assert.Len(t, conf2.Administrators, 10, "There should be 10 admins") {
		return
	}

	for _, admin := range conf2.Administrators {
		err = cl.DeleteConferenceAdmin(&model.DeleteConferenceAdminRequest{
			ConferenceID: conf.ID,
			UserID:       admin.ID,
		})
		if !assert.NoError(t, err, "DeleteConferenceAdmin should succeed") {
			return
		}
	}

	conf3, err := testLookupConference(t, cl, conf.ID, "")
	if err != nil {
		return
	}

	if !assert.Len(t, conf3.Administrators, 0, "There should be 0 admins") {
		return
	}

}

func TestListConference(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)

	user, err := testCreateUser(t, cl, johndoe())
	if err != nil {
		return
	}
	defer testDeleteUser(t, cl, user.ID)

	for i := 0; i < 10; i++ {
		conf, err := testCreateConference(t, cl, &model.CreateConferenceRequest{
			UserID: user.ID,
		})
		if err != nil {
			return
		}

		err = cl.AddConferenceDates(&model.AddConferenceDatesRequest{
			ConferenceID: conf.ID,
			Dates: []model.ConferenceDate{
				model.ConferenceDate{
					Date:  model.NewDate(2016, 3, 22),
					Open:  model.NewWallClock(10, 0),
					Close: model.NewWallClock(19, 0),
				},
			},
		})
		if !assert.NoError(t, err, "AddConferenceDates works") {
			return
		}

		defer testDeleteConference(t, cl, conf.ID)
	}

	in := model.ListConferenceRequest{}
	in.RangeStart.Set("2016-03-21")
	in.RangeEnd.Set("2016-03-23")
	res, err := cl.ListConference(&in)
	if !assert.NoError(t, err, "ListConference should succeed") {
		return
	}

	if !assert.NoError(t, validator.HTTPListConferenceResponse.Validate(res), "Validation should succeed") {
		return
	}

	if !assert.Len(t, res, 10, "ListConference returns 10 rooms") {
		return
	}
}

func TestListRoom(t *testing.T) {
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

	in := model.ListRoomRequest{
		VenueID: venue.ID,
	}
	res, err := cl.ListRoom(&in)
	if !assert.NoError(t, err, "ListRoom should succeed") {
		return
	}

	if !assert.NoError(t, validator.HTTPListRoomResponse.Validate(res), "Validation should succeed") {
		return
	}

	if !assert.Len(t, res, 1, "ListRoom returns 1 rooms") {
		return
	}
}

func TestListSessionByConference(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)

	user, err := testCreateUser(t, cl, johndoe())
	if err != nil {
		return
	}
	defer testDeleteUser(t, cl, user.ID)
	conference, err := testCreateConference(t, cl, yapcasia(user.ID))
	if err != nil {
		return
	}

	for i := 0; i < 10; i++ {
		sin := model.CreateSessionRequest{}
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

	in := model.ListSessionByConferenceRequest{
		ConferenceID: conference.ID,
	}
	res, err := cl.ListSessionByConference(&in)
	if !assert.NoError(t, err, "ListSessionByConference should succeed") {
		return
	}
	if !assert.NoError(t, validator.HTTPListSessionByConferenceResponse.Validate(res), "Validation should succeed") {
		return
	}

	if !assert.Len(t, res, 10, "There should be 10 sessions") {
		return
	}
}

func TestListVenue(t *testing.T) {
	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	cl := client.New(ts.URL)
	in := model.ListVenueRequest{}
	res, err := cl.ListVenue(&in)
	if !assert.NoError(t, err, "ListVenue should succeed") {
		return
	}
	if !assert.NoError(t, validator.HTTPListVenueResponse.Validate(res), "Validation should succeed") {
		return
	}
}
