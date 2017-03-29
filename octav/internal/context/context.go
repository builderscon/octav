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

type isTrustedCaller interface {
	IsTrustedCall() bool
}

func IsTrustedCall(ctx context.Context) bool {
	if tc, ok := ctx.(isTrustedCaller); ok {
		return tc.IsTrustedCall()
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
	pdebug.Printf("GetUser ctx = %#v", ctx)
	if c, ok := ctx.(withUser); ok {
		pdebug.Printf("withUser is true")
		return c.User()
	}
	return nil
}

func GetUserID(ctx context.Context) string {
	if u := GetUser(ctx); u != nil {
pdebug.Printf("GetUser returns %#v", u)
		return u.ID
	}
	return ""
}
