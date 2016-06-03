package service

type ErrInvalidJSONFieldType struct {
	Field string
}

type ErrInvalidFieldType struct {
	Field string
}

type Conference struct{}
type Question struct{}
type Room struct{}
type Session struct{}
type User struct{}
type Venue struct{}
