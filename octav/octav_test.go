package octav_test

import (
	"flag"
	"os"
	"testing"

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
	tx, err := db.Begin()
	if !assert.NoError(t, err, "Transaction starts") {
		return err
	}
	defer tx.AutoRollback()

	if err := v.Delete(tx); !assert.NoError(t, err, "Delete works") {
		return err
	}
	if err := tx.Commit(); !assert.NoError(t, err, "Commit works") {
		return err
	}
	return nil
}

func testVenueDBCreate(t *testing.T, v *db.Venue) error {
	tx, err := db.Begin()
	if !assert.NoError(t, err, "Transaction starts") {
		return err
	}
	defer tx.AutoRollback()

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
	res, err := ctx.HTTPClient.CreateVenue(v)
	if fail {
		if !assert.Error(ctx.T, err, "CreateVenue should fail") {
			return nil, errors.New("expected operation to fail, but succeeded")
		}
		return nil, nil
	}
	if !assert.NoError(ctx.T, err, "CreateVenue should succeed") {
		return nil, err
	}
	return res, nil
}

func testCreateRoomPass(ctx *TestCtx, v *model.CreateRoomRequest) (*model.ObjectID, error) {
	return testCreateRoom(ctx, v, false)
}

func testCreateRoomFail(ctx *TestCtx, v *model.CreateRoomRequest) (*model.ObjectID, error) {
	return testCreateRoom(ctx, v, true)
}

func testCreateRoom(ctx *TestCtx, r *model.CreateRoomRequest, fail bool) (*model.ObjectID, error) {
	res, err := ctx.HTTPClient.CreateRoom(r)
	if fail {
		if !assert.Error(ctx.T, err, "CreateRoom should fail") {
			return nil, errors.New("expected operation to fail, but succeeded")
		}
		return nil, nil
	}
	if !assert.NoError(ctx.T, err, "CreateRoom should succeed") {
		return nil, err
	}
	return res, nil
}

func testCreateConferencePass(ctx *TestCtx, in *model.CreateConferenceRequest) (*model.ObjectID, error) {
	return testCreateConference(ctx, in, false)
}

func testCreateConferenceFail(ctx *TestCtx, in *model.CreateConferenceRequest) (*model.ObjectID, error) {
	return testCreateConference(ctx, in, true)
}

func testCreateConference(ctx *TestCtx, in *model.CreateConferenceRequest, fail bool) (*model.ObjectID, error) {
	res, err := ctx.HTTPClient.CreateConference(in)
	if fail {
		if !assert.Error(ctx.T, err, "CreateConference should fail") {
			return nil, errors.Wrap(err, "expected operation to fail, but succeeded")
		}
		return nil, nil
	}
	if !assert.NoError(ctx.T, err, "CreateConference should succeed") {
		return nil, errors.Wrap(err, "expected operation to succeed, but failed")
	}
	return res, nil
}

func testCreateSessionPass(ctx *TestCtx, in *model.CreateSessionRequest) (*model.ObjectID, error) {
	return testCreateSession(ctx, in, false)
}

func testCreateSessionFail(ctx *TestCtx, in *model.CreateSessionRequest) (*model.ObjectID, error) {
	return testCreateSession(ctx, in, true)
}

func testCreateSession(ctx *TestCtx, in *model.CreateSessionRequest, fail bool) (*model.ObjectID, error) {
	res, err := ctx.HTTPClient.CreateSession(in)
	if fail {
		if !assert.Error(ctx.T, err, "CreateSession should fail") {
			return nil, errors.Wrap(err, "expected operation to fail, but succeeded")
		}
		return nil, nil
	}
	if !assert.NoError(ctx.T, err, "CreateSession should succeed") {
		return nil, errors.Wrap(err, "expected operation to suceed, but failed")
	}
	return res, nil
}

func testDeleteSponsor(ctx *TestCtx, id, userID string) error {
	err := ctx.HTTPClient.DeleteSponsor(&model.DeleteSponsorRequest{
		ID:     id,
		UserID: userID,
	})
	if !assert.NoError(ctx.T, err, "DeleteSponsor should succeed") {
		return err
	}
	return nil
}

func testCreateSponsor(ctx *TestCtx, in *model.AddSponsorRequest) (*model.Sponsor, error) {
	res, err := ctx.HTTPClient.AddSponsor(in)
	if !assert.NoError(ctx.T, err, "CreateSponsor should succeed") {
		return nil, err
	}
	return res, nil
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
	err := ctx.HTTPClient.UpdateSession(in)
	if !assert.NoError(ctx.T, err, "UpdateSession succeeds") {
		return err
	}
	return nil
}

func testDeleteSession(ctx *TestCtx, sessionID, userID string, fail bool) error {
	err := ctx.HTTPClient.DeleteSession(&model.DeleteSessionRequest{ID: sessionID, UserID: userID})
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
}

func testDeleteSessionPass(ctx *TestCtx, sessionID, userID string) error {
	return testDeleteSession(ctx, sessionID, userID, false)
}

func testDeleteSessionFail(ctx *TestCtx, sessionID, userID string) error {
	return testDeleteSession(ctx, sessionID, userID, true)
}

func testUpdateUser(ctx *TestCtx, r *model.UpdateUserRequest, fail bool) error {
	err := ctx.HTTPClient.UpdateUser(r)
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
}

func testUpdateUserPass(ctx *TestCtx, r *model.UpdateUserRequest) error {
	return testUpdateUser(ctx, r, false)
}

func testUpdateUserFail(ctx *TestCtx, r *model.UpdateUserRequest) error {
	return testUpdateUser(ctx, r, true)
}


