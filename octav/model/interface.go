package model

import (
	"errors"
	"mime/multipart"
	"sync"
	"time"

	"github.com/builderscon/octav/octav/db"
	"github.com/lestrrat/go-jsval"
)

const (
	StatusPending  = "pending"
	StatusAccepted = "accepted"
	StatusRejected = "rejected"
	StatusPublic   = "public"
	StatusPrivate  = "private"
)

var ErrInvalidConferenceHour = errors.New("invalid conference hour specification")

type ErrInvalidJSONFieldType struct {
	Field string
	Value interface{}
}

type ErrInvalidFieldType struct {
	Field string
}

// ObjectID is used to return the ID of a newly created object
type ObjectID struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// +model `UpdateRequest:"false"`
type Conference struct {
	LocalizedFields           `json:"-"`
	ID                        string               `json:"id"`
	Title                     string               `json:"title" l10n:"true"`
	Description               string               `json:"description,omitempty" l10n:"true"`
	CFPLeadText               string               `json:"cfp_lead_text,omitempty" l10n:"true"`
	CFPPreSubmitInstructions  string               `json:"cfp_pre_submit_instructions,omitempty" l10n:"true"`
	CFPPostSubmitInstructions string               `json:"cfp_post_submit_instructions,omitempty" l10n:"true"`
	ContactInformation        string               `json:"contact_information,omitempty" l10n:"true"`
	CoverURL                  string               `json:"cover_url"`
	RedirectURL               string               `json:"redirect_url"`
	SeriesID                  string               `json:"series_id,omitempty"`
	Series                    *ConferenceSeries    `json:"series,omitempty" decorate:"true"`
	SubTitle                  string               `json:"sub_title" l10n:"true"`
	Slug                      string               `json:"slug"`
	FullSlug                  string               `json:"full_slug,omitempty"` // Only populated when decorated
	Status                    string               `json:"status"`
	BlogFeedbackAvailable     bool                 `json:"blog_feedback_available"`
	TimetableAvailable        bool                 `json:"timetable_available"`
	Timezone                  string               `json:"timezone"`
	Dates                     ConferenceDateList   `json:"dates,omitempty"`
	Administrators            UserList             `json:"administrators,omitempty" decorate:"true"`
	Venues                    VenueList            `json:"venues,omitempty" decorate:"true"`
	FeaturedSpeakers          FeaturedSpeakerList  `json:"featured_speakers,omitempty" decorate:"true"`
	Sponsors                  SponsorList          `json:"sponsors,omitempty" decorate:"true"`
	SessionTypes              SessionTypeList      `json:"session_types,omitempty" decorate:"true"`
	Tracks                    TrackList            `json:"tracks,omitempty" decorate:"true"`
	ExternalResources         ExternalResourceList `json:"external_resources,omitempty"`
}
type ConferenceList []Conference

// +model
type ConferenceComponent struct {
	ID           string `json:"id"`
	ConferenceID string `json:"conference_id"`
	Name         string `json:"id"`
	Value        string `json:"value"`
}

// +transport
type LookupConferenceComponentRequest struct {
	ID   string            `json:"id" urlenc:"id"`
	Lang jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`

	VerifiedCall bool `json:"-"`
}

type CreateConferenceComponentRequest struct {
	ConferenceID string `json:"conference_id"`
	Name         string `json:"name"`
	Value        string `json:"value"`

	DatabaseOptions []db.InsertOption `json:"-"`
}

type UpdateConferenceComponentRequest struct {
	ID    string            `json:"id"`
	Name  jsval.MaybeString `json:"name"`
	Value jsval.MaybeString `json:"value"`
}

// +model
type ConferenceSeries struct {
	LocalizedFields `json:"-"`
	ID              string `json:"id"`
	Slug            string `json:"slug"`
	Title           string `json:"title" l10n:"true"`
}
type ConferenceSeriesList []ConferenceSeries

// +model
type Room struct {
	LocalizedFields `json:"-"`
	ID              string `json:"id"`
	VenueID         string `json:"venue_id"`
	Name            string `json:"name" l10n:"true"`
	Capacity        uint   `json:"capacity"`
}
type RoomList []Room

// +model
type SessionType struct {
	LocalizedFields       `json:"-"`
	ID                    string    `json:"id"`
	ConferenceID          string    `json:"conference_id"`
	Name                  string    `json:"name" l10n:"true"`
	Abstract              string    `json:"abstract" l10n:"true"`
	Duration              int       `json:"duration"`
	SubmissionStart       time.Time `json:"submission_start,omitempty"`
	SubmissionEnd         time.Time `json:"submission_end,omitempty"`
	IsDefault             bool      `json:"is_default"`
	IsAcceptingSubmission bool      `json:"is_accepting_submission"` // only used to return an easy flag to the client
}
type SessionTypeList []SessionType

// +transport
type LookupSessionTypeRequest struct {
	ID   string            `json:"id" urlenc:"id"`
	Lang jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`

	VerifiedCall bool `json:"-"`
}

// +transport
type AddSessionTypeRequest struct {
	ConferenceID    string            `json:"conference_id"`
	Name            string            `json:"name"`
	Abstract        string            `json:"abstract"`
	Duration        int               `json:"duration"`
	SubmissionStart jsval.MaybeString `json:"submission_start,omitempty"`
	SubmissionEnd   jsval.MaybeString `json:"submission_end,omitempty"`
	LocalizedFields `json:"-"`
	DatabaseOptions []db.InsertOption `json:"-"`
}
type CreateSessionTypeRequest struct {
	*AddSessionTypeRequest
}

