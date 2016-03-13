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
type CreateConferenceRequest struct {
	Title    string                `json:"title" l10n:"true"`
	SubTitle jsval.MaybeString     `json:"sub_title" l10n:"true"`
	Slug     string                `json:"slug"`
	L10N     tools.LocalizedFields `json:"-"`
}
type LookupConferenceRequest struct {
	ID   string            `json:"id" urlenc:"id"`
	Lang jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`
}
type UpdateConferenceRequest struct {
	ID       string            `json:"id"`
	Title    jsval.MaybeString `json:"title,omitempty" l10n:"true"`
	SubTitle jsval.MaybeString `json:"sub_title,omitempty" l10n:"true"`
	Slug     jsval.MaybeString `json:"slug,omitempty"`
	// TODO dates
	L10N tools.LocalizedFields `json:"-"`
}

type Room struct{}
type CreateRoomRequest struct {
	VenueID  jsval.MaybeString     `json:"venue_id"`
	Name     jsval.MaybeString     `json:"name" l10n:"true"`
	Capacity jsval.MaybeUint       `json:"capacity"`
	L10N     tools.LocalizedFields `json:"-"`
}
type UpdateRoomRequest struct {
	ID       string                `json:"id"`
	VenueID  jsval.MaybeString     `json:"venue_id,omitempty"`
	Name     jsval.MaybeString     `json:"name,omitempty" l10n:"true"`
	Capacity jsval.MaybeUint       `json:"capacity,omitempty"`
	L10N     tools.LocalizedFields `json:"-"`
}

type Session struct{}
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
type DeleteSessionRequest struct {
	ID string `json:"id" urlenc:"id"`
}

type User struct{}
type CreateUserRequest struct {
	FirstName  string                `json:"first_name"`
	LastName   string                `json:"last_name"`
	Nickname   string                `json:"nickname"`
	Email      string                `json:"email"`
	TshirtSize string                `json:"tshirt_size"`
	L10N       tools.LocalizedFields `json:"-"`
}
type UpdateUserRequest struct{}

type Venue struct{}
type CreateVenueRequest struct {
	Name      jsval.MaybeString     `json:"name"`
	Address   jsval.MaybeString     `json:"address"`
	Longitude jsval.MaybeFloat      `json:"longitude,omitempty"`
	Latitude  jsval.MaybeFloat      `json:"latitude,omitempty"`
	L10N      tools.LocalizedFields `json:"-"`
}
type UpdateVenueRequest struct {
	ID        jsval.MaybeString     `json:"id"`
	Name      jsval.MaybeString     `json:"name"`
	Address   jsval.MaybeString     `json:"address"`
	Longitude jsval.MaybeFloat      `json:"longitude,omitempty"`
	Latitude  jsval.MaybeFloat      `json:"latitude,omitempty"`
	L10N      tools.LocalizedFields `json:"-"`
}

type DeleteRoomRequest struct {
	ID string `json:"id" urlenc:"id"`
}
type ListRoomRequest struct {
	VenueID string            `json:"venue_id" urlenc:"venue_id"`
	Since   jsval.MaybeString `json:"since" urlenc:"since,omitempty,string"`
	Lang    jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`
	Limit   jsval.MaybeInt    `json:"limit,omitempty" urlenc:"limit,omitempty,int64"`
}
type LookupRoomRequest struct {
	ID   string            `json:"id" urlenc:"id"`
	Lang jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`
}
type LookupSessionRequest struct {
	ID string `json:"id" urlenc:"id"`
	Lang jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`
}

type LookupUserRequest struct {
	ID string `json:"id" urlenc:"id"`
}
type DeleteUserRequest struct {
	ID string `json:"id" urlenc:"id"`
}

type DeleteVenueRequest struct {
	ID string `json:"id" urlenc:"id"`
}
type ListVenueRequest struct {
	Since jsval.MaybeString `json:"since" urlenc:"since,omitempty,string"`
	Lang  jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`
	Limit jsval.MaybeInt    `json:"limit" urlenc:"limit,omitempty,int64"`
}
type LookupVenueRequest struct {
	ID string `json:"id" urlenc:"id"`
}

type DeleteConferenceRequest struct {
	ID string `json:"id" urlenc:"id"`
}
type ListSessionsByConferenceRequest struct {
	ConferenceID string            `json:"conference_id" urlenc:"conference_id"`
	Date         jsval.MaybeString `json:"date" urlenc:"date,omitempty,string"`
}
