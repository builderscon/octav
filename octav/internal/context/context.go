package context

import (
	"context"
	"time"

	"github.com/builderscon/octav/octav/model"
	pdebug "github.com/lestrrat/go-pdebug"
)

type Context interface {
	Deadline() (time.Time, bool)
	Done() <-chan struct{}
	Err() error
	Value(interface{}) interface{}
}

func Background() context.Context {
	return context.Background()
}

func WithCancel(ctx context.Context) (context.Context, func()) {
	return context.WithCancel(ctx)
}

func WithTimeout(ctx context.Context, t time.Duration) (context.Context, func()) {
	return context.WithTimeout(ctx, t)
}

type isInternalCaller interface {
	IsInternalCall() bool
}

func IsInternalCall(ctx context.Context) bool {
	if tc, ok := ctx.(isInternalCaller); ok {
		return tc.IsInternalCall()
	}
	return false
}

type isVerifiedCaller interface {
	IsVerifiedCall() bool
}

func IsVerifiedCall(ctx context.Context) bool {
	if tc, ok := ctx.(isVerifiedCaller); ok {
		return tc.IsVerifiedCall()
	}
	return false
}

type withClientID interface {
	ClientID() string
}

func GetClientID(ctx context.Context) string {
	if c, ok := ctx.(withClientID); ok {
		return c.ClientID()
	}
	return ""
}

type withUser interface {
	User() *model.User
}

func GetUser(ctx context.Context) *model.User {
	if c, ok := ctx.(withUser); ok {
		pdebug.Printf("withUser is true")
		return c.User()
	}
	return nil
}

func GetUserID(ctx context.Context) string {
	if u := GetUser(ctx); u != nil {
		return u.ID
	}
	return ""
}

type requestCtx struct {
	context.Context
	clientID     string
	internalCall bool
	sessionID    string
	user         *model.User
}

func NewRequestCtx(ctx context.Context, clientID string) context.Context {
	return &requestCtx{Context: ctx, clientID: clientID}
}

func (ctx *requestCtx) ClientID() string {
	return ctx.clientID
}

func (ctx *requestCtx) IsVerifiedCall() bool {
	return ctx.clientID != ""
}

func (ctx *requestCtx) IsInternalCall() bool {
	return ctx.internalCall
}

func (ctx *requestCtx) User() *model.User {
	return ctx.user
}

func WithUser(ctx context.Context, sessionID string, user *model.User) context.Context {
	switch ctx.(type) {
	case *requestCtx:
		rctx := ctx.(*requestCtx)
		return &requestCtx{
			Context:      ctx,
			clientID:     rctx.clientID,
			internalCall: rctx.internalCall,
			sessionID:    sessionID,
			user:         user,
		}
	default:
		return &requestCtx{
			Context:   ctx,
			sessionID: sessionID,
			user:      user,
		}
	}
}
func WithInternalCall(ctx context.Context, b bool) context.Context {
	switch ctx.(type) {
	case *requestCtx:
		rctx := ctx.(*requestCtx)
		return &requestCtx{
			Context:      ctx,
			clientID:     rctx.clientID,
			internalCall: b,
			sessionID:    rctx.sessionID,
			user:         rctx.user,
		}
	default:
		return &requestCtx{
			Context:      ctx,
			internalCall: b,
		}
	}
}
