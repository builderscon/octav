package model

import (
	"time"

	"github.com/builderscon/octav/octav/tools"
)

type ErrInvalidJSONFieldType struct {
	Field string
	Value interface{}
}

type ErrInvalidFieldType struct {
	Field string
}

type Conference struct {
	ID       string `json:"id"`
	Title    string `json:"title" l10n:"true"`
	SubTitle string `json:"sub_title" l10n:"true"`
	Slug     string `json:"slug"`
}
type ConferenceList []Conference

type Room struct {
	ID       string                `json:"id"`
	VenueID  string                `json:"venue_id"`
	Name     string                `json:"name" l10n:"true"`
	Capacity uint                  `json:"capacity"`
	L10N     tools.LocalizedFields `json:"-"`
}
type RoomList []RoomL10N

type Session struct {
	ID                string      `json:"id"`
	ConferenceID      string      `json:"conference_id"`
	RoomID            string      `json:"room_id"`
	SpeakerID         string      `json:"speaker_id"`
	Title             string      `json:"title" l10n:"true"`
	Abstract          string      `json:"abstract" l10n:"true"`
	Memo              string      `json:"memo"`
	StartsOn          time.Time   `json:"starts_on"`
	Duration          int         `json:"duration"`
	MaterialLevel     string      `json:"material_level"`
	Tags              TagString   `json:"tags,omitempty" assign:"convert"`
	Category          string      `json:"category,omitempty"`
	SpokenLanguage    string      `json:"spoken_language,omitempty"`
	SlideLanguage     string      `json:"slide_language,omitempty"`
	SlideSubtitles    string      `json:"slide_subtitles,omitempty"`
	SlideURL          string      `json:"slide_url,omitempty"`
	VideoURL          string      `json:"video_url,omitempty"`
	PhotoPermission   string      `json:"photo_permission"`
	VideoPermission   string      `json:"video_permission"`
	SortOrder         int         `json:"-"`
	HasInterpretation bool        `json:"has_interpretation"`
	Status            string      `json:"status"`
	Confirmed         bool        `json:"confirmed"`
	Conference        *Conference `json:"conference"` // only populated for JSON response
	Room              *Room       `json:"room"`       // only populated for JSON response
	Speaker           *User       `json:"speaker"`    // only populated for JSON response
}
type SessionL10NList []SessionL10N

type TagString string

type User struct {
	ID         string                `json:"id"`
	FirstName  string                `json:"first_name"`
	LastName   string                `json:"last_name"`
	Nickname   string                `json:"nickname"`
	Email      string                `json:"email"`
	TshirtSize string                `json:"tshirt_size"`
	L10N       tools.LocalizedFields `json:"-"`
}
type UserList []User

type Venue struct {
	ID        string                `json:"id,omitempty"`
	Name      string                `json:"name" l10n:"true"`
	Address   string                `json:"address" l10n:"true"`
	Longitude float64               `json:"longitude,omitempty"`
	Latitude  float64               `json:"latitude,omitempty"`
	L10N      tools.LocalizedFields `json:"-"`
}
type VenueL10NList []VenueL10N