// +transport
type DeleteSessionTypeRequest struct {
	ID     string `json:"id"`
}

// +transport
type ListSessionTypesByConferenceRequest struct {
	ConferenceID string            `json:"conference_id" urlenc:"conference_id"`
	Since        jsval.MaybeString `json:"since,omitempty" urlenc:"since,omitempty,string"`
	Limit        jsval.MaybeInt    `json:"limit,omitempty" urlenc:"limit,omitempty,int64"`
	Lang         jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`

	VerifiedCall bool `json:"-"`
}

// +transport
type UpdateSessionTypeRequest struct {
	ID              string            `json:"id"`
	Name            jsval.MaybeString `json:"name,omitempty"`
	Abstract        jsval.MaybeString `json:"abstract,omitempty"`
	Duration        jsval.MaybeInt    `json:"duration,omitempty"`
	IsDefault       jsval.MaybeBool   `json:"is_default,omitempty"`
	SubmissionStart jsval.MaybeString `json:"submission_start,omitempty"`
	SubmissionEnd   jsval.MaybeString `json:"submission_end,omitempty"`
	LocalizedFields `json:"-"`
}

// +model
type Session struct {
	LocalizedFields     `json:"-"`
	ID                  string       `json:"id"`
	ConferenceID        string       `json:"conference_id"`
	RoomID              string       `json:"room_id,omitempty"`
	SpeakerID           string       `json:"speaker_id"`
	SessionTypeID       string       `json:"session_type_id"`
	Title               string       `json:"title" l10n:"true"`
	Abstract            string       `json:"abstract" l10n:"true"`
	Memo                string       `json:"memo"`
	StartsOn            time.Time    `json:"starts_on,omitempty"`
	Duration            int          `json:"duration"`
	MaterialLevel       string       `json:"material_level"`
	Tags                TagString    `json:"tags,omitempty" assign:"convert"`
	Category            string       `json:"category,omitempty"`
	SelectionResultSent bool         `json:"selection_result_sent"`
	SpokenLanguage      string       `json:"spoken_language,omitempty"`
	SlideLanguage       string       `json:"slide_language,omitempty"`
	SlideSubtitles      string       `json:"slide_subtitles,omitempty"`
	SlideURL            string       `json:"slide_url,omitempty"`
	VideoURL            string       `json:"video_url,omitempty"`
	PhotoRelease        string       `json:"photo_release"`
	RecordingRelease    string       `json:"recording_release"`
	MaterialsRelease    string       `json:"materials_release"`
	SortOrder           int          `json:"-"`
	HasInterpretation   bool         `json:"has_interpretation"`
	Status              string       `json:"status"`
	Confirmed           bool         `json:"confirmed"`
	Room                *Room        `json:"room,omitempty" decorate:"true"`         // only populated for JSON response
	Speaker             *User        `json:"speaker,omitempty" decorate:"true"`      // only populated for JSON response
	SessionType         *SessionType `json:"session_type,omitempty" decorate:"true"` // only populated for JSON response
}
type SessionList []Session

type TagString string

// +model
type User struct {
	LocalizedFields `json:"-"`
	ID              string `json:"id"`
	AuthVia         string `json:"auth_via,omitempty"`
	AuthUserID      string `json:"auth_user_id,omitempty"`
	AvatarURL       string `json:"avatar_url,omitempty"`
	FirstName       string `json:"first_name,omitempty" l10n:"true"`
	LastName        string `json:"last_name,omitempty" l10n:"true"`
	Lang            string `json:"lang"`
	Nickname        string `json:"nickname"`
	Email           string `json:"email,omitempty"`
	TshirtSize      string `json:"tshirt_size,omitempty"`
	IsAdmin         bool   `json:"is_admin"`
	Timezone        string `json:"timezone"`
}
type UserList []User

// +model
type Venue struct {
	LocalizedFields `json:"-"`
	ID              string   `json:"id,omitempty"`
	Name            string   `json:"name" l10n:"true" decorate:"true"`
	Address         string   `json:"address" l10n:"true" decorate:"true"`
	PlaceID         string   `json:"place_id,omitempty"`
	URL             string   `json:"url,omitempty"`
	Longitude       float64  `json:"longitude,omitempty"`
	Latitude        float64  `json:"latitude,omitempty"`
	Rooms           RoomList `json:"rooms,omitempty"`
}
type VenueList []Venue

// +transport
type LookupConferenceSeriesRequest struct {
	ID   string            `json:"id"`
	Lang jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`

	VerifiedCall bool `json:"-"`
}

// +transport
type CreateConferenceSeriesRequest struct {
	SessionID       string `json:"sid"`
	Slug            string `json:"slug"`
	Title           string `json:"title"`
	LocalizedFields `json:"-"`
	DatabaseOptions []db.InsertOption `json:"-"`
}

// +transport
type UpdateConferenceSeriesRequest struct {
	ID              string            `json:"id"`
	Slug            jsval.MaybeString `json:"slug"`
	Title           jsval.MaybeString `json:"title"`
	LocalizedFields `json:"-"`
}

// +transport
type DeleteConferenceSeriesRequest struct {
	ID     string `json:"id"`
}

