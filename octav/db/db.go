package db

import (
	"bytes"
	"database/sql"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"text/template"

	"github.com/lestrrat/go-pdebug"
	"github.com/lestrrat/go-tx-guard"
)

type DB struct {
	*guard.DB
}

type Tx struct {
	*guard.Tx
}

type dsnvars struct {
	Address   string
	DBName    string
	EnableTLS bool
	Password  string
	Port      int
	Username  string
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

func Init(dsn string) error {
	if dsn == "" {
		dsn = defaultDSN()
	}

	// dsn can be a template, so we can incorporate automatically
	// acquired data
	t, err := template.New("dsn").Parse(dsn)
	if err != nil {
		return err
	}

	vars := defaultDSNVars()
	switch err := trySetupTLS(); err {
	case ErrNoTLSRequested:
		// no op. we weren't requested to do TLS
	case nil:
		// successfully connectd using TLS
		vars.EnableTLS = true
	default:
		// now *this* is an error
		return err
	}

	// This is usually not a good idea, except when using along with
	// stuff like Kubernetes secrets
	if f := os.Getenv("OCTAV_MYSQL_PASSWORD_FILE"); f != "" {
		if v, err := ioutil.ReadFile(f); err == nil {
			vars.Password = string(v)
		}
	}
	if v := os.Getenv("OCTAV_MYSQL_USERNAME"); v != "" {
		vars.Username = v
	}

	if v := os.Getenv("OCTAV_MYSQL_PORT"); v != "" {
		if p, err := strconv.ParseInt(v, 10, 64); err == nil {
			vars.Port = int(p)
		}
	}

	if v := os.Getenv("OCTAV_MYSQL_ADDRESS"); v != "" {
		vars.Address = v
	}

	if v := os.Getenv("OCTAV_MYSQL_DBNAME"); v != "" {
		vars.DBName = v
	}

	buf := bytes.Buffer{}
	if err := t.Execute(&buf, vars); err != nil {
		dsn = buf.String()
	}

	dn := driverName()
	if pdebug.Enabled {
		pdebug.Printf("Connecting to %s %s", dn, dsn)
	}
	conn, err := sql.Open(dn, dsn)
	if err != nil {
		return err
	}

	_db = &DB{&guard.DB{conn}}

	return nil
}

// Begin creates a new transaction (`Tx`) from the current
// global database connection
func Begin() (*Tx, error) {
	tx, err := _db.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}
