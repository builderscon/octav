package service

import (
	"text/template"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

// InTesting grudingly exists to tell if we are running under
// testing mode.
var InTesting bool

type ErrInvalidJSONFieldType struct {
	Field string
}

type ErrInvalidFieldType struct {
	Field string
}

type ClientSvc struct{}
type ConferenceSvc struct {
	Storage StorageClient
}
type ConferenceComponentSvc struct{}
type ConferenceSeriesSvc struct{}
type FeaturedSpeakerSvc struct{}
type MailgunSvc struct {
	defaultSender string
	client        mailgun.Mailgun
}

type QuestionSvc struct{}
type RoomSvc struct{}
type SessionSvc struct{}
type SessionTypeSvc struct{}
type SponsorSvc struct {
	Storage StorageClient
}
type TemplateSvc struct {
	template *template.Template
}
type UserSvc struct{}
type VenueSvc struct{}
