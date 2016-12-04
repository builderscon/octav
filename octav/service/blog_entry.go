package service

import (
	"crypto/sha1"
	"io"

	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
)

func (v *BlogEntrySvc) Init() {}

func (v *BlogEntrySvc) populateRowForCreate(vdb *db.BlogEntry, payload *model.CreateBlogEntryRequest) error {
	vdb.ConferenceID = payload.ConferenceID
	vdb.Title = payload.Title
	vdb.URL = payload.URL

	h := sha1.New()
	io.WriteString(h, payload.URL)
	vdb.URLHash = string(h.Sum(nil))
	return nil
}

func (v *BlogEntrySvc) populateRowForUpdate(vdb *db.BlogEntry, payload *model.UpdateBlogEntryRequest) error {
	if payload.Title.Valid() {
		vdb.Title = payload.Title.String
	}
	if payload.URL.Valid() {
		vdb.URL = payload.URL.String
	}
	return nil
}
