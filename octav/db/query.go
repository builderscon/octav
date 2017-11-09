package db

import (
	"database/sql"
	"sync"

	pdebug "github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

var preparedStmts = make(map[string]*sql.Stmt)
var muPreparedStmts sync.Mutex

func Stmt(tx *sql.Tx, sqltext string) (stmt *sql.Stmt, err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.Stmt %s", sqltext).BindError(&err)
		defer g.End()
	}

	muPreparedStmts.Lock()
	defer muPreparedStmts.Unlock()

	stmt, ok := preparedStmts[sqltext]
	if ok {
		if pdebug.Enabled {
			pdebug.Printf("Cache HIT for statement")
		}
		return tx.Stmt(stmt), nil
	}

	if pdebug.Enabled {
		pdebug.Printf("Cache MISS for statement")
	}
	stmt, err = tx.Prepare(sqltext)
	if err != nil {
		return nil, errors.Wrap(err, `failed to prepare statement`)
	}

	preparedStmts[sqltext] = stmt
	return stmt, nil
}

func Exec(tx *sql.Tx, sqltext string, args ...interface{}) (rows sql.Result, err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.Exec %s %v", sqltext, args).BindError(&err)
		defer g.End()
	}

	stmt, err := Stmt(tx, sqltext)
	if err != nil {
		return nil, errors.Wrap(err, `db.Exec`)
	}
	return stmt.Exec(args...)
}

func Query(tx *sql.Tx, sqltext string, args ...interface{}) (rows *sql.Rows, err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.Query %s %v", sqltext, args).BindError(&err)
		defer g.End()
	}

	stmt, err := Stmt(tx, sqltext)
	if err != nil {
		return nil, errors.Wrap(err, `db.Query`)
	}
	return stmt.Query(args...)
}

func QueryRow(tx *sql.Tx, sqltext string, args ...interface{}) (row *sql.Row, err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("db.QueryRow %s %v", sqltext, args).BindError(&err)
		defer g.End()
	}

	stmt, err := Stmt(tx, sqltext)
	if err != nil {
		return nil, errors.Wrap(err, `db.QueryRow`)
	}
	return stmt.QueryRow(args...), nil
}
