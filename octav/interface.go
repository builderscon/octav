package octav

import (
	"errors"
	"sync"
	"time"
)

var ErrInvalidFieldType = errors.New("placeholder error")

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

type Member struct{}
type Room struct {
	ID       string          `json:"id"`
	VenueID  string          `json:"venue_id"`
	Name     string          `json:"name"`
	Capacity uint            `json:"capacity"`
	L10N     LocalizedFields `json:"-"`
}
type RoomList []Room
type DeleteRoomRequest struct {
	ID string `json:"id" urlenc:"id"`
}
type LookupRoomRequest struct {
	ID string `json:"id" urlenc:"id"`
}

type SessionList []Session
type Session struct {
	ID                string     `json:"id"`
	ConferenceID      string     `json:"conference_id"`
	RoomID            string     `json:"room_id"`
	SpeakerID         string     `json:"speaker_id"`
	Title             string     `json:"title"`
	Abstract          string     `json:"abstract"`
	Memo              string     `json:"memo"`
	StartsOn          time.Time  `json:"starts_on"`
	Duration          int        `json:"duration"`
	MaterialLevel     string     `json:"material_level"`
	Tags              []string   `json:"tags,omitempty"`
	Category          string     `json:"category,omitempty"`
	SpokenLanguage    string     `json:"spoken_language,omitempty"`
	SlideLanguage     string     `json:"slide_language,omitempty"`
	SlideSubtitles    string     `json:"slide_subtitles,omitempty"`
	SlideURL          string     `json:"slide_url,omitempty"`
	VideoURL          string     `json:"video_url,omitempty"`
	PhotoPermission   string     `json:"photo_permission"`
	VideoPermission   string     `json:"video_permission"`
	HasInterpretation bool       `json:"has_interpretation"`
	Status            string     `json:"status"`
	SortOrder         int        `json:"sort_order"`
	Confirmed         bool       `json:"confirmed"`
	Conference        Conference `json:"conference"` // only populated for JSON response
	Room              Room       `json:"room"`       // only populated for JSON response
	Speaker           Member     `json:"speaker"`    // only populated for JSON response
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
type CreateUserRequest struct {
	FirstName  string          `json:"first_name"`
	LastName   string          `json:"last_name"`
	Nickname   string          `json:"nickname"`
	Email      string          `json:"email"`
	TshirtSize string          `json:"tshirt_size"`
	L10N       LocalizedFields `json:"-"`
}

type VenueList []Venue
type Venue struct {
	ID        string          `json:"id,omitempty"`
	Name      string          `json:"name"`
	Address   string          `json:"address"`
	Longitude float64         `json:"longitude,omitempty"`
	Latitude  float64         `json:"latitude,omitempty"`
	L10N      LocalizedFields `json:"-"`
}
type DeleteVenueRequest struct {
	ID string `json:"id" urlenc:"id"`
}
type LookupVenueRequest struct {
	ID string `json:"id" urlenc:"id"`
}

type ConferenceList []Conference
type Conference struct {
	ID       string           `json:"id"`
	Title    string           `json:"title"`
	SubTitle string           `json:"subtitle"`
	Slug     string           `json:"slug"`
	Dates    []ConferenceDate `json:"dates"` // only populated for JSON response
}

type LocalizedFields struct {
	lock sync.RWMutex
	// Language -> field/value
	fields map[string]map[string]string
}
