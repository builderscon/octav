package octav_test

import (
	"os"

	"github.com/builderscon/octav/octav/db"
)

func init() {
	if dsn := os.Getenv("OCTAV_TEST_DSN"); dsn != "" {
		if err := db.Init(dsn); err != nil {
			panic(err.Error())
		}
	}
}
