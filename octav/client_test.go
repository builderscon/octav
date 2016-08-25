package octav_test

import (
	"fmt"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/builderscon/octav/octav"
	"github.com/builderscon/octav/octav/client"
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	"github.com/builderscon/octav/octav/validator"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type TestCtx struct {
	*testing.T
	APIClient  *db.Client
	Superuser  *db.User
	HTTPClient *client.Client
}

func NewTestCtx(t *testing.T) (*TestCtx, error) {
	client := db.Client{
		EID:    tools.UUID(),
		Secret: tools.UUID(), // Todo
		Name:   "Test Client",
	}
	u := newuser()
	superuser := db.User{
		AuthUserID: u.AuthUserID,
		AuthVia:    u.AuthVia,
		EID:        tools.UUID(),
		IsAdmin:    true,
		Nickname:   "root",
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "failed to start DB transaction")
	}
	defer tx.AutoRollback()

	if err = client.Create(tx); err != nil {
		return nil, err
	}

	if err = superuser.Create(tx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "failed to commit changes to DB")
	}

	ctx := &TestCtx{
		T:         t,
		APIClient: &client,
		Superuser: &superuser,
	}

	return ctx, nil
}

func (ctx *TestCtx) Subtest(name string, cb func(*TestCtx)) {
	ctx.T.Run(name, func(t *testing.T) {
		localctx := &TestCtx{
			T:          t,
			APIClient:  ctx.APIClient,
			Superuser:  ctx.Superuser,
			HTTPClient: ctx.HTTPClient,
		}
		cb(localctx)
	})
}

func (ctx *TestCtx) Close() error {
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to start DB transaction")
	}
	defer tx.AutoRollback()

	if cl := ctx.APIClient; cl != nil {
		if err := cl.Delete(tx); err != nil {
			return errors.Wrap(err, "failed to delete client")
		}
	}

	if u := ctx.Superuser; u != nil {
		if err := u.Delete(tx); err != nil {
			return errors.Wrap(err, "failed to delete superuser")
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit changes to DB")
	}
	return nil
}

func (ctx *TestCtx) SetAPIServer(ts *httptest.Server) {
	ctx.HTTPClient = client.New(ts.URL)
	ctx.HTTPClient.BasicAuth.Username = ctx.APIClient.EID
	ctx.HTTPClient.BasicAuth.Password = ctx.APIClient.Secret
}

func bigsight(userID string) *model.CreateVenueRequest {
	lf := model.LocalizedFields{}
	lf.Set("ja", "name", `東京ビッグサイト`)
	lf.Set("ja", "address", `〒135-0063 東京都江東区有明３丁目１０−１`)

	r := model.CreateVenueRequest{}
	r.L10N = lf
	r.Name.Set("Tokyo Bigsight")
	r.Address.Set("Ariake 3-10-1, Koto-ku, Tokyo")
	r.Longitude.Set(35.6320326)
	r.Latitude.Set(139.7976891)
	r.UserID = userID

	return &r
}

func intlConferenceRoom(venueID, userID string) *model.CreateRoomRequest {
	lf := model.LocalizedFields{}
	lf.Set("ja", "name", `国際会議場`)

	r := model.CreateRoomRequest{}
	r.L10N = lf
	r.Capacity.Set(1000)
	r.Name.Set("International Conference Hall")
	r.VenueID.Set(venueID)
	r.UserID = userID

	return &r
}

func testAddConferenceAdmin(ctx *TestCtx, confID, adminID, userID string) error {
	err := ctx.HTTPClient.AddConferenceAdmin(&model.AddConferenceAdminRequest{
		ConferenceID: confID,
		AdminID:      adminID,
		UserID:       userID,
	})
	if !assert.NoError(ctx.T, err, "AddConferenceAdmin should succeed") {
		return err
	}
	return nil
}

func testDeleteConferenceAdmin(ctx *TestCtx, confID, adminID, userID string) error {
	err := ctx.HTTPClient.DeleteConferenceAdmin(&model.DeleteConferenceAdminRequest{
		ConferenceID: confID,
		AdminID:      adminID,
		UserID:       userID,
	})
	if !assert.NoError(ctx.T, err, "DeleteConferenceAdmin should succeed") {
		return err
	}
	return nil
}

func testAddConferenceVenue(ctx *TestCtx, confID, venueID, userID string) error {
	req := model.AddConferenceVenueRequest{
		ConferenceID: confID,
		VenueID:      venueID,
		UserID:       userID,
	}
	err := ctx.HTTPClient.AddConferenceVenue(&req)
	if !assert.NoError(ctx.T, err, "AddConferenceVenue should succeed") {
		return err
	}
	return nil
}

func testDeleteConferenceVenue(ctx *TestCtx, confID, venueID, userID string) error {
	req := model.DeleteConferenceVenueRequest{
		ConferenceID: confID,
		VenueID:      venueID,
		UserID:       userID,
	}
	err := ctx.HTTPClient.DeleteConferenceVenue(&req)
	if !assert.NoError(ctx.T, err, "DeleteConferenceVenue should succeed") {
		return err
	}
	return nil
}

