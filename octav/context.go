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

func (ctx *requestCtx) User() *model.User {
	return ctx.user
}
