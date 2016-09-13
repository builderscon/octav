package db

import (
	"crypto/sha512"
	"database/sql"

	pdebug "github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

var stmtPool = NewStmtPool()

func NewStmtPool() *StmtPool {
	return &StmtPool{
		pool: make(map[StmtKey]*StmtItem),
	}
}

func makeStmtKey(b []byte) StmtKey {
	return StmtKey(sha512.Sum512(b))
}


func (s *StmtPool) Register(key StmtKey, sqltext string) {
	if pdebug.Enabled {
		g := pdebug.Marker("StmtPool.Register %s", sqltext)
		defer g.End()
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.pool[key] = &StmtItem{Text: sqltext}
}

func (s *StmtPool) Get(key StmtKey) (ret *sql.Stmt, err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("StmtPool.Get").BindError(&err)
		defer g.End()
	}

	s.mutex.RLock()
	it, ok := s.pool[key]
	s.mutex.RUnlock()

	if !ok {
		return nil, errors.New("no such statement")
	}

	it.mutex.Lock()
	defer it.mutex.Unlock()
	if it.Stmt != nil {
		return it.Stmt, nil
	}

	stmt, err := _db.Prepare(it.Text)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare statement")
	}
	it.Stmt = stmt

	s.mutex.Lock()
	s.pool[key] = it
	s.mutex.Unlock()
	return it.Stmt, nil
}
