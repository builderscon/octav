package service

// InTesting grudingly exists to tell if we are running under
// testing mode.
var InTesting bool

type ErrInvalidJSONFieldType struct {
	Field string
}

type ErrInvalidFieldType struct {
	Field string
}

type Client struct{}
type Conference struct{
	Storage StorageClient
}
type ConferenceSeries struct{}
type FeaturedSpeaker struct{}
type Question struct{}
type Room struct{}
type Session struct{}
type Sponsor struct {
	Storage StorageClient
}
type User struct{}
type Venue struct{}
