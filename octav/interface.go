package octav

import (
	"sync"
	"time"

	"github.com/lestrrat/go-jsval"
)

type ErrInvalidJSONFieldType struct {
	Field string
}

type ErrInvalidFieldType struct {
	Field string
}

type Date struct {
	Year  int
	Month int
	Day   int
}

type WallClock struct {
	Hour   int
	Minute int
}

// YYYY-MM-DD[HH:MM-HH:MM]
type ConferenceDate struct {
	Date  Date
	Open  WallClock
	Close WallClock
}

type Room struct {
	ID       string          `json:"id"`
	VenueID  string          `json:"venue_id"`
	Name     string          `json:"name" l10n:"true"`
	Capacity uint            `json:"capacity"`
	L10N     LocalizedFields `json:"-"`
}
type RoomList []Room
type DeleteRoomRequest struct {
	ID string `json:"id" urlenc:"id"`
}
type ListRoomRequest struct {
	VenueID string            `json:"venue_id" urlenc:"venue_id"`
	Since   jsval.MaybeString `json:"since" urlenc:"since"`
	Lang    jsval.MaybeString `json:"lang" urlenc:"lang"`
	Limit   jsval.MaybeInt    `json:"limit,omitempty" urlenc:"limit,omitempty"`
}
type LookupRoomRequest struct {
	ID string `json:"id" urlenc:"id"`
}

type Session struct {
	ID                string          `json:"id"`
	ConferenceID      string          `json:"conference_id"`
	RoomID            string          `json:"room_id"`
	SpeakerID         string          `json:"speaker_id"`
	Title             string          `json:"title"`
	Abstract          string          `json:"abstract"`
	Memo              string          `json:"memo"`
	StartsOn          time.Time       `json:"starts_on"`
	Duration          int             `json:"duration"`
	MaterialLevel     string          `json:"material_level"`
	Tags              []string        `json:"tags,omitempty"`
	Category          string          `json:"category,omitempty"`
	SpokenLanguage    string          `json:"spoken_language,omitempty"`
	SlideLanguage     string          `json:"slide_language,omitempty"`
	SlideSubtitles    string          `json:"slide_subtitles,omitempty"`
	SlideURL          string          `json:"slide_url,omitempty"`
	VideoURL          string          `json:"video_url,omitempty"`
	PhotoPermission   string          `json:"photo_permission"`
	VideoPermission   string          `json:"video_permission"`
	HasInterpretation bool            `json:"has_interpretation"`
	Status            string          `json:"status"`
	SortOrder         int             `json:"-"`
	Confirmed         bool            `json:"confirmed"`
	Conference        *Conference     `json:"conference"` // only populated for JSON response
	Room              *Room           `json:"room"`       // only populated for JSON response
	Speaker           *User           `json:"speaker"`    // only populated for JSON response
	L10N              LocalizedFields `json:"-"`
}
type SessionList []Session
type CreateSessionRequest struct {
	ConferenceID    string            `json:"conference_id,omitempty"`
	SpeakerID       string            `json:"speaker_id,omitempty"`
	Title           string            `json:"title,omitempty"`
	Abstract        jsval.MaybeString `json:"abstract,omitempty"`
	Memo            jsval.MaybeString `json:"memo,omitempty"`
	Duration        int               `json:"duration,omitempty"`
	MaterialLevel   jsval.MaybeString `json:"material_level,omitempty"`
	Tags            []string          `json:"tags,omitempty"`
	Category        jsval.MaybeString `json:"category,omitempty"`
	SpokenLanguage  jsval.MaybeString `json:"spoken_language,omitempty"`
	SlideLanguage   jsval.MaybeString `json:"slide_language,omitempty"`
	SlideSubtitles  jsval.MaybeString `json:"slide_subtitles,omitempty"`
	SlideURL        jsval.MaybeString `json:"slide_url,omitempty"`
	VideoURL        jsval.MaybeString `json:"video_url,omitempty"`
	PhotoPermission jsval.MaybeString `json:"photo_permission,omitempty"`
	VideoPermission jsval.MaybeString `json:"video_permission,omitempty"`
	L10N            LocalizedFields   `json:"-"`
}
type LookupSessionRequest struct {
	ID string `json:"id" urlenc:"id"`
}

type User struct {
	ID         string          `json:"id"`
	FirstName  string          `json:"first_name"`
	LastName   string          `json:"last_name"`
	Nickname   string          `json:"nickname"`
	Email      string          `json:"email"`
	TshirtSize string          `json:"tshirt_size"`
	L10N       LocalizedFields `json:"-"`
}
type UserList []User
type CreateUserRequest struct {
	FirstName  string          `json:"first_name"`
	LastName   string          `json:"last_name"`
	Nickname   string          `json:"nickname"`
	Email      string          `json:"email"`
	TshirtSize string          `json:"tshirt_size"`
	L10N       LocalizedFields `json:"-"`
}
type LookupUserRequest struct {
	ID string `json:"id" urlenc:"id"`
}
type DeleteUserRequest struct {
	ID string `json:"id" urlenc:"id"`
}

type Venue struct {
	ID        string          `json:"id,omitempty"`
	Name      string          `json:"name" l10n:"true"`
	Address   string          `json:"address" l10n:"true"`
	Longitude float64         `json:"longitude,omitempty"`
	Latitude  float64         `json:"latitude,omitempty"`
	L10N      LocalizedFields `json:"-"`
}
type VenueList []Venue
type DeleteVenueRequest struct {
	ID string `json:"id" urlenc:"id"`
}
type ListVenueRequest struct {
	Since string `json:"since" urlenc:"since"`
	Lang  string `json:"lang" urlenc:"lang"`
	Limit int    `json:"limit" urlenc:"limit"`
}
type LookupVenueRequest struct {
	ID string `json:"id" urlenc:"id"`
}

type Conference struct {
	ID       string           `json:"id"`
	Title    string           `json:"title"`
	SubTitle string           `json:"sub_title"`
	Slug     string           `json:"slug"`
	Dates    []ConferenceDate `json:"dates"` // only populated for JSON response
	L10N     LocalizedFields  `json:"-"`
}
type ConferenceList []Conference
type CreateConferenceRequest struct {
	Title    string          `json:"title"`
	SubTitle string          `json:"sub_title"`
	Slug     string          `json:"slug"`
	L10N     LocalizedFields `json:"-"`
}
type DeleteConferenceRequest struct {
	ID string `json:"id" urlenc:"id"`
}
type LookupConferenceRequest struct {
	ID string `json:"id" urlenc:"id"`
}

type LocalizedFields struct {
	lock sync.RWMutex
	// Language -> field/value
	fields map[string]map[string]string
}
