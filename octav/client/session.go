package client

import (
	"net/http"

	"github.com/builderscon/octav/octav/model"
	pdebug "github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

type Session struct {
	*Client
}

const sessionIDHeaderKey = "X-Octav-Session-ID"

func NewSession(c *Client, token, userID string) (*Session, error) {
	in := model.CreateClientSessionRequest{
		AccessToken: token,
		UserID:      userID,
	}

	res, err := c.CreateClientSession(&in)
	if err != nil {
		return nil, errors.Wrap(err, `failed to create session`)
	}

	sid := res.SessionID
	c.SetMutator(func(r *http.Request) error {
		if pdebug.Enabled {
			pdebug.Printf("Setting `%s` to `%s`", sessionIDHeaderKey, sid)
		}
		r.Header.Set(sessionIDHeaderKey, sid)
		return nil
	})
	return &Session{Client: c}, nil
}
