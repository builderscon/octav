package octav_test

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/builderscon/octav/octav/client"
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/service"
	"github.com/builderscon/octav/octav/tools"
	"github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	flag.Parse()
	// We don't call m.Run() directly here so that we can
	// make sure that defer() gets fired.
	service.InTesting = true
	os.Exit(setupAndRun(m))
}

func setupAndRun(m *testing.M) int {
	if pdebug.Enabled {
		if fn := os.Getenv("OCTAV_DEBUG_FILE"); fn != "" {
			f, err := os.Create(fn)
			if err == nil {
				// Before setting the output file, notify via
				// the standard channel what's going to happen
				pdebug.Printf("OCTAV_DEBUG_FILE (%s) is available. Redirecting pdebug output", fn)
				pdebug.DefaultCtx.Writer = f
				defer f.Close()
			} else {
				pdebug.DefaultCtx.Writer = os.Stderr
				pdebug.Printf("Failed to open file '%s': %s", fn, err)
			}
		}
	}

	var dsn string
	if dsn = os.Getenv("OCTAV_TEST_DSN"); dsn != "" {
		if pdebug.Enabled {
			pdebug.Printf("OCTAV_TEST_DSN (%s) is available", dsn)
		}
	}

	if pdebug.Enabled {
		pdebug.Printf("Initializing database...")
	}
	if err := db.Init(dsn); err != nil {
		panic(err.Error())
	}

	service.User().EnableVerify = false
	service.DefaultCacheMagic = tools.UUID()

	return m.Run()
}

func TestVenueDB(t *testing.T) {
	v := db.Venue{EID: tools.UUID()}
	if err := testVenueDBCreate(t, &v); err != nil {
		return
	}
	defer testVenueDBDelete(t, &v)
}

func testVenueDBDelete(t *testing.T, v *db.Venue) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if !assert.NoError(t, err, "Transaction starts") {
		return err
	}

	if err := v.Delete(tx); !assert.NoError(t, err, "Delete works") {
		return err
	}
	if err := tx.Commit(); !assert.NoError(t, err, "Commit works") {
		return err
	}
	return nil
}

func testVenueDBCreate(t *testing.T, v *db.Venue) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if !assert.NoError(t, err, "Transaction starts") {
		return err
	}

	if err := v.Create(tx); !assert.NoError(t, err, "Create works") {
		return err
	}
	if err := tx.Commit(); !assert.NoError(t, err, "Commit works") {
		return err
	}
	return nil
}

func testCreateVenuePass(ctx *TestCtx, v *model.CreateVenueRequest) (*model.ObjectID, error) {
	return testCreateVenue(ctx, v, false)
}

func testCreateVenueFail(ctx *TestCtx, v *model.CreateVenueRequest) (*model.ObjectID, error) {
	return testCreateVenue(ctx, v, true)
}

func testCreateVenue(ctx *TestCtx, v *model.CreateVenueRequest, fail bool) (*model.ObjectID, error) {
	var id *model.ObjectID
	err := withSession(ctx, v.UserID, func(s *client.Session) error {
		res, err := s.CreateVenue(v)
		if fail {
			if !assert.Error(ctx.T, err, "CreateVenue should fail") {
				return errors.New("expected operation to fail, but succeeded")
			}
			return nil
		}
		if !assert.NoError(ctx.T, err, "CreateVenue should succeed") {
			return err
		}
		id = res
		return nil
	})
	return id, err
}

func testCreateRoomPass(ctx *TestCtx, v *model.CreateRoomRequest) (*model.ObjectID, error) {
	return testCreateRoom(ctx, v, false)
}

func testCreateRoomFail(ctx *TestCtx, v *model.CreateRoomRequest) (*model.ObjectID, error) {
	return testCreateRoom(ctx, v, true)
}

func testCreateRoom(ctx *TestCtx, r *model.CreateRoomRequest, fail bool) (*model.ObjectID, error) {
	var id *model.ObjectID
	err := withSession(ctx, r.UserID, func(s *client.Session) error {
		res, err := s.Client.CreateRoom(r)
		if fail {
			if !assert.Error(ctx.T, err, "CreateRoom should fail") {
				return errors.New("expected operation to fail, but succeeded")
			}
			return nil
		}
		if !assert.NoError(ctx.T, err, "CreateRoom should succeed") {
			return err
		}
		id = res
		return nil
	})
	return id, err
}

