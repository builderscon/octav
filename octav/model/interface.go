package model

import (
	"errors"
	"mime/multipart"
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
	tools.LocalizedFields `json:"-"`
	ID                    string              `json:"id"`
	Title                 string              `json:"title" l10n:"true"`
	Description           string              `json:"description,omitempty" l10n:"true"`
	SeriesID              string              `json:"series_id,omitempty"`
	Series                *ConferenceSeries   `json:"series,omitempty" decorate:"true"`
	SubTitle              string              `json:"sub_title" l10n:"true"`
	Slug                  string              `json:"slug"`
	Dates                 ConferenceDateList  `json:"dates,omitempty"`
	Administrators        UserList            `json:"administrators,omitempty" decorate:"true"`
	Venues                VenueList           `json:"venues,omitempty" decorate:"true"`
	FeaturedSpeakers      FeaturedSpeakerList `json:"featured_speakers,omitempty" decorate:"true"`
	Sponsors              SponsorList         `json:"sponsors,omitempty" decorate:"true"`
}
type ConferenceList []Conference

// +model
type ConferenceSeries struct {
	tools.LocalizedFields `json:"-"`
	ID                    string `json:"id"`
	Slug                  string `json:"slug"`
	Title                 string `json:"title" l10n:"true"`
}
type ConferenceSeriesList []ConferenceSeries

// +model
type Room struct {
	tools.LocalizedFields `json:"-"`
	ID                    string `json:"id"`
	VenueID               string `json:"venue_id"`
	Name                  string `json:"name" l10n:"true"`
	Capacity              uint   `json:"capacity"`
}
type RoomList []Room

// +model
type Session struct {
	tools.LocalizedFields `json:"-"`
	ID                    string      `json:"id"`
	ConferenceID          string      `json:"conference_id"`
	RoomID                string      `json:"room_id,omitempty"`
	SpeakerID             string      `json:"speaker_id"`
	Title                 string      `json:"title" l10n:"true"`
	Abstract              string      `json:"abstract" l10n:"true"`
	Memo                  string      `json:"memo"`
	StartsOn              time.Time   `json:"starts_on"`
	Duration              int         `json:"duration"`
	MaterialLevel         string      `json:"material_level"`
	Tags                  TagString   `json:"tags,omitempty" assign:"convert"`
	Category              string      `json:"category,omitempty"`
	SpokenLanguage        string      `json:"spoken_language,omitempty"`
	SlideLanguage         string      `json:"slide_language,omitempty"`
	SlideSubtitles        string      `json:"slide_subtitles,omitempty"`
	SlideURL              string      `json:"slide_url,omitempty"`
	VideoURL              string      `json:"video_url,omitempty"`
	PhotoPermission       string      `json:"photo_permission"`
	VideoPermission       string      `json:"video_permission"`
	SortOrder             int         `json:"-"`
	HasInterpretation     bool        `json:"has_interpretation"`
	Status                string      `json:"status"`
	Confirmed             bool        `json:"confirmed"`
	Conference            *Conference `json:"conference,omitempy" decorate:"true"` // only populated for JSON response
	Room                  *Room       `json:"room,omitempty" decorate:"true"`      // only populated for JSON response
	Speaker               *User       `json:"speaker,omitempty" decorate:"true"`   // only populated for JSON response
}
type SessionList []Session

type TagString string

// +model
type User struct {
	tools.LocalizedFields `json:"-"`
	ID                    string `json:"id"`
	AuthVia               string `json:"auth_via"`
	AuthUserID            string `json:"auth_user_id"`
	AvatarURL             string `json:"avatar_url,omitempty"`
	FirstName             string `json:"first_name,omitempty" l10n:"true"`
	LastName              string `json:"last_name,omitempty" l10n:"true"`
	Nickname              string `json:"nickname"`
	Email                 string `json:"email,omitempty"`
	TshirtSize            string `json:"tshirt_size,omitempty"`
	IsAdmin               bool   `json:"is_admin"`
}
type UserList []User

