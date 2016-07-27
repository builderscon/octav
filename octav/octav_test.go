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

func testCreateVenuePass(ctx *TestCtx, v *model.CreateVenueRequest) (*model.Venue, error) {
	return testCreateVenue(ctx, v, false)
}

func testCreateVenueFail(ctx *TestCtx, v *model.CreateVenueRequest) (*model.Venue, error) {
	return testCreateVenue(ctx, v, true)
}

func testCreateVenue(ctx *TestCtx, v *model.CreateVenueRequest, fail bool) (*model.Venue, error) {
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

func testCreateRoomPass(ctx *TestCtx, v *model.CreateRoomRequest) (*model.Room, error) {
	return testCreateRoom(ctx, v, false)
}

func testCreateRoomFail(ctx *TestCtx, v *model.CreateRoomRequest) (*model.Room, error) {
	return testCreateRoom(ctx, v, true)
}

func testCreateRoom(ctx *TestCtx, r *model.CreateRoomRequest, fail bool) (*model.Room, error) {
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

func testCreateConferencePass(ctx *TestCtx, in *model.CreateConferenceRequest) (*model.Conference, error) {
	return testCreateConference(ctx, in, false)
}

func testCreateConferenceFail(ctx *TestCtx, in *model.CreateConferenceRequest) (*model.Conference, error) {
	return testCreateConference(ctx, in, true)
}

func testCreateConference(ctx *TestCtx, in *model.CreateConferenceRequest, fail bool) (*model.Conference, error) {
	res, err := ctx.HTTPClient.CreateConference(in)
	if fail {
		if !assert.Error(ctx.T, err, "CreateConference should fail") {
			return nil, errors.New("expected operation to fail, but succeeded")
		}
		return nil, nil
	}
	if !assert.NoError(ctx.T, err, "CreateConference should succeed") {
		return nil, err
	}
	return res, nil
}