func testCreateConferencePass(ctx *TestCtx, in *model.CreateConferenceRequest) (*model.ObjectID, error) {
	return testCreateConference(ctx, in, false)
}

func testCreateConferenceFail(ctx *TestCtx, in *model.CreateConferenceRequest) (*model.ObjectID, error) {
	return testCreateConference(ctx, in, true)
}

func testCreateConference(ctx *TestCtx, in *model.CreateConferenceRequest, fail bool) (*model.ObjectID, error) {
	var id *model.ObjectID
	err := withSession(ctx, in.UserID, func(s *client.Session) error {
		res, err := s.CreateConference(in)
		if fail {
			if !assert.Error(ctx.T, err, "CreateConference should fail") {
				return errors.Wrap(err, "expected operation to fail, but succeeded")
			}
			return nil
		}
		if !assert.NoError(ctx.T, err, "CreateConference should succeed") {
			return errors.Wrap(err, "expected operation to succeed, but failed")
		}
		id = res
		return nil
	})
	return id, err
}

func testCreateSessionPass(ctx *TestCtx, in *model.CreateSessionRequest) (*model.ObjectID, error) {
	return testCreateSession(ctx, in, false)
}

func testCreateSessionFail(ctx *TestCtx, in *model.CreateSessionRequest) (*model.ObjectID, error) {
	return testCreateSession(ctx, in, true)
}

func testCreateSession(ctx *TestCtx, in *model.CreateSessionRequest, fail bool) (*model.ObjectID, error) {
	var id *model.ObjectID
	err := withSession(ctx, in.UserID, func(s *client.Session) error {
		res, err := s.Client.CreateSession(in)
		if fail {
			if !assert.Error(ctx.T, err, "CreateSession should fail") {
				return errors.Wrap(err, "expected operation to fail, but succeeded")
			}
			return nil
		}
		if !assert.NoError(ctx.T, err, "CreateSession should succeed") {
			return errors.Wrap(err, "expected operation to suceed, but failed")
		}
		id = res
		return nil
	})
	return id, err
}

func testDeleteSponsor(ctx *TestCtx, id, userID string) error {
	s, err := ctx.getSession(userID)
	if err != nil {
		return errors.Wrap(err, `failed to find active session`)
	}

	err = s.DeleteSponsor(&model.DeleteSponsorRequest{
		ID:     id,
		UserID: userID,
	})
	if !assert.NoError(ctx.T, err, "DeleteSponsor should succeed") {
		return err
	}
	return nil
}

func testCreateSponsor(ctx *TestCtx, in *model.AddSponsorRequest) (*model.Sponsor, error) {
	var sponsor *model.Sponsor
	err := withSession(ctx, in.UserID, func(s *client.Session) error {
		res, err := s.AddSponsor(in)
		if !assert.NoError(ctx.T, err, "CreateSponsor should succeed") {
			return err
		}
		sponsor = res
		return nil
	})
	return sponsor, err
}

func testLookupSession(ctx *TestCtx, id, lang string) (*model.Session, error) {
	r := &model.LookupSessionRequest{ID: id}
	if lang != "" {
		r.Lang.Set(lang)
	}
	v, err := ctx.HTTPClient.LookupSession(r)
	if !assert.NoError(ctx.T, err, "LookupSession succeeds") {
		return nil, err
	}
	return v, nil
}

func testUpdateSession(ctx *TestCtx, in *model.UpdateSessionRequest) error {
	return withSession(ctx, in.UserID, func(s *client.Session) error {
		err := s.UpdateSession(in)
		if !assert.NoError(ctx.T, err, "UpdateSession succeeds") {
			return err
		}
		return nil
	})
}

func testDeleteSession(ctx *TestCtx, sessionID, userID string, fail bool) error {
	return withSession(ctx, userID, func(s *client.Session) error {
		err := s.DeleteSession(&model.DeleteSessionRequest{ID: sessionID, UserID: userID})
		if fail {
			if !assert.Error(ctx.T, err, "DeleteSession should fail") {
				return errors.New("expected operation to fail, but succeeded")
			}
			return nil
		}

		if !assert.NoError(ctx.T, err, "DeleteSession should be successful") {
			return err
		}
		return nil
	})
}