func testLookupRoom(ctx *TestCtx, id, lang string) (*model.Room, error) {
	r := &model.LookupRoomRequest{ID: id}
	if lang != "" {
		r.Lang.Set(lang)
	}
	venue, err := ctx.HTTPClient.LookupRoom(r)
	if !assert.NoError(ctx.T, err, "LookupRoom succeeds") {
		return nil, err
	}
	return venue, nil
}

func testCreateUser(ctx *TestCtx, in *model.CreateUserRequest) (*model.User, error) {
	res, err := ctx.HTTPClient.CreateUser(in)
	if !assert.NoError(ctx.T, err, "CreateUser should succeed") {
		return nil, err
	}
	return res, nil
}

func testLookupUser(ctx *TestCtx, id, lang string) (*model.User, error) {
	r := &model.LookupUserRequest{ID: id}
	if lang != "" {
		r.Lang.Set(lang)
	}
	user, err := ctx.HTTPClient.LookupUser(r)
	if !assert.NoError(ctx.T, err, "LookupUser succeeds") {
		return nil, err
	}
	return user, nil
}

func testDeleteUser(ctx *TestCtx, targetID, userID string) error {
	err := ctx.HTTPClient.DeleteUser(&model.DeleteUserRequest{ID: targetID, UserID: userID})
	if !assert.NoError(ctx.T, err, "DeleteUser should succeed") {
		return err
	}
	return nil
}

func testLookupVenue(ctx *TestCtx, id, lang string) (*model.Venue, error) {
	r := &model.LookupVenueRequest{ID: id}
	if lang != "" {
		r.Lang.Set(lang)
	}
	venue, err := ctx.HTTPClient.LookupVenue(r)
	if !assert.NoError(ctx.T, err, "LookupVenue succeeds") {
		return nil, err
	}
	return venue, nil
}

func testUpdateRoom(ctx *TestCtx, in *model.UpdateRoomRequest) error {
	err := ctx.HTTPClient.UpdateRoom(in)
	if !assert.NoError(ctx.T, err, "UpdateRoom succeeds") {
		return err
	}
	return nil
}

func testDeleteRoom(ctx *TestCtx, roomID, userID string) error {
	err := ctx.HTTPClient.DeleteRoom(&model.DeleteRoomRequest{ID: roomID, UserID: userID})
	if !assert.NoError(ctx.T, err, "DeleteRoom should be successful") {
		return err
	}
	return err
}

func testUpdateVenue(ctx *TestCtx, in *model.UpdateVenueRequest) error {
	err := ctx.HTTPClient.UpdateVenue(in)
	if !assert.NoError(ctx.T, err, "UpdateVenue succeeds") {
		return err
	}
	return nil
}

func testDeleteVenue(ctx *TestCtx, venueID, userID string) error {
	err := ctx.HTTPClient.DeleteVenue(&model.DeleteVenueRequest{ID: venueID, UserID: userID})
	if !assert.NoError(ctx.T, err, "DeleteVenue should be successful") {
		return err
	}
	return err
}

func yapcasia(uid string) *model.CreateConferenceSeriesRequest {
	return &model.CreateConferenceSeriesRequest{
		UserID: uid,
		Slug:   "yapcasia",
	}
}

func yapcasiaTokyo(seriesID, userID string) *model.CreateConferenceRequest {
	r := &model.CreateConferenceRequest{
		Title:       "YAPC::Asia Tokyo",
		SeriesID:    seriesID,
		Slug:        "2015",
		UserID:      userID,
	}
	r.L10N.Set("ja", "description", "最後のYAPC::Asia Tokyo")
	r.Description.Set("The last YAPC::Asia Tokyo")
	return r
}

func testCreateConferenceSeries(ctx *TestCtx, in *model.CreateConferenceSeriesRequest) (*model.ObjectID, error) {
	res, err := ctx.HTTPClient.CreateConferenceSeries(in)
	if !assert.NoError(ctx.T, err, "CreateConferenceSeries should succeed") {
		return nil, err
	}
	return res, nil
}

func testDeleteConferenceSeries(ctx *TestCtx, id, userID string) error {
	err := ctx.HTTPClient.DeleteConferenceSeries(&model.DeleteConferenceSeriesRequest{ID: id, UserID: userID})
	if !assert.NoError(ctx.T, err, "DeleteConferenceSeries should be successful") {
		return err
	}
	return err
}

func testAddConferenceSeriesAdmin(ctx *TestCtx, id, adminID, userID string) error {
	err := ctx.HTTPClient.AddConferenceSeriesAdmin(&model.AddConferenceSeriesAdminRequest{SeriesID: id, AdminID: adminID, UserID: userID})
	if !assert.NoError(ctx.T, err, "AddConferenceSeriesAdmin should be successful") {
		return err
	}
	return err
}

func testLookupConference(ctx *TestCtx, id, lang string) (*model.Conference, error) {
	r := &model.LookupConferenceRequest{ID: id}
	if lang != "" {
		r.Lang.Set(lang)
	}
	conference, err := ctx.HTTPClient.LookupConference(r)
	if !assert.NoError(ctx.T, err, "LookupConference succeeds") {
		return nil, err
	}
	return conference, nil
}

func testUpdateConference(ctx *TestCtx, in *model.UpdateConferenceRequest) error {
	err := ctx.HTTPClient.UpdateConference(in, nil)
	if !assert.NoError(ctx.T, err, "UpdateConference succeeds") {
		return err
	}
	return nil
}

