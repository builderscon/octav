package service

import (
	"github.com/builderscon/octav/octav/tools"
	"github.com/lestrrat/go-jsval"
)

type ErrInvalidJSONFieldType struct {
	Field string
}

type ErrInvalidFieldType struct {
	Field string
}

type Conference struct{}

// +transport
type CreateConferenceRequest struct {
	Title    string                `json:"title" l10n:"true"`
	SubTitle jsval.MaybeString     `json:"sub_title" l10n:"true"`
	Slug     string                `json:"slug"`
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

// +transport
type DeleteConferenceRequest struct {
	ID string `json:"id" urlenc:"id"`
}

type Room struct{}

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

type Session struct{}

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

type User struct{}

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
	ID        string                `json:"id"`
	FirstName jsval.MaybeString     `json:"first_name,omitempty"`
	LastName  jsval.MaybeString     `json:"last_name,omitempty"`
	Nickname  jsval.MaybeString     `json:"nickname,omitempty"`
	Email     jsval.MaybeString     `json:"email,omitempty"`
	L10N      tools.LocalizedFields `json:"-"`
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

type Venue struct{}

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