// +transport
type ListConferenceSeriesRequest struct {
	Since jsval.MaybeString `json:"since,omitempty" urlenc:"since,omitempty,string"`
	Limit jsval.MaybeInt    `json:"limit,omitempty" urlenc:"limit,omitempty,int64"`
	Lang  jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`
}

// +transport
type AddConferenceSeriesAdminRequest struct {
	SeriesID string `json:"series_id"`
	AdminID  string `json:"admin_id"` // new ID to add
}

// +transport
type ListConferenceSeriesReponse []ConferenceSeries

// +transport
type CreateConferenceRequest struct {
	Title                     string            `json:"title" l10n:"true"`
	CFPLeadText               jsval.MaybeString `json:"cfp_lead_text" l10n:"true"`
	CFPPreSubmitInstructions  jsval.MaybeString `json:"cfp_pre_submit_instructions" l10n:"true"`
	CFPPostSubmitInstructions jsval.MaybeString `json:"cfp_post_submit_instructions" l10n:"true"`
	ContactInformation        jsval.MaybeString `json:"contact_information" l10n:"true"`
	Description               jsval.MaybeString `json:"description" l10n:"true"`
	SeriesID                  string            `json:"series_id"`
	SubTitle                  jsval.MaybeString `json:"sub_title" l10n:"true"`
	Slug                      string            `json:"slug"`
	Timezone                  jsval.MaybeString `json:"timezone"`
	LocalizedFields           `json:"-"`
	DatabaseOptions           []db.InsertOption `json:"-"`
}

// +transport
type LookupConferenceRequest struct {
	ID   string            `json:"id" urlenc:"id"`
	Lang jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`

	VerifiedCall bool `json:"-"`
}

// +transport
type LookupConferenceBySlugRequest struct {
	Slug string            `json:"slug"`
	Lang jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`

	VerifiedCall bool `json:"-"`
}

// +transport
type UpdateConferenceRequest struct {
	ID                        string            `json:"id"`
	Title                     jsval.MaybeString `json:"title,omitempty" l10n:"true"`
	Description               jsval.MaybeString `json:"description" l10n:"true"`
	CFPLeadText               jsval.MaybeString `json:"cfp_lead_text" l10n:"true"`
	CFPPreSubmitInstructions  jsval.MaybeString `json:"cfp_pre_submit_instructions" l10n:"true"`
	CFPPostSubmitInstructions jsval.MaybeString `json:"cfp_post_submit_instructions" l10n:"true"`
	ContactInformation        jsval.MaybeString `json:"contact_information" l10n:"true"`
	MultipartForm             *multipart.Form   `json:"-"`
	RedirectURL               jsval.MaybeString `json:"redirect_url,omitempty"`
	SeriesID                  jsval.MaybeString `json:"series_id,omitempty"`
	Slug                      jsval.MaybeString `json:"slug,omitempty"`
	SubTitle                  jsval.MaybeString `json:"sub_title,omitempty" l10n:"true"`
	Status                    jsval.MaybeString `json:"status,omitempty"`
	BlogFeedbackAvailable     jsval.MaybeBool   `json:"blog_feedback_available,omitempty"`
	TimetableAvailable        jsval.MaybeBool   `json:"timetable_available,omitempty"`
	Timezone                  jsval.MaybeString `json:"timezone,omitempty"`
	LocalizedFields           `json:"-"`

	// These fields are only used internally
	CoverURL jsval.MaybeString `json:"-"`
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

// +model `LookupRequest:"false" UpdateRequest:"false"`
type ConferenceDate struct {
	ID    string
	Open  time.Time
	Close time.Time
}
type ConferenceDateList []ConferenceDate

// +transport
type ListConferenceDateRequest struct {
	ConferenceID string `json:"conference_id"`
}

// +transport
type CreateConferenceDateRequest struct {
	ConferenceID    string            `json:"conference_id"`
	Date            ConferenceDate    `json:"date" extract:"true"`
	DatabaseOptions []db.InsertOption `json:"-"`
}

// +transport
type ListConferenceAdminRequest struct {
	ConferenceID string            `json:"conference_id"`
	Lang         jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`
}

// +transport
type AddConferenceAdminRequest struct {
	ConferenceID string `json:"conference_id"`
	AdminID      string `json:"admin_id"`
}

// +transport
type AddConferenceVenueRequest struct {
	ConferenceID    string            `json:"conference_id"`
	VenueID         string            `json:"venue_id"`
	DatabaseOptions []db.InsertOption `json:"-"`
}

// +transport
type DeleteConferenceDateRequest struct {
	ConferenceID string `json:"conference_id"`
	Date         string `json:"date"`
}

// +transport
type DeleteConferenceAdminRequest struct {
	ConferenceID string `json:"conference_id"`
	AdminID      string `json:"admin_id"`
}

// +transport
type DeleteConferenceVenueRequest struct {
	ConferenceID string `json:"conference_id"`
	VenueID      string `json:"venue_id"`
}

// +transport
type DeleteConferenceRequest struct {
	ID     string `json:"id"`
}

// +transport
type AddVenueRoomRequest struct {
	VenueID         string            `json:"venue_id"`
	RoomID          string            `json:"room_id"`
	DatabaseOptions []db.InsertOption `json:"-"`
}

// +transport
type DeleteVenueRoomRequest struct {
	VenueID string `json:"venue_id"`
	RoomID  string `json:"room_id"`
}