func testMakeConferencePublic(ctx *TestCtx, conferenceID, userID string) error {
	r := model.UpdateConferenceRequest{
		ID:     conferenceID,
		UserID: userID,
	}
	r.Status.Set("public")
	return testUpdateConference(ctx, &r)
}

func testDeleteConference(ctx *TestCtx, id string) error {
	err := ctx.HTTPClient.DeleteConference(&model.DeleteConferenceRequest{ID: id})
	if !assert.NoError(ctx.T, err, "DeleteConference should be successful") {
		return err
	}
	return err
}

func larrywall(confID, userID string) *model.AddFeaturedSpeakerRequest {
	r := &model.AddFeaturedSpeakerRequest{
		ConferenceID: confID,
		DisplayName:  `Larry Wall (TimToady)`,
		Description:  `Larry Wall is a computer programmer and author, most widely known as the creator of the Perl programming language.`,
		UserID:       userID,
	}
	r.AvatarURL.Set(`https://upload.wikimedia.org/wikipedia/commons/b/b3/Larry_Wall_YAPC_2007.jpg`)
	return r
}

func testCreateFeaturedSpeaker(ctx *TestCtx, in *model.AddFeaturedSpeakerRequest) (*model.FeaturedSpeaker, error) {
	res, err := ctx.HTTPClient.AddFeaturedSpeaker(in)
	if !assert.NoError(ctx.T, err, "CreateFeaturedSpeaker should succeed") {
		return nil, err
	}
	return res, nil
}

func testDeleteFeaturedSpeaker(ctx *TestCtx, id, userID string) error {
	err := ctx.HTTPClient.DeleteFeaturedSpeaker(&model.DeleteFeaturedSpeakerRequest{
		ID:     id,
		UserID: userID,
	})
	if !assert.NoError(ctx.T, err, "DeleteFeaturedSpeaker should succeed") {
		return err
	}
	return nil
}

func buildersconinc(confID, userID string) *model.AddSponsorRequest {
	r := &model.AddSponsorRequest{
		ConferenceID: confID,
		Name:         "builderscon",
		URL:          "http://builderscon.io",
		GroupName:    "tier-1",
		UserID:       userID,
	}
	return r
}

func TestConferenceCRUD(t *testing.T) {
	ctx, err := NewTestCtx(t)
	if !assert.NoError(t, err, "failed to create test ctx") {
		return
	}
	defer ctx.Close()

	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	ctx.SetAPIServer(ts)

	user, err := testCreateUser(ctx, johndoe())
	if err != nil {
		return
	}
	defer testDeleteUser(ctx, user.ID, ctx.Superuser.EID)

	series, err := testCreateConferenceSeries(ctx, yapcasia(ctx.Superuser.EID))
	if err != nil {
		return
	}
	defer testDeleteConferenceSeries(ctx, series.ID, ctx.Superuser.EID)

	if err := testAddConferenceSeriesAdmin(ctx, series.ID, user.ID, ctx.Superuser.EID); err != nil {
		return
	}

	res, err := testCreateConferencePass(ctx, yapcasiaTokyo(series.ID, user.ID))
	if err != nil {
		return
	}
	defer testDeleteConference(ctx, res.ID)

	if !assert.NoError(ctx.T, validator.HTTPCreateConferenceResponse.Validate(res), "Validation should succeed") {
		return
	}

	conf1, err := testLookupConference(ctx, res.ID, "")
	if err != nil {
		return
	}

	if !assert.NotEmpty(ctx.T, conf1.Description, "Description should not be empty") {
		return
	}

	if !assert.Len(ctx.T, conf1.Administrators, 1, "There should be 1 administrator") {
		return
	}

	conf2, err := testLookupConference(ctx, res.ID, "ja")
	if err != nil {
		return
	}

	if !assert.NotEqual(ctx.T, conf1.Description, conf2.Description, "description should localized") {
		return
	}

	conf3, err := testLookupConference(ctx, res.ID, "ja")
	if err != nil {
		return
	}

	if !assert.Equal(ctx.T, conf2, conf3, "LookupConference is the same") {
		return
	}

	in := model.UpdateConferenceRequest{ID: res.ID, UserID: user.ID}
	in.SubTitle.Set("Big Bang!")
	in.L10N.Set("ja", "title", "ヤップシー エイジア")
	in.L10N.Set("ja", "cfp_lead_text", "ばっちこい！")
	in.L10N.Set("ja", "cfp_pre_submit_instructions", "事前にこれを読んでね")
	in.L10N.Set("ja", "cfp_post_submit_instructions", "応募したらこれを読んでね")
	if err := testUpdateConference(ctx, &in); err != nil {
		return
	}

	conf4, err := testLookupConference(ctx, res.ID, "ja")
	if err != nil {
		return
	}
	if !assert.Equal(ctx.T, conf4.SubTitle, "Big Bang!", "Conference.SubTitle is the same as the conference updated") {
		return
	}

	if !assert.Equal(ctx.T, "ヤップシー エイジア", conf4.Title, "Conference.title#ja is the same as the conference updated") {
		return
	}

	if !assert.Equal(ctx.T, "ばっちこい！", conf4.CFPLeadText, "Conference.cfp_lead_text#ja is the same as the conference updated") {
		return
	}

	if !assert.Equal(ctx.T, "事前にこれを読んでね", conf4.CFPPreSubmitInstructions, "Conference.cfp_pre_submit_instructions#ja is the same as the conference updated") {
		return
	}

	if !assert.Equal(ctx.T, "応募したらこれを読んでね", conf4.CFPPostSubmitInstructions, "Conference.cfp_post_submit_instructions#ja is the same as the conference updated") {
		return
	}

	venue, err := testCreateVenuePass(ctx, bigsight(user.ID))
	if err != nil {
		return
	}

	if err := testAddConferenceVenue(ctx, res.ID, venue.ID, user.ID); err != nil {
		return
	}
	defer testDeleteConferenceVenue(ctx, res.ID, venue.ID, user.ID)

	// add a featured speaker
	fs, err := testCreateFeaturedSpeaker(ctx, larrywall(res.ID, user.ID))
	if err != nil {
		return
	}
	defer testDeleteFeaturedSpeaker(ctx, fs.ID, user.ID)

	sp, err := testCreateSponsor(ctx, buildersconinc(res.ID, user.ID))
	if err != nil {
		return
	}
	defer testDeleteSponsor(ctx, sp.ID, user.ID)
}

