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

type Conference struct {}
type CreateConferenceRequest struct {
	Title    string            `json:"title"`
	SubTitle jsval.MaybeString `json:"sub_title"`
	Slug     string            `json:"slug"`
	L10N     tools.LocalizedFields   `json:"-"`
}

type Session struct {}
type CreateSessionRequest struct {
	ConferenceID    jsval.MaybeString `json:"conference_id,omitempty"`
	SpeakerID       jsval.MaybeString `json:"speaker_id,omitempty"`
	Title           jsval.MaybeString `json:"title,omitempty"`
	Abstract        jsval.MaybeString `json:"abstract,omitempty"`
	Memo            jsval.MaybeString `json:"memo,omitempty"`
	Duration        jsval.MaybeInt    `json:"duration,omitempty"`
	MaterialLevel   jsval.MaybeString `json:"material_level,omitempty"`
	Tags            jsval.MaybeString `json:"tags,omitempty"`
	Category        jsval.MaybeString `json:"category,omitempty"`
	SpokenLanguage  jsval.MaybeString `json:"spoken_language,omitempty"`
	SlideLanguage   jsval.MaybeString `json:"slide_language,omitempty"`
	SlideSubtitles  jsval.MaybeString `json:"slide_subtitles,omitempty"`
	SlideURL        jsval.MaybeString `json:"slide_url,omitempty"`
	VideoURL        jsval.MaybeString `json:"video_url,omitempty"`
	PhotoPermission jsval.MaybeString `json:"photo_permission,omitempty"`
	VideoPermission jsval.MaybeString `json:"video_permission,omitempty"`
	L10N            tools.LocalizedFields   `json:"-"`
}


