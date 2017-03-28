package octav

import (
	"context"

	"github.com/builderscon/octav/octav/model"
)

type requestCtx struct {
	context.Context
	clientID    string
	sessionID   string
	trustedCall bool
	user        *model.User
}

func (ctx *requestCtx) ClientID() string {
	return ctx.clientID
}

func (ctx *requestCtx) IsTrustedCall() bool {
	return ctx.trustedCall
}

type isTrustedCaller interface {
	IsTrustedCall() bool
}

func isTrustedCall(ctx context.Context) bool {
	if tc, ok := ctx.(isTrustedCaller); ok {
		return tc.IsTrustedCall()
	}
	return false
}

type withClientID interface {
	ClientID() string
}

func getClientID(ctx context.Context) string {
	if c, ok := ctx.(withClientID); ok {
		return c.ClientID()
	}
	return ""
}