func TestRoomCRUD(t *testing.T) {
	ctx, err := NewTestCtx(t)
	if !assert.NoError(t, err, "failed to create test ctx") {
		return
	}
	defer ctx.Close()

	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	ctx.SetAPIServer(ts)

	venue, err := testCreateVenuePass(ctx, bigsight(ctx.Superuser.EID))
	if err != nil {
		return
	}

	res, err := testCreateRoomPass(ctx, intlConferenceRoom(venue.ID, ctx.Superuser.EID))
	if err != nil {
		return
	}

	if !assert.NotEmpty(ctx.T, res.ID, "Returned structure has ID") {
		return
	}

	if !assert.NoError(ctx.T, validator.HTTPCreateRoomResponse.Validate(res), "Validation should succeed") {
		return
	}

	room1, err := testLookupRoom(ctx, res.ID, "ja")
	if err != nil {
		return
	}

	room2, err := testLookupRoom(ctx, res.ID, "ja")
	if err != nil {
		return
	}

	if !assert.Equal(ctx.T, room2, room1, "LookupRoom is the same as the room created") {
		return
	}

	in := model.UpdateRoomRequest{ID: res.ID, UserID: ctx.Superuser.EID}
	in.L10N.Set("ja", "name", "国際会議場")
	if err := testUpdateRoom(ctx, &in); err != nil {
		return
	}

	room3, err := testLookupRoom(ctx, res.ID, "ja")
	if err != nil {
		return
	}

	if !assert.Equal(ctx.T, "国際会議場", room3.Name, "Room.name#ja is the same as the conference updated") {
		return
	}

	if err := testDeleteRoom(ctx, res.ID, ctx.Superuser.EID); err != nil {
		return
	}

	if err := testDeleteVenue(ctx, venue.ID, ctx.Superuser.EID); err != nil {
		return
	}
}

func bconsession(cid, speakerID, userID, sessionTypeID string) *model.CreateSessionRequest {
	in := model.CreateSessionRequest{}
	in.ConferenceID = cid
	in.SpeakerID.Set(speakerID)
	in.Title.Set("How To Write A Conference Backend")
	in.SessionTypeID = sessionTypeID
	in.Abstract.Set("Use lots of reflection and generate lots of code")
	in.UserID = userID
	return &in
}

func testStartSubmission(ctx *TestCtx, sessionTypeID, userID string, ref time.Time) error {
	r := &model.UpdateSessionTypeRequest{
		ID:     sessionTypeID,
		UserID: userID,
	}
	r.SubmissionStart.Set(ref.Add(-1 * 24 * time.Hour).Format(time.RFC3339))
	r.SubmissionEnd.Set(ref.Add(24 * time.Hour).Format(time.RFC3339))
	err := ctx.HTTPClient.UpdateSessionType(r)
	if !assert.NoError(ctx.T, err, "StartSessionSubmission should be successful") {
		return err
	}
	return err
}