// +model
type Venue struct {
	tools.LocalizedFields `json:"-"`
	ID                    string   `json:"id,omitempty"`
	Name                  string   `json:"name" l10n:"true" decorate:"true"`
	Address               string   `json:"address" l10n:"true" decorate:"true"`
	Longitude             float64  `json:"longitude,omitempty"`
	Latitude              float64  `json:"latitude,omitempty"`
	Rooms                 RoomList `json:"rooms,omitempty"`
}
type VenueList []Venue

// +transport
type LookupConferenceSeriesRequest struct {
	ID   string            `json:"id"`
	Lang jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`
}

// +transport
type CreateConferenceSeriesRequest struct {
	UserID string                `json:"user_id"`
	Slug   string                `json:"slug"`
	Title  string                `json:"title"`
	L10N   tools.LocalizedFields `json:"-"`
}

// +transport
type UpdateConferenceSeriesRequest struct {
	ID    string                `json:"id"`
	Slug  jsval.MaybeString     `json:"slug"`
	Title jsval.MaybeString     `json:"title"`
	L10N  tools.LocalizedFields `json:"-"`
}

// +transport
type DeleteConferenceSeriesRequest struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
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
	UserID   string `json:"user_id"`  // ID of the operator
}

// +transport
type ListConferenceSeriesReponse []ConferenceSeries

// +transport
type CreateConferenceRequest struct {
	Title       string                `json:"title" l10n:"true"`
	Description jsval.MaybeString     `json:"description" l10n:"true"`
	SeriesID    string                `json:"series_id"`
	SubTitle    jsval.MaybeString     `json:"sub_title" l10n:"true"`
	Slug        string                `json:"slug"`
	UserID      string                `json:"user_id"`
	L10N        tools.LocalizedFields `json:"-"`
}

// +transport
type LookupConferenceRequest struct {
	ID   string            `json:"id" urlenc:"id"`
	Lang jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`
}

// +transport
type LookupConferenceBySlugRequest struct {
	Slug string            `json:"slug"`
	Lang jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`
}

// +transport
type UpdateConferenceRequest struct {
	ID          string            `json:"id"`
	Title       jsval.MaybeString `json:"title,omitempty" l10n:"true"`
	Description jsval.MaybeString `json:"description" l10n:"true"`
	SeriesID    jsval.MaybeString `json:"series_id,omitempty"`
	Slug        jsval.MaybeString `json:"slug,omitempty"`
	SubTitle    jsval.MaybeString `json:"sub_title,omitempty" l10n:"true"`
	Status      jsval.MaybeString `json:"status,omitempty"`
	UserID      string            `json:"user_id"`
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
	UserID       string             `json:"user_id"`
}

// +transport
type AddConferenceAdminRequest struct {
	ConferenceID string `json:"conference_id"`
	AdminID      string `json:"admin_id"`
	UserID       string `json:"user_id"`
}

// +transport
type AddConferenceVenueRequest struct {
	ConferenceID string `json:"conference_id"`
	VenueID      string `json:"venue_id"`
	UserID       string `json:"user_id"`
}

// +transport
type DeleteConferenceDatesRequest struct {
	ConferenceID string   `json:"conference_id"`
	Dates        DateList `json:"dates" extract:"true"`
	UserID       string   `json:"user_id"`
}

// +transport
type DeleteConferenceAdminRequest struct {
	ConferenceID string `json:"conference_id"`
	AdminID      string `json:"admin_id"`
	UserID       string `json:"user_id"`
}

// +transport
type DeleteConferenceVenueRequest struct {
	ConferenceID string `json:"conference_id"`
	VenueID      string `json:"venue_id"`
	UserID       string `json:"user_id"`
}

// +transport
type DeleteConferenceRequest struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

// +transport
type AddVenueRoomRequest struct {
	VenueID string `json:"venue_id"`
	RoomID  string `json:"room_id"`
}

// +transport
type DeleteVenueRoomRequest struct {
	VenueID string `json:"venue_id"`
	RoomID  string `json:"room_id"`
}

// +transport
type ListConferenceRequest struct {
	RangeEnd   jsval.MaybeString `json:"range_end,omitempty" urlenc:"range_end,omitempty,string"`
	RangeStart jsval.MaybeString `json:"range_start,omitempty" urlenc:"range_start,omitempty,string"`
	Since      jsval.MaybeString `json:"since,omitempty" urlenc:"since,omitempty,string"`
	Status     jsval.MaybeString `json:"status,omitempty" urlenc:"status,omitempty,string"`
	Lang       jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`
	Limit      jsval.MaybeInt    `json:"limit,omitempty" urlenc:"limit,omitempty,int64"`
}