// +transport
type ListConferenceRequest struct {
	Organizers []string          `json:"organizers" urlenc:"organizers,omitempty"`
	RangeEnd   jsval.MaybeString `json:"range_end,omitempty" urlenc:"range_end,omitempty,string"`
	RangeStart jsval.MaybeString `json:"range_start,omitempty" urlenc:"range_start,omitempty,string"`
	Since      jsval.MaybeString `json:"since,omitempty" urlenc:"since,omitempty,string"`
	Status     []string          `json:"status" urlenc:"status,omitempty"`
	Lang       jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`
	Limit      jsval.MaybeInt    `json:"limit,omitempty" urlenc:"limit,omitempty,int64"`
}

// +transport
type ListConferenceReponse []Conference

// +transport
type CreateRoomRequest struct {
	VenueID         jsval.MaybeString `json:"venue_id"`
	Name            jsval.MaybeString `json:"name" l10n:"true"`
	Capacity        jsval.MaybeUint   `json:"capacity"`
	LocalizedFields `json:"-"`
	DatabaseOptions []db.InsertOption `json:"-"`
}

// +transport
type LookupRoomRequest struct {
	ID   string            `json:"id" urlenc:"id"`
	Lang jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`

	VerifiedCall bool `json:"-"`
}

// +transport
type UpdateRoomRequest struct {
	ID              string            `json:"id"`
	VenueID         jsval.MaybeString `json:"venue_id,omitempty"`
	Name            jsval.MaybeString `json:"name,omitempty" l10n:"true"`
	Capacity        jsval.MaybeUint   `json:"capacity,omitempty"`
	LocalizedFields `json:"-"`
}

// +transport
type DeleteRoomRequest struct {
	ID     string `json:"id" urlenc:"id"`
}

// +transport
type ListRoomRequest struct {
	VenueID string            `json:"venue_id" urlenc:"venue_id"`
	Since   jsval.MaybeString `json:"since" urlenc:"since,omitempty,string"`
	Lang    jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`
	Limit   jsval.MaybeInt    `json:"limit,omitempty" urlenc:"limit,omitempty,int64"`

	VerifiedCall bool `json:"-"`
}

// +transport
type CreateSessionRequest struct {
	ConferenceID     string            `json:"conference_id"`
	SpeakerID        jsval.MaybeString `json:"speaker_id,omitempty"`
	SessionTypeID    string            `json:"session_type_id"`
	Title            jsval.MaybeString `json:"title,omitempty" l10n:"true"`
	Abstract         jsval.MaybeString `json:"abstract,omitempty" l10n:"true"`
	Memo             jsval.MaybeString `json:"memo,omitempty"`
	MaterialLevel    jsval.MaybeString `json:"material_level,omitempty"`
	Tags             jsval.MaybeString `json:"tags,omitempty"`
	Category         jsval.MaybeString `json:"category,omitempty"`
	SpokenLanguage   jsval.MaybeString `json:"spoken_language,omitempty"`
	SlideLanguage    jsval.MaybeString `json:"slide_language,omitempty"`
	SlideSubtitles   jsval.MaybeString `json:"slide_subtitles,omitempty"`
	SlideURL         jsval.MaybeString `json:"slide_url,omitempty"`
	VideoURL         jsval.MaybeString `json:"video_url,omitempty"`
	PhotoRelease     jsval.MaybeString `json:"photo_release,omitempty"`
	RecordingRelease jsval.MaybeString `json:"recording_release,omitempty"`
	MaterialsRelease jsval.MaybeString `json:"materials_release,omitempty"`
	LocalizedFields  `json:"-"`
	Duration         int               `json:"-"` // This is not sent from the client, but is used internally
	DatabaseOptions  []db.InsertOption `json:"-"`
}

// +transport
type LookupSessionRequest struct {
	ID   string            `json:"id" urlenc:"id"`
	Lang jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`

	VerifiedCall bool `json:"-"`
}

// +transport
type UpdateSessionRequest struct {
	ID                  string            `json:"id"`
	ConferenceID        jsval.MaybeString `json:"conference_id,omitempty"`
	SpeakerID           jsval.MaybeString `json:"speaker_id,omitempty"`
	SessionTypeID       jsval.MaybeString `json:"session_type_id,omitempty"`
	RoomID              jsval.MaybeString `json:"room_id,omitempty"`
	Title               jsval.MaybeString `json:"title,omitempty" l10n:"true"`
	Abstract            jsval.MaybeString `json:"abstract,omitempty" l10n:"true"`
	Memo                jsval.MaybeString `json:"memo,omitempty"`
	Duration            jsval.MaybeInt    `json:"duration,omitempty"`
	MaterialLevel       jsval.MaybeString `json:"material_level,omitempty"`
	Tags                jsval.MaybeString `json:"tags,omitempty"`
	Category            jsval.MaybeString `json:"category,omitempty"`
	SpokenLanguage      jsval.MaybeString `json:"spoken_language,omitempty"`
	SlideLanguage       jsval.MaybeString `json:"slide_language,omitempty"`
	SlideSubtitles      jsval.MaybeString `json:"slide_subtitles,omitempty"`
	SlideURL            jsval.MaybeString `json:"slide_url,omitempty"`
	VideoURL            jsval.MaybeString `json:"video_url,omitempty"`
	PhotoRelease        jsval.MaybeString `json:"photo_release,omitempty"`
	RecordingRelease    jsval.MaybeString `json:"recording_release,omitempty"`
	MaterialsRelease    jsval.MaybeString `json:"materials_release,omitempty"`
	SortOrder           jsval.MaybeInt    `json:"sort_order,omitempty"`
	HasInterpretation   jsval.MaybeBool   `json:"has_interpretation,omitempty"`
	SelectionResultSent jsval.MaybeBool   `json:"selection_result_sent,omitempty"`
	Status              jsval.MaybeString `json:"status,omitempty"`
	StartsOn            jsval.MaybeTime   `json:"starts_on,omitempty"`
	Confirmed           jsval.MaybeBool   `json:"confirmed,omitempty"`
	LocalizedFields     `json:"-"`
}