func TestSessionCRUD(t *testing.T) {
	ctx, err := NewTestCtx(t)
	if !assert.NoError(t, err, "failed to create test ctx") {
		return
	}
	defer ctx.Close()

	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	ctx.SetAPIServer(ts)

	organizer, err := testCreateUser(ctx, johndoe())
	if err != nil {
		return
	}
	defer testDeleteUser(ctx, organizer.ID, ctx.Superuser.EID)

	user, err := testCreateUser(ctx, johndoe())
	if err != nil {
		return
	}
	defer testDeleteUser(ctx, user.ID, ctx.Superuser.EID)

	series, err := testCreateConferenceSeries(ctx, yapcasia(ctx.Superuser.EID))
	if err != nil {
		return
	}
	defer testDeleteConferenceSeries(ctx, series.ID, ctx.Superuser.EID)

	if err := testAddConferenceSeriesAdmin(ctx, series.ID, organizer.ID, ctx.Superuser.EID); err != nil {
		return
	}

	conference, err := testCreateConferencePass(ctx, yapcasiaTokyo(series.ID, organizer.ID))
	if err != nil {
		return
	}
	defer testDeleteConference(ctx, conference.ID)

	// Make sure the conference is public
	if err := testMakeConferencePublic(ctx, conference.ID, organizer.ID); err != nil {
		return
	}

	list, err := ctx.HTTPClient.ListSessionTypesByConference(&model.ListSessionTypesByConferenceRequest{
		ConferenceID: conference.ID,
	})
	if !assert.NoError(ctx.T, err) {
		return
	}
	if !assert.True(ctx.T, len(list) > 0) {
		return
	}

	// Make this one session type to NOT accept talks right now
	ctx.Subtest("Proposal submission should be rejected if out of range", func(ctx *TestCtx) {
		stype := list[1]
		// Set time at 1 month ago
		if err := testStartSubmission(ctx, stype.ID, organizer.ID, time.Now().Add(-1*24*30*time.Hour)); err != nil {
			return
		}
		_, err := testCreateSessionFail(ctx, bconsession(conference.ID, user.ID, user.ID, stype.ID))
		if err != nil {
			return
		}
	})

	stype := list[0]
	if err := testStartSubmission(ctx, stype.ID, organizer.ID, time.Now()); err != nil {
		return
	}

	res, err := testCreateSessionPass(ctx, bconsession(conference.ID, user.ID, user.ID, stype.ID))
	if err != nil {
		return
	}
	defer testDeleteSessionPass(ctx, res.ID, user.ID)

	if !assert.NoError(ctx.T, validator.HTTPCreateSessionResponse.Validate(res), "Validation should succeed") {
		return
	}

	session1, err := testLookupSession(ctx, res.ID, "ja")
	if err != nil {
		return
	}

	session2, err := testLookupSession(ctx, res.ID, "ja")
	if err != nil {
		return
	}
	if !assert.Equal(ctx.T, session2, session1, "LookupSession is the same as the room created") {
		return
	}

	if !assert.NotEmpty(ctx.T, session1.Speaker.Email, "email should NOT be empty for authenticated requests") {
		return
	}

	in := model.UpdateSessionRequest{ID: res.ID, UserID: user.ID}
	in.L10N.Set("ja", "title", "カンファレンス用ソフトウェアの作り方")
	if err := testUpdateSession(ctx, &in); err != nil {
		return
	}

	session3, err := testLookupSession(ctx, res.ID, "ja")
	if err != nil {
		return
	}

	if !assert.Equal(ctx.T, "カンファレンス用ソフトウェアの作り方", session3.Title, "Session.title#ja is the same as the conference updated") {
		return
	}

	ctx.Subtest("Proposals should not be deleted if accepted", func(ctx *TestCtx) {
		res, err := testCreateSessionPass(ctx, bconsession(conference.ID, user.ID, user.ID, stype.ID))
		if err != nil {
			return
		}

		in := model.UpdateSessionRequest{
			ID: res.ID,
			UserID: ctx.Superuser.EID,
		}
		in.Status.Set(model.StatusAccepted)
		if err := testUpdateSession(ctx, &in); err != nil {
			testDeleteSessionPass(ctx, res.ID, ctx.Superuser.EID)
			return
		}

		if err := testDeleteSessionFail(ctx, res.ID, user.ID); err != nil {
			return
		}

		if err := testDeleteSessionPass(ctx, res.ID, ctx.Superuser.EID); err != nil {
			return
		}
	})
}

var ghidL = sync.Mutex{}
var ghid = 0

func newuser() *model.CreateUserRequest {
	ghidL.Lock()
	defer ghidL.Unlock()

	ghid++
	r := model.CreateUserRequest{}

	r.AuthVia = "github"
	r.AuthUserID = strconv.Itoa(ghid)
	return &r
}

func johndoe() *model.CreateUserRequest {
	r := newuser()

	lf := model.LocalizedFields{}
	lf.Set("ja", "first_name", "ジョン")
	lf.Set("ja", "last_name", "ドー")

	r.Nickname = tools.UUID()
	r.AuthVia = "github"
	r.AuthUserID = tools.RandomString(32)
	r.FirstName.Set("John")
	r.LastName.Set("Doe")
	r.Email.Set("john.doe@example.com")
	r.TshirtSize.Set("XL")
	r.L10N = lf
	return r
}

