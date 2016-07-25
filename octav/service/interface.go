package service

import "google.golang.org/cloud/storage"

type ErrInvalidJSONFieldType struct {
	Field string
}

type ErrInvalidFieldType struct {
	Field string
}

type Client struct{}
type Conference struct{}
type ConferenceSeries struct{}
type FeaturedSpeaker struct{}
type Question struct{}
type Room struct{}
type Session struct{}
type Sponsor struct {
	Storage *storage.Client
}
type User struct{}
type Venue struct{}