// +transport
type DeleteSessionRequest struct {
	ID     string `json:"id" urlenc:"id"`
}

// +transport
type CreateUserRequest struct {
	FirstName       jsval.MaybeString `json:"first_name,omitempty" l10n:"true"`
	LastName        jsval.MaybeString `json:"last_name,omitempty" l10n:"true"`
	Nickname        string            `json:"nickname"`
	Email           jsval.MaybeString `json:"email,omitempty"`
	AuthVia         string            `json:"auth_via"`
	AuthUserID      string            `json:"auth_user_id"`
	AvatarURL       jsval.MaybeString `json:"avatar_url,omitempty"`
	TshirtSize      jsval.MaybeString `json:"tshirt_size,omitempty"`
	LocalizedFields `json:"-"`
	DatabaseOptions []db.InsertOption `json:"-"`
}

// +transport
type UpdateUserRequest struct {
	ID              string            `json:"id"`
	FirstName       jsval.MaybeString `json:"first_name,omitempty" l10n:"true"`
	LastName        jsval.MaybeString `json:"last_name,omitempty" l10n:"true"`
	Nickname        jsval.MaybeString `json:"nickname,omitempty"`
	Lang            jsval.MaybeString `json:"lang,omitempty"`
	Email           jsval.MaybeString `json:"email,omitempty"`
	AuthVia         jsval.MaybeString `json:"auth_via,omitempty"`
	AuthUserID      jsval.MaybeString `json:"auth_user_id,omitempty"`
	AvatarURL       jsval.MaybeString `json:"avatar_url,omitempty"`
	TshirtSize      jsval.MaybeString `json:"tshirt_size,omitempty"`
	LocalizedFields `json:"-"`
}

// +transport
type LookupUserRequest struct {
	ID   string            `json:"id" urlenc:"id"`
	Lang jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`

	VerifiedCall bool `json:"-"`
}

// +transport
type LookupUserByAuthUserIDRequest struct {
	AuthVia    string            `json:"auth_via" urlenc:"auth_via"`
	AuthUserID string            `json:"auth_user_id" urlenc:"auth_user_id"`
	Lang       jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`

	VerifiedCall bool `json:"-"`
}

// +transport
type DeleteUserRequest struct {
	ID     string `json:"id"`
}

// +transport
type ListUserRequest struct {
	Pattern jsval.MaybeString `json:"pattern" urlenc:"pattern,omitempty,string"`
	Since   jsval.MaybeString `json:"since" urlenc:"since,omitempty,string"`
	Lang    jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`
	Limit   jsval.MaybeInt    `json:"limit" urlenc:"limit,omitempty,int64"`

	VerifiedCall bool `json:"-"`
}

// +transport
type CreateVenueRequest struct {
	Name            jsval.MaybeString `json:"name" l10n:"true"`
	Address         jsval.MaybeString `json:"address" l10n:"true"`
	Longitude       jsval.MaybeFloat  `json:"longitude,omitempty"`
	Latitude        jsval.MaybeFloat  `json:"latitude,omitempty"`
	LocalizedFields `json:"-"`
	DatabaseOptions []db.InsertOption `json:"-"`
}

// +transport
type UpdateVenueRequest struct {
	ID              string            `json:"id"`
	Name            jsval.MaybeString `json:"name,omitempty" l10n:"true"`
	Address         jsval.MaybeString `json:"address,omitempty" l10n:"true"`
	Longitude       jsval.MaybeFloat  `json:"longitude,omitempty"`
	Latitude        jsval.MaybeFloat  `json:"latitude,omitempty"`
	LocalizedFields `json:"-"`
}

// +transport
type DeleteVenueRequest struct {
	ID     string `json:"id" urlenc:"id"`
}

// +transport
type ListVenueRequest struct {
	Since jsval.MaybeString `json:"since" urlenc:"since,omitempty,string"`
	Lang  jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`
	Limit jsval.MaybeInt    `json:"limit" urlenc:"limit,omitempty,int64"`

	VerifiedCall bool `json:"-"`
}

// +transport
type LookupVenueRequest struct {
	ID   string            `json:"id" urlenc:"id"`
	Lang jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`

	VerifiedCall bool `json:"-"`
}

// +transport
type ListSessionsRequest struct {
	ConferenceID jsval.MaybeString `json:"conference_id" urlenc:"conference_id,omitempty,string"`
	Confirmed    []bool            `json:"confirmed" urlenc:"confirmed,omitempty"`
	RangeStart   jsval.MaybeString `json:"range_start" urlenc:"range_start,omitempty,string"`
	RangeEnd     jsval.MaybeString `json:"range_end" urlenc:"range_end,omitempty,string"`
	SpeakerID    jsval.MaybeString `json:"speaker_id" urlenc:"speaker_id,omitempty,string"`
	Status       []string          `json:"status" urlenc:"status,omitempty"`
	Lang         jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`

	VerifiedCall bool `json:"-"`
}