func testDeleteSessionPass(ctx *TestCtx, sessionID, userID string) error {
	return testDeleteSession(ctx, sessionID, userID, false)
}

func testDeleteSessionFail(ctx *TestCtx, sessionID, userID string) error {
	return testDeleteSession(ctx, sessionID, userID, true)
}

func testUpdateUser(ctx *TestCtx, r *model.UpdateUserRequest, fail bool) error {
	return withSession(ctx, r.UserID, func(s *client.Session) error {
		err := s.UpdateUser(r)
		if fail {
			if !assert.Error(ctx.T, err, "UpdateUser should fail") {
				return errors.New("expected operation to fail, but succeeded")
			}
			return nil
		}

		if !assert.NoError(ctx.T, err, "UpdateUser should be successful") {
			return err
		}
		return nil
	})
}

func testUpdateUserPass(ctx *TestCtx, r *model.UpdateUserRequest) error {
	return testUpdateUser(ctx, r, false)
}

func testUpdateUserFail(ctx *TestCtx, r *model.UpdateUserRequest) error {
	return testUpdateUser(ctx, r, true)
}

func testCreateConferenceSeries(ctx *TestCtx, in *model.CreateConferenceSeriesRequest) (*model.ObjectID, error) {
	s, err := ctx.getSession(in.UserID)
	if err != nil {
		return nil, errors.Wrap(err, `failed to find active session`)
	}

	res, err := s.CreateConferenceSeries(in)
	if !assert.NoError(ctx.T, err, "CreateConferenceSeries should succeed") {
		return nil, err
	}
	return res, nil
}

func testDeleteConferenceSeries(ctx *TestCtx, id, userID string) error {
	return withSession(ctx, userID, func(s *client.Session) error {
		err := s.DeleteConferenceSeries(&model.DeleteConferenceSeriesRequest{ID: id, UserID: userID})
		if !assert.NoError(ctx.T, err, "DeleteConferenceSeries should be successful") {
			return err
		}
		return nil
	})
}

func withSession(ctx *TestCtx, userID string, cb func(*client.Session) error) error {
	s, err := ctx.getSession(userID)
	if !assert.NoError(ctx.T, err, `failed to get session`) {
		return errors.Wrap(err, `failed to get session`)
	}

	return cb(s)
}

