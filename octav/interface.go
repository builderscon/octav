package octav

import "time"

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
	ID       string `json:"id"`
	VenueID  string `json:"venue_id"`
	Name     string `json:"name"`
	Capacity uint   `json:"capacity"`
}
type RoomList []Room
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
type User struct{}

type VenueList []Venue
type Venue struct {
	ID   string
	Name string
}

type ConferenceList []Conference
type Conference struct {
	ID       string           `json:"id"`
	Title    string           `json:"title"`
	SubTitle string           `json:"subtitle"`
	Slug     string           `json:"slug"`
	Dates    []ConferenceDate `json:"dates"` // only populated for JSON response
}
