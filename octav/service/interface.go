package service

import (
	"context"
	"io"
	"sync"

	"cloud.google.com/go/storage"

	"github.com/dghubble/go-twitter/twitter"
	tmplbox "github.com/lestrrat/go-tmplbox"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

// InTesting grudingly exists to tell if we are running under
// testing mode.
var InTesting bool

type CallOption interface {
	Get() interface{}
}

type ObjectList interface {
	Next() bool
	Object() interface{}
	Error() error
}

type GoogleStorageObjectList struct {
	elements <-chan interface{}
	err      error
	mu       sync.Mutex
	next     interface{}
}

type WithObjectAttrs storage.ObjectAttrs
type WithQueryPrefix string

type StorageClient interface {
	URLFor(string) string
	List(ctx context.Context, options ...CallOption) (ObjectList, error)
	Move(ctx context.Context, src, dst string, options ...CallOption) error
	Upload(ctx context.Context, name string, src io.Reader, options ...CallOption) error
	Download(ctx context.Context, name string, dst io.Writer) error
	DeleteObjects(ctx context.Context, list ObjectList) error
}

type GoogleStorageClient struct {
	bucketName string
	clientOnce sync.Once
	Client     *storage.Client
}

var MediaStorage StorageClient
var CredentialStorage StorageClient

type ErrInvalidJSONFieldType struct {
	Field string
}

type ErrInvalidFieldType struct {
	Field string
}

type ClientSvc struct{}
type ConferenceSvc struct {
	mediaStorage      StorageClient
	credentialStorage StorageClient
}
type ConferenceComponentSvc struct{}
type ConferenceDateSvc struct{}
type ConferenceSeriesSvc struct{}

// +PostCreateHook
// +PostUpdateHook
// +PostDeleteHook
type ConferenceVenueSvc struct{}

// +PreUpdateFromPayloadHook
type FeaturedSpeakerSvc struct{}
type MailgunSvc struct {
	defaultSender string
	client        mailgun.Mailgun
}

type LocalizedStringSvc struct{}
type QuestionSvc struct{}

// +PostCreateHook
// +PostUpdateHook
// +PreUpdateFromPayloadHook
type RoomSvc struct{}

// +PreUpdateFromPayloadHook
type SessionSvc struct{}

// +PostUpdateHook
type SessionTypeSvc struct{}
type SponsorSvc struct {
	mediaStorage StorageClient
}
type TemplateSvc struct {
	box *tmplbox.Box
}
type TwitterSvc struct {
	*twitter.Client
}

// +PreUpdateFromPayloadHook
// +PostLookupFromPayloadHook
type UserSvc struct {
	EnableVerify bool
}

type VenueSvc struct{}

type TrackSvc struct {}
