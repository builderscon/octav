package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"

	youtube "google.golang.org/api/youtube/v3"

	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

var youtubeSvc YoutubeSvc

func Youtube() *YoutubeSvc {
	return &youtubeSvc
}

func (v *YoutubeSvc) Client(ctx context.Context, confID string) (*youtube.Service, error) {
	credentialsKey := "conferences/" + confID + "/credentials/youtube"
	var credentialsBuf bytes.Buffer
	if err := CredentialStorage.Download(ctx, credentialsKey, &credentialsBuf); err != nil {
		return nil, errors.Wrap(err, "failed to download youtube credentials")
	}

	var token oauth2.Token
	if err := json.NewDecoder(&credentialsBuf).Decode(&token); err != nil {
		return nil, errors.Wrap(err, `failed to decode youtube credentials`)
	}

	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&token))
	return youtube.New(client)
}

func (v *YoutubeSvc) UploadThumbnailFromPayload(ctx context.Context, tx *db.Tx, payload *model.SetSessionVideoCoverRequest) error {
	sv := Session()

	var session model.Session
	if err := sv.Lookup(tx, &session, payload.ID); err != nil {
		return errors.Wrap(err, `failed to lookup session`)
	}

	videoID, err := sv.VideoID(&session)
	if err != nil {
		return errors.Wrap(err, `failed to get video ID`)
	}

	if payload.MultipartForm == nil || payload.MultipartForm.File == nil {
		return errors.New(`no image uploaded`)
	}

	const field = "image"
	fhs := payload.MultipartForm.File[field]
	if len(fhs) == 0 {
		return nil
	}

	var imgf multipart.File
	imgf, err = fhs[0].Open()
	if err != nil {
		return errors.Wrap(err, "failed to open image file from multipart form")
	}

	var imgbuf bytes.Buffer
	if _, err := io.Copy(&imgbuf, imgf); err != nil {
		return errors.Wrap(err, "failed to copy image image data to memory")
	}

	cl, err := v.Client(ctx, session.ConferenceID)
	if err != nil {
		return errors.Wrap(err, `failed to create youtube client`)
	}

	if _, err := cl.Thumbnails.Set(videoID).Media(&imgbuf).Do(); err != nil {
		return errors.Wrap(err, `failed to set thumbnail`)
	}
	return nil
}
