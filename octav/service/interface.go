package service

import (
	"sync"

	"google.golang.org/cloud/storage"
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

type Client struct{}
type Conference struct{}
type ConferenceSeries struct{}
type FeaturedSpeaker struct{}
type Question struct{}
type Room struct{}
type Session struct{}
type Sponsor struct {
	bucketOnce      sync.Once
	storageOnce     sync.Once
	MediaBucketName string
	Storage         *storage.Client
}
type User struct{}
type Venue struct{}
