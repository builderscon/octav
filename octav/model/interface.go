package model

import (
	"errors"
	"time"

	"github.com/builderscon/octav/octav/tools"
	"github.com/lestrrat/go-jsval"
)

var ErrInvalidConferenceHour = errors.New("invalid conference hour specification")

type ErrInvalidJSONFieldType struct {
	Field string
	Value interface{}
}

type ErrInvalidFieldType struct {
	Field string
}

// +model
type Conference struct {
	ID             string             `json:"id"`
	Title          string             `json:"title" l10n:"true"`
	SubTitle       string             `json:"sub_title" l10n:"true"`
	Slug           string             `json:"slug"`
	Dates          ConferenceDateList `json:"dates,omitempty"`
	Administrators UserList           `json:"administrators,omitempty"`
}
type ConferenceL10NList []ConferenceL10N
type ConferenceList []Conference

// +model
type Room struct {
	ID       string `json:"id"`
	VenueID  string `json:"venue_id"`
	Name     string `json:"name" l10n:"true"`
	Capacity uint   `json:"capacity"`
}
type RoomList []RoomL10N

// +model
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

// +model
type User struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name" l10n:"true"`
	LastName   string `json:"last_name" l10n:"true"`
	Nickname   string `json:"nickname"`
	Email      string `json:"email"`
	TshirtSize string `json:"tshirt_size"`
}
type UserList []User

// +model
type Venue struct {
	ID        string  `json:"id,omitempty"`
	Name      string  `json:"name" l10n:"true"`
	Address   string  `json:"address" l10n:"true"`
	Longitude float64 `json:"longitude,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
}
type VenueL10NList []VenueL10N

// +transport
type CreateConferenceRequest struct {
	Title    string                `json:"title" l10n:"true"`
	SubTitle jsval.MaybeString     `json:"sub_title" l10n:"true"`
	Slug     string                `json:"slug"`
	UserID   string                `json:"user_id"`
	L10N     tools.LocalizedFields `json:"-"`
}

// +transport
type LookupConferenceRequest struct {
	ID   string            `json:"id" urlenc:"id"`
	Lang jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`
}

// +transport
type UpdateConferenceRequest struct {
	ID       string            `json:"id"`
	Title    jsval.MaybeString `json:"title,omitempty" l10n:"true"`
	SubTitle jsval.MaybeString `json:"sub_title,omitempty" l10n:"true"`
	Slug     jsval.MaybeString `json:"slug,omitempty"`
	// TODO dates
	L10N tools.LocalizedFields `json:"-"`
}

// Date is used to store simple dates YYYY-MM-DD
type Date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}
type DateList []Date

// WallClock is used to store simple time HH:MM
type WallClock struct {
	hour   int
	minute int
	Valid  bool // True if set
}

// YYYY-MM-DD[HH:MM-HH:MM]
type ConferenceDate struct {
	Date  Date
	Open  WallClock
	Close WallClock
}
type ConferenceDateList []ConferenceDate

// +transport
type AddConferenceDatesRequest struct {
	ConferenceID string             `json:"conference_id"`
	Dates        ConferenceDateList `json:"dates" extract:"true"`
}

// +transport
type AddConferenceAdminRequest struct {
	ConferenceID string `json:"conference_id"`
	UserID       string `json:"user_id"`
}

// +transport
type DeleteConferenceDatesRequest struct {
	ConferenceID string   `json:"conference_id"`
	Dates        DateList `json:"dates" extract:"true"`
}

// +transport
type DeleteConferenceAdminRequest struct {
	ConferenceID string   `json:"conference_id"`
	UserID       string   `json:"user_id"`
}

// +transport
type DeleteConferenceRequest struct {
	ID string `json:"id" urlenc:"id"`
}

