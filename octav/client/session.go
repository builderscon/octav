package client

import (
	"context"
	"net/http"
	"time"

	"github.com/builderscon/octav/octav/model"
	pdebug "github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

type Session struct {
	*Client
	sid     string
	expires time.Time
}

const sessionIDHeaderKey = "X-Octav-Session-ID"

func (s *Session) newSessionID(token, userID string) (string, time.Time, error) {
	in := model.CreateClientSessionRequest{
		AccessToken: token,
		UserID:      userID,
	}

	res, err := s.Client.CreateClientSession(&in)
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, `failed to create session`)
	}

	expires, err := time.Parse(time.RFC3339, res.Expires)
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, `failed to parse expires field`)
	}
	return res.SessionID, expires, nil
}

func (s *Session) updateSessionID(token, userID string) error {
	sid, expires, err := s.newSessionID(token, userID)
	if err != nil {
		return errors.Wrap(err, `failed to create new session ID`)
	}

	s.sid = sid
	s.expires = expires
	return nil
}

func (s *Session) modifyRequest(r *http.Request) error {
	if pdebug.Enabled {
		pdebug.Printf("Setting `%s` to `%s`", sessionIDHeaderKey, s.sid)
	}
	r.Header.Set(sessionIDHeaderKey, s.sid)
	return nil
}

func (s *Session) periodicUpdate(ctx context.Context, token, userID string) {
	for {
		next := time.Until(s.expires)
		next = next - (next % time.Second)
		t := time.NewTimer(next)
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
		t.Stop()

		s.updateSessionID(token, userID)
	}
}

func NewSession(ctx context.Context, c *Client, token, userID string) (*Session, error) {
	var s Session
	s.Client = c

	s.updateSessionID(token, userID)
	s.Client.SetMutator(s.modifyRequest)
	go s.periodicUpdate(ctx, token, userID)

	return &s, nil
}