// +transport
type ListConferenceReponse []Conference

// +transport
type CreateRoomRequest struct {
	VenueID  jsval.MaybeString     `json:"venue_id"`
	Name     jsval.MaybeString     `json:"name" l10n:"true"`
	Capacity jsval.MaybeUint       `json:"capacity"`
	L10N     tools.LocalizedFields `json:"-"`
	UserID   string                `json:"user_id"`
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
	UserID   string                `json:"user_id"`
}

// +transport
type DeleteRoomRequest struct {
	ID     string `json:"id" urlenc:"id"`
	UserID string `json:"user_id"`
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
	UserID          string                `json:"user_id"`
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
	UserID            string                `json:"user_id"`
}

// +transport
type DeleteSessionRequest struct {
	ID     string `json:"id" urlenc:"id"`
	UserID string `json:"user_id"`
}

// +transport
type CreateUserRequest struct {
	FirstName  jsval.MaybeString     `json:"first_name,omitempty" l18n:"true"`
	LastName   jsval.MaybeString     `json:"last_name,omitempty" l18n:"true"`
	Nickname   string                `json:"nickname"`
	Email      jsval.MaybeString     `json:"email,omitempty"`
	AuthVia    string                `json:"auth_via"`
	AuthUserID string                `json:"auth_user_id"`
	AvatarURL  jsval.MaybeString     `json:"avatar_url,omitempty"`
	TshirtSize jsval.MaybeString     `json:"tshirt_size,omitempty"`
	L10N       tools.LocalizedFields `json:"-"`
}

// +transport
type UpdateUserRequest struct {
	ID         string                `json:"id"`
	FirstName  jsval.MaybeString     `json:"first_name,omitempty"`
	LastName   jsval.MaybeString     `json:"last_name,omitempty"`
	Nickname   jsval.MaybeString     `json:"nickname,omitempty"`
	Email      jsval.MaybeString     `json:"email,omitempty"`
	AuthVia    jsval.MaybeString     `json:"auth_via,omitempty"`
	AuthUserID jsval.MaybeString     `json:"auth_user_id,omitempty"`
	AvatarURL  jsval.MaybeString     `json:"avatar_url,omitempty"`
	TshirtSize jsval.MaybeString     `json:"tshirt_size,omitempty"`
	UserID     string                `json:"user_id"`
	L10N       tools.LocalizedFields `json:"-"`
}

// +transport
type LookupUserRequest struct {
	ID   string            `json:"id" urlenc:"id"`
	Lang jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`
}

// +transport
type LookupUserByAuthUserIDRequest struct {
	AuthVia    string            `json:"auth_via" urlenc:"auth_via"`
	AuthUserID string            `json:"auth_user_id" urlenc:"auth_user_id"`
	Lang       jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`
}

// +transport
type DeleteUserRequest struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

