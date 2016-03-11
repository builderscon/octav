package octav_test

import (
	"flag"
	"os"
	"testing"

	"github.com/builderscon/octav/octav"
	"github.com/builderscon/octav/octav/db"
	"github.com/lestrrat/go-pdebug"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	flag.Parse()
	// We don't call m.Run() directly here so that we can
	// make sure that defer() gets fired.
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

	if dsn := os.Getenv("OCTAV_TEST_DSN"); dsn != "" {
		if pdebug.Enabled {
			pdebug.Printf("OCTAV_TEST_DSN (%s) is available. Initializing database", dsn)
		}
		if err := db.Init(dsn); err != nil {
			panic(err.Error())
		}
	}


	return m.Run()
}

func TestVenueDB(t *testing.T) {
	v := db.Venue{EID: octav.UUID()}
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