func TestCreateUser(t *testing.T) {
	ctx, err := NewTestCtx(t)
	if !assert.NoError(t, err, "failed to create test ctx") {
		return
	}
	defer ctx.Close()

	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	ctx.SetAPIServer(ts)

	res, err := testCreateUser(ctx, johndoe())
	if err != nil {
		return
	}

	if !assert.NotEmpty(ctx.T, res.ID, "Returned structure has ID") {
		return
	}

	if !assert.NoError(ctx.T, validator.HTTPCreateUserResponse.Validate(res), "Validation should succeed") {
		return
	}

	res2, err := testLookupUser(ctx, res.ID, "")
	if err != nil {
		return
	}

	if !assert.Equal(ctx.T, res2, res, "LookupUser is the same as the user created") {
		return
	}

	res3, err := testLookupUser(ctx, res.ID, "ja")
	if err != nil {
		return
	}

	if !assert.Equal(ctx.T, "ジョン", res3.FirstName, "User.first_name#ja is localized") {
		return
	}

	if !assert.Equal(ctx.T, "ドー", res3.LastName, "User.last_name#ja is localized") {
		return
	}

	res4, err := ctx.HTTPClient.LookupUserByAuthUserID(&model.LookupUserByAuthUserIDRequest{
		AuthVia:    res.AuthVia,
		AuthUserID: res.AuthUserID,
	})
	if !assert.NoError(ctx.T, err, "LookupUserByAuthUserID should succeed") {
		return
	}

	if !assert.Equal(ctx.T, res4, res) {
		return
	}

	clientID, clientSecret := ctx.HTTPClient.BasicAuth.Username, ctx.HTTPClient.BasicAuth.Password
	ctx.HTTPClient.BasicAuth.Username = ""
	ctx.HTTPClient.BasicAuth.Password = ""
	res5, err := testLookupUser(ctx, res.ID, "")
	if err != nil {
		return
	}

	res.Email = ""
	res.TshirtSize = ""
	res.AuthVia = ""
	res.AuthUserID = ""
	if !assert.Equal(ctx.T, res5, res, "Lookup user should be the same (email and tshirt-size should be empty)") {
		return
	}

	ctx.HTTPClient.BasicAuth.Username = clientID
	ctx.HTTPClient.BasicAuth.Password = clientSecret
	if err := testDeleteUser(ctx, res.ID, ctx.Superuser.EID); err != nil {
		return
	}
}

func TestVenueCRUD(t *testing.T) {
	ctx, err := NewTestCtx(t)
	if !assert.NoError(t, err, "failed to create test ctx") {
		return
	}
	defer ctx.Close()

	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	ctx.SetAPIServer(ts)

	res, err := testCreateVenuePass(ctx, bigsight(ctx.Superuser.EID))
	if err != nil {
		return
	}

	if !assert.NotEmpty(ctx.T, res.ID, "Returned structure has ID") {
		return
	}

	if !assert.NoError(ctx.T, validator.HTTPCreateVenueResponse.Validate(res), "Validation should succeed") {
		return
	}

	room1, err := testLookupVenue(ctx, res.ID, "ja")
	if err != nil {
		return
	}

	room2, err := testLookupVenue(ctx, res.ID, "ja")
	if err != nil {
		return
	}

	if !assert.Equal(ctx.T, room2, room1, "LookupVenue is the same as the venue created") {
		return
	}

	in := model.UpdateVenueRequest{ID: res.ID, UserID: ctx.Superuser.EID}
	in.L10N.Set("ja", "name", "東京ビッグサイト")
	if err := testUpdateVenue(ctx, &in); err != nil {
		return
	}

	room3, err := testLookupVenue(ctx, res.ID, "ja")
	if err != nil {
		return
	}

	if !assert.Equal(ctx.T, "東京ビッグサイト", room3.Name, "Venue.name#ja is the same as the object updated") {
		return
	}

	if err := testDeleteVenue(ctx, res.ID, ctx.Superuser.EID); err != nil {
		return
	}
}

func TestDeleteConferenceDates(t *testing.T) {
	ctx, err := NewTestCtx(t)
	if !assert.NoError(t, err, "failed to create test ctx") {
		return
	}
	defer ctx.Close()

	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	ctx.SetAPIServer(ts)

	user, err := testCreateUser(ctx, johndoe())
	if err != nil {
		return
	}
	defer testDeleteUser(ctx, user.ID, ctx.Superuser.EID)

	series, err := testCreateConferenceSeries(ctx, &model.CreateConferenceSeriesRequest{
		UserID: ctx.Superuser.EID,
		Slug:   tools.RandomString(8),
	})
	if err != nil {
		return
	}
	defer testDeleteConferenceSeries(ctx, series.ID, ctx.Superuser.EID)

	if err := testAddConferenceSeriesAdmin(ctx, series.ID, user.ID, ctx.Superuser.EID); err != nil {
		return
	}

	conf, err := testCreateConferencePass(ctx, &model.CreateConferenceRequest{
		UserID:   user.ID,
		SeriesID: series.ID,
		Slug:     tools.RandomString(8),
	})
	if err != nil {
		return
	}
	defer testDeleteConference(ctx, conf.ID)

	err = ctx.HTTPClient.AddConferenceDates(&model.AddConferenceDatesRequest{
		ConferenceID: conf.ID,
		UserID:       user.ID,
		Dates: []model.ConferenceDate{
			model.ConferenceDate{
				Date:  model.NewDate(2016, 3, 22),
				Open:  model.NewWallClock(10, 0),
				Close: model.NewWallClock(19, 0),
			},
		},
	})
	if !assert.NoError(ctx.T, err, "AddConferenceDates works") {
		return
	}

	err = ctx.HTTPClient.DeleteConferenceDates(&model.DeleteConferenceDatesRequest{
		ConferenceID: conf.ID,
		Dates:        []model.Date{model.NewDate(2016, 3, 22)},
		UserID:       user.ID,
	})
	if !assert.NoError(ctx.T, err, "DeleteConferenceDates works") {
		return
	}

	conf2, err := testLookupConference(ctx, conf.ID, "")
	if err != nil {
		return
	}

	if !assert.Len(ctx.T, conf2.Dates, 0, "There should be no dates set") {
		return
	}
}