// +model
type Question struct {
	ID        string
	SessionID string
	Body      string
}

// +transport
type LookupQuestionRequest struct {
	ID string `json:"id"`

	VerifiedCall bool `json:"-"`
}

// +transport
type CreateQuestionRequest struct {
	SessionID       string            `json:"session_id" urlenc:"session_id"`
	Body            string            `json:"body" urlenc:"body"`
	DatabaseOptions []db.InsertOption `json:"-"`
}

// +transport
type UpdateQuestionRequest struct {
	ID        string            `json:"id" urlenc:"id"`
	SessionID jsval.MaybeString `json:"session_id" urlenc:"session_id"`
	Body      jsval.MaybeString `json:"body" urlenc:"body"`
}

// +transport
type DeleteQuestionRequest struct {
	ID string `json:"id" urlenc:"id"`
}

// +transport
type ListQuestionRequest struct {
	Since jsval.MaybeString `json:"since" urlenc:"since,omitempty,string"`
	Lang  jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`
	Limit jsval.MaybeInt    `json:"limit" urlenc:"limit,omitempty,int64"`
}

// +transport
type CreateSessionSurveyResponseRequest struct {
	SessionID          jsval.MaybeString `json:"session_id"`
	UserPriorKnowledge int               `json:"user_prior_knowledge"`
	SpeakerKnowledge   int               `json:"speaker_knowledge"`
	MaterialQuality    int               `json:"material_quality"`
	OverallRating      int               `json:"overall_rating"`
	CommentGood        jsval.MaybeString `json:"comment_good" urlenc:"comment_good,omitempty,string"`
	CommentImprovement jsval.MaybeString `json:"comment_improvement" urlenc:"comment_improvement,omitempty,string"`
}

// +model `CreateRequest:"false" UpdateRequest:"false"`
type Client struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
	Name   string `json:"name"`
}

// +transport
type CreateClientRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// +transport
type LookupClientRequest struct {
	ID string `json:"id"`

	VerifiedCall bool `json:"-"`
}

// +transport
type UpdateClientRequest struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
	Name   string `json:"name"`
}

// +model
type FeaturedSpeaker struct {
	LocalizedFields `json:"-"`
	ID              string `json:"id"`
	ConferenceID    string `json:"conference_id"`
	SpeakerID       string `json:"speaker_id"`
	AvatarURL       string `json:"avatar_url"`
	DisplayName     string `json:"display_name" l10n:"true"`
	Description     string `json:"description" l10n:"true"`
}
type FeaturedSpeakerList []FeaturedSpeaker

// +transport
type LookupFeaturedSpeakerRequest struct {
	ID   string            `json:"id"`
	Lang jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`

	VerifiedCall bool `json:"-"`
}

// +transport
type ListFeaturedSpeakersRequest struct {
	ConferenceID string            `json:"conference_id"`
	Since        jsval.MaybeString `json:"since" urlenc:"since,omitempty,string"`
	Lang         jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`
	Limit        jsval.MaybeInt    `json:"limit" urlenc:"limit,omitempty,int64"`

	VerifiedCall bool `json:"-"`
}

// +transport
type AddFeaturedSpeakerRequest struct {
	ConferenceID    string            `json:"conference_id"`
	SpeakerID       jsval.MaybeString `json:"speaker_id"`
	AvatarURL       jsval.MaybeString `json:"avatar_url"`
	DisplayName     string            `json:"display_name" l10n:"true"`
	Description     string            `json:"description" l10n:"true"`
	LocalizedFields `json:"-"`
	DatabaseOptions []db.InsertOption `json:"-"`
}
type CreateFeaturedSpeakerRequest struct {
	*AddFeaturedSpeakerRequest
}

// +transport
type UpdateFeaturedSpeakerRequest struct {
	ID              string            `json:"id"`
	SpeakerID       jsval.MaybeString `json:"speaker_id,omitempty"`
	AvatarURL       jsval.MaybeString `json:"avatar_url,omitempty"`
	DisplayName     jsval.MaybeString `json:"display_name,omitempty" l10n:"true"`
	Description     jsval.MaybeString `json:"description,omitempty" l10n:"true"`
	LocalizedFields `json:"-"`
}

// +transport
type DeleteFeaturedSpeakerRequest struct {
	ID     string `json:"id"`
}

// +model
type Sponsor struct {
	LocalizedFields `json:"-"`
	ID              string `json:"id"`
	ConferenceID    string `json:"conference_id"`
	Name            string `json:"name" l10n:"true"`
	LogoURL         string `json:"logo_url,omitempty"`
	URL             string `json:"url"`
	GroupName       string `json:"group_name"`
	SortOrder       int    `json:"sort_order"`
}
type SponsorList []Sponsor

// +transport
type LookupSponsorRequest struct {
	ID   string            `json:"id"`
	Lang jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`

	VerifiedCall bool `json:"-"`
}

