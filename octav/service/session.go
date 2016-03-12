package service

import (
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/tools"
)

func (v *Session) Create(tx *db.Tx, payload CreateSessionRequest, vdb *db.Session) error {
	vdb.EID = tools.UUID()
	vdb.ConferenceID = payload.ConferenceID.String
	vdb.SpeakerID = payload.SpeakerID.String
	vdb.Duration = int(payload.Duration.Int)
	vdb.HasInterpretation = false
	vdb.Status = "pending"
	vdb.SortOrder = 0
	vdb.Confirmed = false

	if payload.Title.Valid() {
		vdb.Title.Valid = true
		vdb.Title.String = payload.Title.String
	}

	if payload.Abstract.Valid() {
		vdb.Abstract.Valid = true
		vdb.Abstract.String = payload.Abstract.String
	}

	if payload.Memo.Valid() {
		vdb.Memo.Valid = true
		vdb.Memo.String = payload.Memo.String
	}

	if payload.MaterialLevel.Valid() {
		vdb.MaterialLevel.Valid = true
		vdb.MaterialLevel.String = payload.MaterialLevel.String
	}

	if payload.Category.Valid() {
		vdb.Category.Valid = true
		vdb.Category.String = payload.Category.String
	}

	if payload.SpokenLanguage.Valid() {
		vdb.SpokenLanguage.Valid = true
		vdb.SpokenLanguage.String = payload.SpokenLanguage.String
	}

	if payload.SlideLanguage.Valid() {
		vdb.SlideLanguage.Valid = true
		vdb.SlideLanguage.String = payload.SlideLanguage.String
	}

	if payload.SlideSubtitles.Valid() {
		vdb.SlideSubtitles.Valid = true
		vdb.SlideSubtitles.String = payload.SlideSubtitles.String
	}

	if payload.SlideURL.Valid() {
		vdb.SlideURL.Valid = true
		vdb.SlideURL.String = payload.SlideURL.String
	}

	if payload.VideoURL.Valid() {
		vdb.VideoURL.Valid = true
		vdb.VideoURL.String = payload.VideoURL.String
	}

	if payload.PhotoPermission.Valid() {
		vdb.PhotoPermission.Valid = true
		vdb.PhotoPermission.String = payload.PhotoPermission.String
	}

	if payload.VideoPermission.Valid() {
		vdb.VideoPermission.Valid = true
		vdb.VideoPermission.String = payload.VideoPermission.String
	}

	if payload.Tags.Valid() {
		vdb.Tags.Valid = true
		vdb.Tags.String = string(payload.Tags.String)
	}

	if err := vdb.Create(tx); err != nil {
		return err
	}

	if err := payload.L10N.CreateLocalizedStrings(tx, "Conference", vdb.EID); err != nil {
		return err
	}
	return nil
}
