package db

import (
	"database/sql"
	"os"
	"strconv"

	"github.com/lestrrat/go-pdebug"
	"github.com/lestrrat/go-tx-guard"
	"github.com/pkg/errors"
)

type DB struct {
	*guard.DB
}

type Tx struct {
	*guard.Tx
}

var _db *DB // global database connection
var ErrNoTLSRequested = errors.New("TLS environment variables not set")
var Trace bool

func init() {
	if f := os.Getenv("OCTAV_TRACE_DB"); f != "" {
		if b, err := strconv.ParseBool(f); b && err == nil {
			Trace = true
		}
	}
}

func Init(dsn string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.Init").BindError(&err)
		defer g.End()
	}

	if dsn == "" {
		dsn, err = ConfigureDSN()
		if err != nil {
			return err
		}
	}

	dn := driverName()
	conn, err := sql.Open(dn, dsn)
	if err != nil {
		return err
	}

	if err := onConnect(conn); err != nil {
		return err
	}

	_db = &DB{&guard.DB{conn}}

	return nil
}

// Begin creates a new transaction (`Tx`) from the current
// global database connection
func Begin() (*Tx, error) {
	if _db == nil {
		return nil, errors.New("database has not been initialized")
	}

	tx, err := _db.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}