// +transport
type ListSponsorsRequest struct {
	ConferenceID string            `json:"conference_id"`
	GroupName    jsval.MaybeString `json:"group_name" urlenc:"group_name,omitempty,string"`
	Since        jsval.MaybeString `json:"since" urlenc:"since,omitempty,string"`
	Lang         jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`
	Limit        jsval.MaybeInt    `json:"limit" urlenc:"limit,omitempty,int64"`

	VerifiedCall bool `json:"-"`
}

// +transport
type AddSponsorRequest struct {
	ConferenceID    string            `json:"conference_id"`
	Name            string            `json:"name" l10n:"true"`
	URL             string            `json:"url"`
	LogoURL         jsval.MaybeString `json:"logo_url"`
	GroupName       string            `json:"group_name"`
	SortOrder       int               `json:"sort_order"`
	LocalizedFields `json:"-"`
	DatabaseOptions []db.InsertOption `json:"-"`
}
type CreateSponsorRequest struct {
	*AddSponsorRequest
}

// +transport
type UpdateSponsorRequest struct {
	ID              string            `json:"id"`
	Name            jsval.MaybeString `json:"name,omitempty" l10n:"true"`
	URL             jsval.MaybeString `json:"url,omitempty"`
	GroupName       jsval.MaybeString `json:"group_name,omitempty"`
	MultipartForm   *multipart.Form   `json:"-"`
	SortOrder       jsval.MaybeInt    `json:"sort_order,omitempty"`
	LocalizedFields `json:"-"`
	LogoURL         jsval.MaybeString `json:"-"`
}

// +transport
type DeleteSponsorRequest struct {
	ID     string `json:"id"`
}

// +transport
type ListConferencesByOrganizerRequest struct {
	OrganizerID []string          `json:"organizer_id" urlenc:"organizer_id,omitempty"`
	Since       jsval.MaybeString `json:"since" urlenc:"since,omitempty,string"`
	Status      []string          `json:"status" urlenc:"status,omitempty"`
	Lang        jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`
	Limit       jsval.MaybeInt    `json:"limit" urlenc:"limit,omitempty,int64"`
}

type LocalizedFields struct {
	lock sync.RWMutex
	// Language -> field/value
	fields map[string]map[string]string
}

// +transport
type CreateTemporaryEmailRequest struct {
	TargetID string            `json:"target_id"` // ID of the user to register the email for
	Email    string            `json:"email"`
	Lang     jsval.MaybeString `json:"lang"`
}

// +transport
type CreateTemporaryEmailResponse struct {
	ConfirmationKey string `json:"confirmation_key,omitempty"`
}

// +transport
type ConfirmTemporaryEmailRequest struct {
	TargetID        string `json:"target_id"` // ID of the user to register the email for
	ConfirmationKey string `json:"confirmation_key"`
}

// +transport
type ListConferenceCredentialRequest struct {
	ConferenceID string `json:"conference_id"`
}

// +transport
type AddConferenceCredentialRequest struct {
	ConferenceID string `json:"conference_id"`
	Type         string `json:"type"`
	Data         string `json:"data"`
}

// +transport
type TweetAsConferenceRequest struct {
	ConferenceID string `json:"conference_id"`
	Tweet        string `json:"tweet"`
}

type JSONTime time.Time
type JSONTimeList []JSONTime

// +transport
type GetConferenceScheduleRequest struct {
	ConferenceID string            `json:"conference_id"`
	Lang         jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`
}

// +transport
type VerifyUserRequest struct {
	ID     string `json:"id"`      // ID of the user being verified
}

// +transport
type SendSelectionResultNotificationRequest struct {
	Force     bool   `json:"force"`      // true to force sending notification after we have already done so for this session
	SessionID string `json:"session_id"` // ID of the session to which we're making the notification

	VerifiedCall bool `json:"-"`
}

// +transport
type SendAllSelectionResultNotificationRequest struct {
	Force        bool   `json:"force"`         // true to force sending notification after we have already done so for this session
	ConferenceID string `json:"conference_id"` // ID of the conference to which we're making the notification

	VerifiedCall bool `json:"-"`
}

// +transport
type SendSelectionResultNotificationResponse struct {
	Message string `json:"message"`
}

// +transport
type SendAllSelectionResultNotificationResponse struct {
	Message string `json:"message"`
}

// +model `CreateRequest:"false" UpdateRequest:"false" Update:"false" Delete:"false" Lookup:"false" LookupRequest:"false"`
type LocalizedString struct {
	ParentID   string // EID of the parent object
	ParentType string // Type of the parent object
	Name       string
	Language   string
	Localized  string
}

// +model
type Track struct {
	LocalizedFields `json:"-"`
	ID              string `json:"id"`
	ConferenceID    string `json:"conference_id,omitempty"`
	RoomID          string `json:"room_id"`
	SortOrder       int    `json:"sort_order"`
	Name            string `json:"name" l10n:"true"`
}
type TrackList []Track

// +transport
type LookupTrackRequest struct {
	ID          string            `json:"id"`
	Lang        jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`
	VerifiedCall bool              `json:"-"`
}

// +transport
type CreateTrackRequest struct {
	ConferenceID    string            `json:"conference_id"`
	RoomID          string            `json:"room_id"`
	Name            jsval.MaybeString `json:"name,omitempty" l10n:"true"`
	SortOrder       jsval.MaybeInt    `json:"sort_order,omitempty"`
	LocalizedFields `json:"-"`
	DatabaseOptions []db.InsertOption `json:"-"`
}