func TestConferenceAdmins(t *testing.T) {
	ctx, err := NewTestCtx(t)
	if !assert.NoError(t, err, "failed to create test ctx") {
		return
	}
	defer ctx.Close()

	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	ctx.SetAPIServer(ts)

	user, err := testCreateUser(ctx, johndoe())
	if err != nil {
		return
	}
	defer testDeleteUser(ctx, user.ID, ctx.Superuser.EID)

	series, err := testCreateConferenceSeries(ctx, &model.CreateConferenceSeriesRequest{
		UserID: ctx.Superuser.EID,
		Slug:   tools.RandomString(8),
	})
	if err != nil {
		return
	}
	defer testDeleteConferenceSeries(ctx, series.ID, ctx.Superuser.EID)

	if err := testAddConferenceSeriesAdmin(ctx, series.ID, user.ID, ctx.Superuser.EID); err != nil {
		return
	}

	conf, err := testCreateConferencePass(ctx, &model.CreateConferenceRequest{
		UserID:   user.ID,
		SeriesID: series.ID,
		Slug:     tools.RandomString(8),
	})
	if err != nil {
		return
	}
	defer testDeleteConference(ctx, conf.ID)

	for i := 0; i < 9; i++ {
		extraAdmin, err := testCreateUser(ctx, johndoe())
		if err != nil {
			return
		}
		defer testDeleteUser(ctx, extraAdmin.ID, ctx.Superuser.EID)

		if err := testAddConferenceAdmin(ctx, conf.ID, extraAdmin.ID, user.ID); err != nil {
			return
		}
	}

	conf2, err := testLookupConference(ctx, conf.ID, "")
	if err != nil {
		return
	}

	if !assert.Len(ctx.T, conf2.Administrators, 10, "There should be 10 admins") {
		return
	}

	for _, admin := range conf2.Administrators {
		if err := testDeleteConferenceAdmin(ctx, conf.ID, admin.ID, user.ID); err != nil {
			return
		}
	}

	conf3, err := testLookupConference(ctx, conf.ID, "")
	if err != nil {
		return
	}

	if !assert.Len(ctx.T, conf3.Administrators, 0, "There should be 0 admins") {
		return
	}
}

func TestListConference(t *testing.T) {
	ctx, err := NewTestCtx(t)
	if !assert.NoError(t, err, "failed to create test ctx") {
		return
	}
	defer ctx.Close()

	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	ctx.SetAPIServer(ts)

	user, err := testCreateUser(ctx, johndoe())
	if err != nil {
		return
	}
	defer testDeleteUser(ctx, user.ID, ctx.Superuser.EID)

	series, err := testCreateConferenceSeries(ctx, &model.CreateConferenceSeriesRequest{
		UserID: ctx.Superuser.EID,
		Slug:   tools.RandomString(8),
	})
	if err != nil {
		return
	}
	defer testDeleteConferenceSeries(ctx, series.ID, ctx.Superuser.EID)

	if err := testAddConferenceSeriesAdmin(ctx, series.ID, user.ID, ctx.Superuser.EID); err != nil {
		return
	}

	confs := make([]*model.ObjectID, 10)
	for i := 0; i < 10; i++ {
		lf := model.LocalizedFields{}
		lf.Set("ja", "title", `リストカンファレンステスト`)

		conf, err := testCreateConferencePass(ctx, &model.CreateConferenceRequest{
			L10N:     lf,
			SeriesID: series.ID,
			Slug:     tools.RandomString(8),
			Title:    "ListConference Test",
			UserID:   user.ID,
		})
		if err != nil {
			return
		}
		confs[i] = conf

		err = ctx.HTTPClient.AddConferenceDates(&model.AddConferenceDatesRequest{
			ConferenceID: conf.ID,
			UserID:       user.ID,
			Dates: []model.ConferenceDate{
				model.ConferenceDate{
					Date:  model.NewDate(2016, 3, 22),
					Open:  model.NewWallClock(10, 0),
					Close: model.NewWallClock(19, 0),
				},
			},
		})
		if !assert.NoError(ctx.T, err, "AddConferenceDates works") {
			return
		}

		defer testDeleteConference(ctx, conf.ID)

		req := &model.UpdateConferenceRequest{
			ID:     confs[i].ID,
			UserID: user.ID,
		}
		req.Status.Set("public")
		if err := testUpdateConference(ctx, req); err != nil {
			return
		}
	}

	in := model.ListConferenceRequest{}
	in.Lang.Set("ja")
	in.RangeStart.Set("2016-03-21")
	in.RangeEnd.Set("2016-03-23")
	res, err := ctx.HTTPClient.ListConference(&in)
	if !assert.NoError(ctx.T, err, "ListConference should succeed") {
		return
	}

	if !assert.NoError(ctx.T, validator.HTTPListConferenceResponse.Validate(res), "Validation should succeed") {
		return
	}

	if !assert.Len(ctx.T, res, 10, "ListConference returns 10 rooms") {
		return
	}

	for _, c := range res {
		if !assert.Equal(t, c.Title, "リストカンファレンステスト", "Title is in Japanese") {
			return
		}
	}

	// Make some of them private
	for i := 0; i < 5; i++ {
		req := &model.UpdateConferenceRequest{
			ID:     confs[i*2].ID,
			UserID: user.ID,
		}
		req.Status.Set("private")
		if err := testUpdateConference(ctx, req); err != nil {
			return
		}
	}

	res, err = ctx.HTTPClient.ListConference(&in)
	if !assert.NoError(ctx.T, err, "ListConference should succeed") {
		return
	}

	if !assert.Len(ctx.T, res, 5, "ListConference returns 5 conferences") {
		return
	}

	for _, c := range res {
		if !assert.Equal(t, c.Title, "リストカンファレンステスト", "Title is in Japanese") {
			return
		}
	}
}