// +transport
type ListConferencesRequest struct {
	RangeEnd   jsval.MaybeString `json:"range_end,omitempty" urlenc:"range_end,omitempty,string"`
	RangeStart jsval.MaybeString `json:"range_start,omitempty" urlenc:"range_start,omitempty,string"`
	Since      jsval.MaybeString `json:"since,omitempty" urlenc:"since,omitempty,string"`
	Lang       jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`
	Limit      jsval.MaybeInt    `json:"limit,omitempty" urlenc:"limit,omitempty,int64"`
}

// +transport
type CreateRoomRequest struct {
	VenueID  jsval.MaybeString     `json:"venue_id"`
	Name     jsval.MaybeString     `json:"name" l10n:"true"`
	Capacity jsval.MaybeUint       `json:"capacity"`
	L10N     tools.LocalizedFields `json:"-"`
}

// +transport
type LookupRoomRequest struct {
	ID   string            `json:"id" urlenc:"id"`
	Lang jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`
}

// +transport
type UpdateRoomRequest struct {
	ID       string                `json:"id"`
	VenueID  jsval.MaybeString     `json:"venue_id,omitempty"`
	Name     jsval.MaybeString     `json:"name,omitempty" l10n:"true"`
	Capacity jsval.MaybeUint       `json:"capacity,omitempty"`
	L10N     tools.LocalizedFields `json:"-"`
}

// +transport
type DeleteRoomRequest struct {
	ID string `json:"id" urlenc:"id"`
}

// +transport
type ListRoomRequest struct {
	VenueID string            `json:"venue_id" urlenc:"venue_id"`
	Since   jsval.MaybeString `json:"since" urlenc:"since,omitempty,string"`
	Lang    jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`
	Limit   jsval.MaybeInt    `json:"limit,omitempty" urlenc:"limit,omitempty,int64"`
}

// +transport
type CreateSessionRequest struct {
	ConferenceID    jsval.MaybeString     `json:"conference_id,omitempty"`
	SpeakerID       jsval.MaybeString     `json:"speaker_id,omitempty"`
	Title           jsval.MaybeString     `json:"title,omitempty"`
	Abstract        jsval.MaybeString     `json:"abstract,omitempty"`
	Memo            jsval.MaybeString     `json:"memo,omitempty"`
	Duration        jsval.MaybeInt        `json:"duration,omitempty"`
	MaterialLevel   jsval.MaybeString     `json:"material_level,omitempty"`
	Tags            jsval.MaybeString     `json:"tags,omitempty"`
	Category        jsval.MaybeString     `json:"category,omitempty"`
	SpokenLanguage  jsval.MaybeString     `json:"spoken_language,omitempty"`
	SlideLanguage   jsval.MaybeString     `json:"slide_language,omitempty"`
	SlideSubtitles  jsval.MaybeString     `json:"slide_subtitles,omitempty"`
	SlideURL        jsval.MaybeString     `json:"slide_url,omitempty"`
	VideoURL        jsval.MaybeString     `json:"video_url,omitempty"`
	PhotoPermission jsval.MaybeString     `json:"photo_permission,omitempty"`
	VideoPermission jsval.MaybeString     `json:"video_permission,omitempty"`
	L10N            tools.LocalizedFields `json:"-"`
}

// +transport
type LookupSessionRequest struct {
	ID   string            `json:"id" urlenc:"id"`
	Lang jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`
}