func testAddConferenceSeriesAdmin(ctx *TestCtx, id, adminID, userID string) error {
	return withSession(ctx, userID, func(s *client.Session) error {
		err := s.AddConferenceSeriesAdmin(&model.AddConferenceSeriesAdminRequest{SeriesID: id, AdminID: adminID, UserID: userID})
		if !assert.NoError(ctx.T, err, "AddConferenceSeriesAdmin should be successful") {
			return err
		}
		return nil
	})
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
	s, err := ctx.getSession(in.UserID)
	if !assert.NoError(ctx.T, err, `failed to get session`) {
		return errors.Wrap(err, `failed to get session`)
	}
	err = s.UpdateConference(in, nil)
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

func testDeleteConference(ctx *TestCtx, id, userID string) error {
	return withSession(ctx, userID, func(s *client.Session) error {
		err := s.DeleteConference(&model.DeleteConferenceRequest{ID: id, UserID: userID})
		if !assert.NoError(ctx.T, err, "DeleteConference should be successful") {
			return err
		}
		return nil
	})
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
	return withSession(ctx, userID, func(s *client.Session) error {
		err := s.DeleteUser(&model.DeleteUserRequest{ID: targetID, UserID: userID})
		if !assert.NoError(ctx.T, err, "DeleteUser should succeed") {
			return err
		}
		return nil
	})
}

func testAddConferenceVenue(ctx *TestCtx, confID, venueID, userID string) error {
	return withSession(ctx, userID, func(s *client.Session) error {
		req := model.AddConferenceVenueRequest{
			ConferenceID: confID,
			VenueID:      venueID,
			UserID:       userID,
		}
		err := s.AddConferenceVenue(&req)
		if !assert.NoError(ctx.T, err, "AddConferenceVenue should succeed") {
			return err
		}
		return nil
	})
}

func testDeleteConferenceVenue(ctx *TestCtx, confID, venueID, userID string) error {
	return withSession(ctx, userID, func(s *client.Session) error {
		req := model.DeleteConferenceVenueRequest{
			ConferenceID: confID,
			VenueID:      venueID,
			UserID:       userID,
		}
		err := s.DeleteConferenceVenue(&req)
		if !assert.NoError(ctx.T, err, "DeleteConferenceVenue should succeed") {
			return err
		}
		return nil
	})
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

func testCreateFeaturedSpeaker(ctx *TestCtx, in *model.AddFeaturedSpeakerRequest) (*model.FeaturedSpeaker, error) {
	var speaker *model.FeaturedSpeaker
	return speaker, withSession(ctx, in.UserID, func(s *client.Session) error {
		res, err := s.AddFeaturedSpeaker(in)
		if !assert.NoError(ctx.T, err, "CreateFeaturedSpeaker should succeed") {
			return err
		}
		speaker = res
		return nil
	})
}

func testDeleteFeaturedSpeaker(ctx *TestCtx, id, userID string) error {
	s, err := ctx.getSession(userID)
	if !assert.NoError(ctx.T, err, `failed to get session`) {
		return errors.New(`failed to get session`)
	}
	err = s.DeleteFeaturedSpeaker(&model.DeleteFeaturedSpeakerRequest{
		ID:     id,
		UserID: userID,
	})
	if !assert.NoError(ctx.T, err, "DeleteFeaturedSpeaker should succeed") {
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

func testUpdateRoom(ctx *TestCtx, in *model.UpdateRoomRequest) error {
	return withSession(ctx, in.UserID, func(s *client.Session) error {
		err := s.UpdateRoom(in)
		if !assert.NoError(ctx.T, err, "UpdateRoom succeeds") {
			return err
		}
		return nil
	})
}

func testDeleteRoom(ctx *TestCtx, roomID, userID string) error {
	return withSession(ctx, userID, func(s *client.Session) error {
		err := s.DeleteRoom(&model.DeleteRoomRequest{ID: roomID, UserID: userID})
		if !assert.NoError(ctx.T, err, "DeleteRoom should be successful") {
			return err
		}
		return err
	})
}

func testAddConferenceAdmin(ctx *TestCtx, confID, adminID, userID string) error {
	return withSession(ctx, userID, func(s *client.Session) error {
		err := s.AddConferenceAdmin(&model.AddConferenceAdminRequest{
			ConferenceID: confID,
			AdminID:      adminID,
			UserID:       userID,
		})
		if !assert.NoError(ctx.T, err, "AddConferenceAdmin should succeed") {
			return err
		}
		return nil
	})
}

func testDeleteConferenceAdmin(ctx *TestCtx, confID, adminID, userID string) error {
	return withSession(ctx, userID, func(s *client.Session) error {
		err := s.DeleteConferenceAdmin(&model.DeleteConferenceAdminRequest{
			ConferenceID: confID,
			AdminID:      adminID,
			UserID:       userID,
		})
		if !assert.NoError(ctx.T, err, "DeleteConferenceAdmin should succeed") {
			return err
		}
		return nil
	})
}

func testUpdateVenue(ctx *TestCtx, in *model.UpdateVenueRequest) error {
	return withSession(ctx, in.UserID, func(s *client.Session) error {
		err := s.UpdateVenue(in)
		if !assert.NoError(ctx.T, err, "UpdateVenue succeeds") {
			return err
		}
		return nil
	})
}

func testDeleteVenue(ctx *TestCtx, venueID, userID string) error {
	return withSession(ctx, userID, func(s *client.Session) error {
		err := s.DeleteVenue(&model.DeleteVenueRequest{ID: venueID, UserID: userID})
		if !assert.NoError(ctx.T, err, "DeleteVenue should be successful") {
			return err
		}
		return nil
	})
}

func testStartSubmission(ctx *TestCtx, sessionTypeID, userID string, ref time.Time) error {
	return withSession(ctx, userID, func(s *client.Session) error {
		r := &model.UpdateSessionTypeRequest{
			ID:     sessionTypeID,
			UserID: userID,
		}
		r.SubmissionStart.Set(ref.Add(-1 * 24 * time.Hour).Format(time.RFC3339))
		r.SubmissionEnd.Set(ref.Add(24 * time.Hour).Format(time.RFC3339))
		err := s.UpdateSessionType(r)
		if !assert.NoError(ctx.T, err, "StartSessionSubmission should be successful") {
			return err
		}
		return nil
	})
}