// +transport
type ListUserRequest struct {
	Since jsval.MaybeString `json:"since" urlenc:"since,omitempty,string"`
	Lang  jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`
	Limit jsval.MaybeInt    `json:"limit" urlenc:"limit,omitempty,int64"`
}

// +transport
type CreateVenueRequest struct {
	Name      jsval.MaybeString     `json:"name"`
	Address   jsval.MaybeString     `json:"address"`
	Longitude jsval.MaybeFloat      `json:"longitude,omitempty"`
	Latitude  jsval.MaybeFloat      `json:"latitude,omitempty"`
	L10N      tools.LocalizedFields `json:"-"`
	UserID    string                `json:"user_id"`
}

// +transport
type UpdateVenueRequest struct {
	ID        string                `json:"id"`
	Name      jsval.MaybeString     `json:"name,omitempty"`
	Address   jsval.MaybeString     `json:"address,omitempty"`
	Longitude jsval.MaybeFloat      `json:"longitude,omitempty"`
	Latitude  jsval.MaybeFloat      `json:"latitude,omitempty"`
	L10N      tools.LocalizedFields `json:"-"`
	UserID    string                `json:"user_id"`
}

// +transport
type DeleteVenueRequest struct {
	ID     string `json:"id" urlenc:"id"`
	UserID string `json:"user_id"`
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
type ListSessionByConferenceRequest struct {
	ConferenceID string            `json:"conference_id" urlenc:"conference_id"`
	Date         jsval.MaybeString `json:"date" urlenc:"date,omitempty,string"`
	Lang         jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`
}

// +model
type Question struct {
	ID        string
	SessionID string
	UserID    string
	Body      string
}

// +transport
type LookupQuestionRequest struct {
	ID string `json:"id"`
}

// +transport
type CreateQuestionRequest struct {
	SessionID string `json:"session_id" urlenc:"session_id"`
	UserID    string `json:"user_id" urlenc:"user_id"`
	Body      string `json:"body" urlenc:"body"`
}

// +transport
type UpdateQuestionRequest struct {
	ID        string            `json:"id" urlenc:"id"`
	SessionID jsval.MaybeString `json:"session_id" urlenc:"session_id"`
	UserID    jsval.MaybeString `json:"user_id" urlenc:"user_id"`
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
	UserID             jsval.MaybeString `json:"user_id"`
	SessionID          jsval.MaybeString `json:"session_id"`
	UserPriorKnowledge int               `json:"user_prior_knowledge"`
	SpeakerKnowledge   int               `json:"speaker_knowledge"`
	MaterialQuality    int               `json:"material_quality"`
	OverallRating      int               `json:"overall_rating"`
	CommentGood        jsval.MaybeString `json:"comment_good" urlenc:"comment_good,omitempty,string"`
	CommentImprovement jsval.MaybeString `json:"comment_improvement" urlenc:"comment_improvement,omitempty,string"`
}

// +model
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
}

// +transport
type UpdateClientRequest struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
	Name   string `json:"name"`
}

// +model
type FeaturedSpeaker struct {
	tools.LocalizedFields `json:"-"`
	ID                    string `json:"id"`
	ConferenceID          string `json:"conference_id"`
	SpeakerID             string `json:"speaker_id"`
	AvatarURL             string `json:"avatar_url"`
	DisplayName           string `json:"display_name" l10n:"true"`
	Description           string `json:"description" l10n:"true"`
}
type FeaturedSpeakerList []FeaturedSpeaker

// +transport
type LookupFeaturedSpeakerRequest struct {
	ID   string            `json:"id"`
	Lang jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`
}

// +transport
type ListFeaturedSpeakersRequest struct {
	ConferenceID string            `json:"conference_id"`
	Since        jsval.MaybeString `json:"since" urlenc:"since,omitempty,string"`
	Lang         jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`
	Limit        jsval.MaybeInt    `json:"limit" urlenc:"limit,omitempty,int64"`
}