// +transport
type UpdateSessionRequest struct {
	ID                string                `json:"id"`
	ConferenceID      jsval.MaybeString     `json:"conference_id,omitempty"`
	SpeakerID         jsval.MaybeString     `json:"speaker_id,omitempty"`
	Title             jsval.MaybeString     `json:"title,omitempty"`
	Abstract          jsval.MaybeString     `json:"abstract,omitempty"`
	Memo              jsval.MaybeString     `json:"memo,omitempty"`
	Duration          jsval.MaybeInt        `json:"duration,omitempty"`
	MaterialLevel     jsval.MaybeString     `json:"material_level,omitempty"`
	Tags              jsval.MaybeString     `json:"tags,omitempty"`
	Category          jsval.MaybeString     `json:"category,omitempty"`
	SpokenLanguage    jsval.MaybeString     `json:"spoken_language,omitempty"`
	SlideLanguage     jsval.MaybeString     `json:"slide_language,omitempty"`
	SlideSubtitles    jsval.MaybeString     `json:"slide_subtitles,omitempty"`
	SlideURL          jsval.MaybeString     `json:"slide_url,omitempty"`
	VideoURL          jsval.MaybeString     `json:"video_url,omitempty"`
	PhotoPermission   jsval.MaybeString     `json:"photo_permission,omitempty"`
	VideoPermission   jsval.MaybeString     `json:"video_permission,omitempty"`
	SortOrder         jsval.MaybeInt        `json:"sort_order,omitempty"`
	HasInterpretation jsval.MaybeBool       `json:"has_interpretation,omitempty"`
	Status            jsval.MaybeString     `json:"status,omitempty"`
	Confirmed         jsval.MaybeBool       `json:"confirmed,omitempty"`
	L10N              tools.LocalizedFields `json:"-"`
}

// +transport
type DeleteSessionRequest struct {
	ID string `json:"id" urlenc:"id"`
}

// +transport
type CreateUserRequest struct {
	FirstName  string                `json:"first_name" l18n:"true"`
	LastName   string                `json:"last_name" l18n:"true"`
	Nickname   string                `json:"nickname"`
	Email      string                `json:"email"`
	TshirtSize string                `json:"tshirt_size"`
	L10N       tools.LocalizedFields `json:"-"`
}

// +transport
type UpdateUserRequest struct {
	ID         string                `json:"id"`
	FirstName  jsval.MaybeString     `json:"first_name,omitempty"`
	LastName   jsval.MaybeString     `json:"last_name,omitempty"`
	Nickname   jsval.MaybeString     `json:"nickname,omitempty"`
	Email      jsval.MaybeString     `json:"email,omitempty"`
	TshirtSize jsval.MaybeString     `json:"tshirt_size,omitempty"`
	L10N       tools.LocalizedFields `json:"-"`
}

// +transport
type LookupUserRequest struct {
	ID   string            `json:"id" urlenc:"id"`
	Lang jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`
}

// +transport
type DeleteUserRequest struct {
	ID string `json:"id"`
}

// +transport
type CreateVenueRequest struct {
	Name      jsval.MaybeString     `json:"name"`
	Address   jsval.MaybeString     `json:"address"`
	Longitude jsval.MaybeFloat      `json:"longitude,omitempty"`
	Latitude  jsval.MaybeFloat      `json:"latitude,omitempty"`
	L10N      tools.LocalizedFields `json:"-"`
}

// +transport
type UpdateVenueRequest struct {
	ID        string                `json:"id"`
	Name      jsval.MaybeString     `json:"name,omitempty"`
	Address   jsval.MaybeString     `json:"address,omitempty"`
	Longitude jsval.MaybeFloat      `json:"longitude,omitempty"`
	Latitude  jsval.MaybeFloat      `json:"latitude,omitempty"`
	L10N      tools.LocalizedFields `json:"-"`
}

// +transport
type DeleteVenueRequest struct {
	ID string `json:"id" urlenc:"id"`
}

// +transport
type ListVenueRequest struct {
	Since jsval.MaybeString `json:"since" urlenc:"since,omitempty,string"`
	Lang  jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`
	Limit jsval.MaybeInt    `json:"limit" urlenc:"limit,omitempty,int64"`
}

// +transport
type LookupVenueRequest struct {
	ID   string            `json:"id" urlenc:"id"`
	Lang jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`
}

// +transport
type ListSessionsByConferenceRequest struct {
	ConferenceID string            `json:"conference_id" urlenc:"conference_id"`
	Date         jsval.MaybeString `json:"date" urlenc:"date,omitempty,string"`
}