// +transport
type UpdateTrackRequest struct {
	ID              string            `json:"id"`
	Name            jsval.MaybeString `json:"name,omitempty" l10n:"true"`
	RoomID          jsval.MaybeString `json:"room_id,omitempty"`
	SortOrder       jsval.MaybeInt    `json:"sort_order,omitempty"`
	LocalizedFields `json:"-"`
}

// +transport
type DeleteTrackRequest struct {
	ID              string `json:"id"`
	LocalizedFields `json:"-"`
}

// +model
type ConferenceVenue struct {
	ConferenceID string
	VenueID      string
}

// +transport
type CreateConferenceVenueRequest struct {
	ConferenceID    string
	VenueID         string
	DatabaseOptions []db.InsertOption `json:"-"`
}

// +transport
type UpdateConferenceVenueRequest struct {
	ConferenceID string
	VenueID      string
}

// +model
type BlogEntry struct {
	ID           string `json:"id"`
	ConferenceID string `json:"conference_id,omitempty"`
	Status       string `json:"status,omitempty"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	URLHash      string `json:"status,omitempty"`
}
type BlogEntryList []BlogEntry

// +transport
type CreateBlogEntryRequest struct {
	ConferenceID    string            `json:"conference_id"`
	Status          string            `json:"status"`
	Title           string            `json:"title"`
	URL             string            `json:"url"`
	DatabaseOptions []db.InsertOption `json:"-"`
}

// +transport
type UpdateBlogEntryRequest struct {
	ID     string            `json:"id"`
	Status jsval.MaybeString `json:"status,omitempty"`
	Title  jsval.MaybeString `json:"title,omitempty"`
	URL    jsval.MaybeString `json:"url,omitempty"`
}

// +transport
type LookupBlogEntryRequest struct {
	ID     string `json:"id"`
}

// +transport
type DeleteBlogEntryRequest struct {
	ID     string `json:"id"`
}

// +transport
type ListBlogEntriesRequest struct {
	ConferenceID string            `json:"conference_id"`
	Status       []string          `json:"status" urlenc:"status,omitempty"`
	Lang         jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`
	VerifiedCall  bool              `json:"-"`
}

// +transport
type AddConferenceStaffRequest struct {
	ConferenceID string `json:"conference_id"`
	StaffID      string `json:"staff_id"`
}

// +transport
type DeleteConferenceStaffRequest struct {
	ConferenceID string `json:"conference_id"`
	StaffID      string `json:"staff_id"`
}

// +transport
type ListConferenceStaffRequest struct {
	ConferenceID string            `json:"conference_id"`
	Lang         jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`
}

// +model
type ExternalResource struct {
	LocalizedFields `json:"-"`
	ID              string `json:"id"`
	ConferenceID    string `json:"conference_id"`
	Description     string `json:"description" l10n:"true"`
	ImageURL        string `json:"image_url"`
	Title           string `json:"title" l10n:"true"`
	URL             string `json:"url"`
	SortOrder       int    `json:"sort_order"`
}
type ExternalResourceList []ExternalResource

// +transport
type LookupExternalResourceRequest struct {
	ID    string            `json:"id"`
	Since jsval.MaybeString `json:"since,omitempty" urlenc:"since,omitempty,string"`
	Limit jsval.MaybeInt    `json:"limit,omitempty" urlenc:"limit,omitempty,int64"`
	Lang  jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`

	VerifiedCall bool `json:"-"`
}

// +transport
type ListExternalResourceRequest struct {
	ConferenceID string            `json:"conference_id"`
	Since        jsval.MaybeString `json:"since,omitempty" urlenc:"since,omitempty,string"`
	Limit        jsval.MaybeInt    `json:"limit,omitempty" urlenc:"limit,omitempty,int64"`
	Lang         jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`

	VerifiedCall bool `json:"-"`
}

// +transport
type CreateExternalResourceRequest struct {
	ConferenceID    string            `json:"conference_id"`
	Description     jsval.MaybeString `json:"description,omitempty" l10n:"true"`
	Title           string            `json:"title" l10n:"true"`
	URL             string            `json:"url"`
	LocalizedFields `json:"-"`
	SortOrder       jsval.MaybeInt    `json:"sort_order,omitempty"`
	DatabaseOptions []db.InsertOption `json:"-"`
}

// +transport
type UpdateExternalResourceRequest struct {
	ID              string            `json:"id"`
	Description     jsval.MaybeString `json:"description,omitempty" l10n:"true"`
	Title           jsval.MaybeString `json:"title,omitempty" l10n:"true"`
	URL             jsval.MaybeString `json:"url,omitempty"`
	LocalizedFields `json:"-"`
	SortOrder       jsval.MaybeInt `json:"sort_order,omitempty"`
}

// +transport
type DeleteExternalResourceRequest struct {
	ID     string `json:"id"`
}

// +transport
type SetSessionVideoCoverRequest struct {
	ID            string          `json:"id"`
	MultipartForm *multipart.Form `json:"-"`
}

// +transport
type CreateClientSessionRequest struct {
	AccessToken string `json:"access_token"`
	AuthVia     string `json:"auth_via"`
}

// +transport
type CreateClientSessionResponse struct {
	SessionID string `json:"sid"`
	Expires   string `json:"expires"`
}

// +transport
type LookupUserAvatarRequest struct {
	ID string `json:"id"`
}
