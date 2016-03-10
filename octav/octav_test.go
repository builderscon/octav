package octav_test

import (
	"os"
	"testing"

	"github.com/builderscon/octav/octav"
	"github.com/builderscon/octav/octav/db"
	"github.com/stretchr/testify/assert"
)

func init() {
	if dsn := os.Getenv("OCTAV_TEST_DSN"); dsn != "" {
		if err := db.Init(dsn); err != nil {
			panic(err.Error())
		}
	}
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