func TestListRoom(t *testing.T) {
	ctx, err := NewTestCtx(t)
	if !assert.NoError(t, err, "failed to create test ctx") {
		return
	}
	defer ctx.Close()

	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	ctx.SetAPIServer(ts)

	venue, err := testCreateVenuePass(ctx, bigsight(ctx.Superuser.EID))
	if err != nil {
		return
	}

	room, err := testCreateRoomPass(ctx, intlConferenceRoom(venue.ID, ctx.Superuser.EID))
	if err != nil {
		return
	}
	defer testDeleteRoom(ctx, room.ID, ctx.Superuser.EID)

	in := model.ListRoomRequest{
		VenueID: venue.ID,
	}
	res, err := ctx.HTTPClient.ListRoom(&in)
	if !assert.NoError(ctx.T, err, "ListRoom should succeed") {
		return
	}

	if !assert.NoError(ctx.T, validator.HTTPListRoomResponse.Validate(res), "Validation should succeed") {
		return
	}

	if !assert.Len(ctx.T, res, 1, "ListRoom returns 1 rooms") {
		return
	}
}

func TestListSessions(t *testing.T) {
	ctx, err := NewTestCtx(t)
	if !assert.NoError(t, err, "failed to create test ctx") {
		return
	}
	defer ctx.Close()

	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	ctx.SetAPIServer(ts)

	user, err := testCreateUser(ctx, johndoe())
	if err != nil {
		return
	}
	defer testDeleteUser(ctx, user.ID, ctx.Superuser.EID)

	series, err := testCreateConferenceSeries(ctx, yapcasia(ctx.Superuser.EID))
	if err != nil {
		return
	}
	defer testDeleteConferenceSeries(ctx, series.ID, ctx.Superuser.EID)

	if err := testAddConferenceSeriesAdmin(ctx, series.ID, user.ID, ctx.Superuser.EID); err != nil {
		return
	}

	conference, err := testCreateConferencePass(ctx, yapcasiaTokyo(series.ID, user.ID))
	if err != nil {
		return
	}
	defer testDeleteConference(ctx, conference.ID)

	if err := testMakeConferencePublic(ctx, conference.ID, user.ID); err != nil {
		return
	}

	list, err := ctx.HTTPClient.ListSessionTypesByConference(&model.ListSessionTypesByConferenceRequest{
		ConferenceID: conference.ID,
	})
	if !assert.NoError(ctx.T, err) {
		return
	}
	if !assert.True(ctx.T, len(list) > 0) {
		return
	}

	stype := list[0]
	if err := testStartSubmission(ctx, stype.ID, user.ID, time.Now()); err != nil {
		return
	}

	for i := 0; i < 10; i++ {
		sin := model.CreateSessionRequest{}
		sin.ConferenceID = conference.ID
		sin.SpeakerID.Set(user.ID)
		sin.Title.Set(fmt.Sprintf("Title %d", i))
		sin.SessionTypeID = stype.ID
		sin.Abstract.Set("Use lots of reflection and generate lots of code")
		sin.UserID = user.ID
		s, err := testCreateSessionPass(ctx, &sin)
		if err != nil {
			return
		}
		defer testDeleteSessionPass(ctx, s.ID, user.ID)
	}

	in := model.ListSessionsRequest{}
	in.ConferenceID.Set(conference.ID)
	res, err := ctx.HTTPClient.ListSessions(&in)
	if !assert.NoError(ctx.T, err, "ListSessions should succeed") {
		return
	}
	t.Logf("%#v", res)
	if !assert.NoError(ctx.T, validator.HTTPListSessionsResponse.Validate(res), "Validation should succeed") {
		return
	}

	if !assert.Len(ctx.T, res, 10, "There should be 10 sessions") {
		return
	}

	in = model.ListSessionsRequest{}
	_, err = ctx.HTTPClient.ListSessions(&in)
	if !assert.Error(ctx.T, err, "Query without conference_id/speaker_id should be an error") {
		return
	}

	if !assert.Equal(ctx.T, err.Error(), "no query specified (one of conference_id/speaker_id is required)") {
		return
	}
}

func TestListVenue(t *testing.T) {
	ctx, err := NewTestCtx(t)
	if !assert.NoError(t, err, "failed to create test ctx") {
		return
	}
	defer ctx.Close()

	ts := httptest.NewServer(octav.New())
	defer ts.Close()

	ctx.SetAPIServer(ts)

	in := model.ListVenueRequest{}
	res, err := ctx.HTTPClient.ListVenue(&in)
	if !assert.NoError(ctx.T, err, "ListVenue should succeed") {
		return
	}
	if !assert.NoError(ctx.T, validator.HTTPListVenueResponse.Validate(res), "Validation should succeed") {
		return
	}
}
