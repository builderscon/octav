package service

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/builderscon/octav/octav/cache"
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/internal/errors"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	pdebug "github.com/lestrrat/go-pdebug"
)

func (v *ClientSvc) Init() {
}

func (v *ClientSvc) populateRowForCreate(vdb *db.Client, payload *model.CreateClientRequest) error {
	vdb.EID = tools.RandomString(64)
	vdb.Secret = tools.RandomString(64)
	vdb.Name = payload.Name
	return nil
}

func (v *ClientSvc) populateRowForUpdate(vdb *db.Client, payload *model.UpdateClientRequest) error {
	vdb.Secret = payload.Secret
	vdb.Name = payload.Name
	return nil
}

func (v *ClientSvc) Authenticate(ctx context.Context, clientID, clientSecret string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Client.Authenticate").BindError(&err)
		defer g.End()
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}

	var vdb db.Client
	if err := vdb.LoadByEID(tx, clientID); err != nil {
		return errors.Wrap(err, "failed to load client ID")
	}

	if vdb.Secret != clientSecret {
		return errors.WithHTTPCode(errors.New("invalid secret"), http.StatusForbidden)
	}
	return nil
}

func clientSessionKey(sessionID, clientID string) string {
	buf := tools.GetBuffer()
	defer tools.ReleaseBuffer(buf)

	buf.WriteString(`client-session.`)
	buf.WriteString(clientID)
	buf.WriteString(`.`)
	buf.WriteString(sessionID)
	return buf.String()
}

func (v *ClientSvc) LoadClientSession(ctx context.Context, tx *sql.Tx, sessionID, clientID string, u *model.User) (err error) {
	key := clientSessionKey(sessionID, clientID)
	if pdebug.Enabled {
		g := pdebug.Marker("service.Client.LoadClientSession %s", key).BindError(&err)
		defer g.End()
	}

	// load the session
	cache := Cache()

	var userID string
	if err := cache.Get(key, &userID); err != nil {
		return errors.Wrap(err, `failed to fetch session`)
	}

	user := User()
	if err := user.Lookup(ctx, tx, u, userID); err != nil {
		return errors.Wrap(err, `failed to load user`)
	}
	return nil
}

func (v *ClientSvc) CreateClientSession(ctx context.Context, tx *sql.Tx, sessionID, clientID, userID string) (err error) {
	key := clientSessionKey(sessionID, clientID)
	if pdebug.Enabled {
		g := pdebug.Marker("service.Client.CreateClientSession %s", key).BindError(&err)
		defer g.End()
	}

	c := Cache()
	if err := c.Set(key, userID, cache.WithExpires(time.Hour)); err != nil {
		return errors.Wrap(err, `failed to fetch session`)
	}

	return nil
}