// +transport
type AddFeaturedSpeakerRequest struct {
	ConferenceID string                `json:"conference_id"`
	SpeakerID    jsval.MaybeString     `json:"speaker_id"`
	AvatarURL    jsval.MaybeString     `json:"avatar_url"`
	DisplayName  string                `json:"display_name" l18n:"true"`
	Description  string                `json:"description" l18n":"true"`
	L10N         tools.LocalizedFields `json:"-"`
	UserID       string                `json:"user_id"`
}
type CreateFeaturedSpeakerRequest struct {
	AddFeaturedSpeakerRequest
}

// +transport
type UpdateFeaturedSpeakerRequest struct {
	ID          string                `json:"id"`
	SpeakerID   jsval.MaybeString     `json:"speaker_id,omitempty"`
	AvatarURL   jsval.MaybeString     `json:"avatar_url,omitempty"`
	DisplayName jsval.MaybeString     `json:"display_name,omitempty" l18n:"true"`
	Description jsval.MaybeString     `json:"description,omitempty" l18n":"true"`
	L10N        tools.LocalizedFields `json:"-"`
	UserID      string                `json:"user_id"`
}

// +transport
type DeleteFeaturedSpeakerRequest struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

// +model
type Sponsor struct {
	tools.LocalizedFields `json:"-"`
	ID                    string `json:"id"`
	ConferenceID          string `json:"conference_id"`
	Name                  string `json:"name" l10n:"true"`
	LogoURL1              string `json:"logo_url1"`
	LogoURL2              string `json:"logo_url2,omitempty"`
	LogoURL3              string `json:"logo_url3,omitempty"`
	URL                   string `json:"url"`
	GroupName             string `json:"group_name"`
	SortOrder             int    `json:"sort_order"`
}
type SponsorList []Sponsor

// +transport
type LookupSponsorRequest struct {
	ID   string            `json:"id"`
	Lang jsval.MaybeString `json:"lang,omitempty" urlenc:"lang,omitempty,string"`
}

// +transport
type ListSponsorsRequest struct {
	ConferenceID string            `json:"conference_id"`
	GroupName    jsval.MaybeString `json:"group_name" urlenc:"group_name,omitempty,string"`
	Since        jsval.MaybeString `json:"since" urlenc:"since,omitempty,string"`
	Lang         jsval.MaybeString `json:"lang" urlenc:"lang,omitempty,string"`
	Limit        jsval.MaybeInt    `json:"limit" urlenc:"limit,omitempty,int64"`
}

// +transport
type AddSponsorRequest struct {
	ConferenceID  string                `json:"conference_id"`
	Name          string                `json:"name"`
	MultipartForm *multipart.Form       `json:"-"`
	LogoURL1      string                `json:"logo_url1"`
	LogoURL2      jsval.MaybeString     `json:"logo_url2,omitempty"`
	LogoURL3      jsval.MaybeString     `json:"logo_url3,omitempty"`
	URL           string                `json:"url"`
	GroupName     string                `json:"group_name"`
	SortOrder     int                   `json:"sort_order"`
	L10N          tools.LocalizedFields `json:"-"`
	UserID        string                `json:"user_id"`
}
type CreateSponsorRequest struct {
	AddSponsorRequest
}

// +transport
type UpdateSponsorRequest struct {
	ID        string                `json:"id"`
	Name      jsval.MaybeString     `json:"name,omitempty"`
	LogoURL1  jsval.MaybeString     `json:"logo_url1,omitempty"`
	LogoURL2  jsval.MaybeString     `json:"logo_url2,omitempty"`
	LogoURL3  jsval.MaybeString     `json:"logo_url3,omitempty"`
	URL       jsval.MaybeString     `json:"url,omitempty"`
	GroupName jsval.MaybeString     `json:"group_name,omitempty"`
	SortOrder jsval.MaybeInt        `json:"sort_order,omitempty"`
	L10N      tools.LocalizedFields `json:"-"`
	UserID    string                `json:"user_id"`
}

// +transport
type DeleteSponsorRequest struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}
