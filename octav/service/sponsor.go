package service

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/internal/errors"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	"github.com/lestrrat/go-pdebug"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"google.golang.org/cloud/storage"
)

func (v *Sponsor) getMediaBucketName() string {
	v.bucketOnce.Do(func() {
		if v.MediaBucketName == "" {
			v.MediaBucketName = os.Getenv("GOOGLE_STORAGE_MEDIA_BUCKET")
		}
	})
	return v.MediaBucketName
}

func (v *Sponsor) getStorageClient(ctx context.Context) *storage.Client {
	v.storageOnce.Do(func() {
		if v.Storage == nil {
			client, err := defaultStorageClient(ctx)
			if err != nil {
				panic(err.Error())
			}
			v.Storage = client
		}
	})
	return v.Storage
}

func (v *Sponsor) populateRowForCreate(vdb *db.Sponsor, payload model.CreateSponsorRequest) error {
	vdb.EID = tools.UUID()

	vdb.ConferenceID = payload.ConferenceID
	vdb.Name = payload.Name
	vdb.URL = payload.URL
	vdb.GroupName = payload.GroupName
	vdb.SortOrder = payload.SortOrder

	return nil
}

func (v *Sponsor) populateRowForUpdate(vdb *db.Sponsor, payload model.UpdateSponsorRequest) error {
	if payload.Name.Valid() {
		vdb.Name = payload.Name.String
	}

	if payload.LogoURL1.Valid() {
		vdb.LogoURL1.Valid = true
		vdb.LogoURL1.String = payload.LogoURL1.String
	}

	if payload.LogoURL2.Valid() {
		vdb.LogoURL2.Valid = true
		vdb.LogoURL2.String = payload.LogoURL2.String
	}

	if payload.LogoURL3.Valid() {
		vdb.LogoURL3.Valid = true
		vdb.LogoURL3.String = payload.LogoURL3.String
	}

	if payload.URL.Valid() {
		vdb.URL = payload.URL.String
	}

	if payload.GroupName.Valid() {
		vdb.GroupName = payload.GroupName.String
	}

	if payload.SortOrder.Valid() {
		vdb.SortOrder = int(payload.SortOrder.Int)
	}

	return nil
}

type finalizeFunc func() error

func (ff finalizeFunc) FinalizeFunc() func() error {
	return ff
}

// Ignorable always returns true, otherwise the caller will have to
// bail out immediately
func (ff finalizeFunc) Ignorable() bool {
	return true
}

func (ff finalizeFunc) Error() string {
	return "operation needs finalization"
}

func (v *Sponsor) UploadImagesFromPayload(ctx context.Context, tx *db.Tx, row *db.Sponsor, payload *model.UpdateSponsorRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Sponsor.UploadImagesFromPayload").BindError(&err)
		defer g.End()
	}

	// There's nothing to do
	if payload.MultipartForm == nil || payload.MultipartForm.File == nil {
		return nil
	}

	bucketName := v.getMediaBucketName()
	prefix := "http://storage.googleapis.com/" + bucketName

	finalizers := make([]func() error, 0, 3)
	for _, field := range []string{"logo1", "logo2", "logo3"} {
		fhs := payload.MultipartForm.File[field]
		if len(fhs) == 0 {
			continue
		}

		var imgf multipart.File
		imgf, err = fhs[0].Open()
		if err != nil {
			return errors.Wrap(err, "failed to open logo file from multipart form")
		}

		var imgbuf bytes.Buffer
		if _, err := io.Copy(&imgbuf, imgf); err != nil {
			return errors.Wrap(err, "failed to copy logo image data to memory")
		}
		imgtyp := http.DetectContentType(imgbuf.Bytes())

		// Only work with image/png or image/jpeg
		var suffix string
		switch imgtyp {
		case "image/png":
			suffix = "png"
		case "image/jpeg":
			suffix = "jpeg"
		default:
			return errors.Errorf("Unsupported image type %s", imgtyp)
		}

		// TODO: Validate the image
		// TODO: Avoid Google Storage hardcoding?
		// Upload this to a temporary location, then upon successful write to DB
		// rename it to $conference_id/$sponsor_id
		storagecl := v.getStorageClient(ctx)
		tmpname := time.Now().UTC().Format("2006-01-02") + "/" + tools.RandomString(64) + "." + suffix
		wc := storagecl.Bucket(bucketName).Object(tmpname).NewWriter(ctx)
		wc.ContentType = imgtyp
		wc.ACL = []storage.ACLRule{{storage.AllUsers, storage.RoleReader}}

		if pdebug.Enabled {
			pdebug.Printf("Writing '%s' to %s", field, tmpname)
		}

		if _, err := io.Copy(wc, &imgbuf); err != nil {
			return errors.Wrap(err, "failed to write image to temporary location")
		}
		// Note: DO NOT defer wc.Close(), as it's part of the write operation.
		// If wc.Close() does not complete w/o errors. the write failed
		if err := wc.Close(); err != nil {
			return errors.Wrap(err, "failed to write image to temporary location")
		}

		dstname := "conferences/" + row.ConferenceID + "/" + row.EID + "-" + field + "." + suffix
		switch field {
		case "logo1":
			payload.LogoURL1.Set(prefix + "/" + dstname)
		case "logo2":
			payload.LogoURL2.Set(prefix + "/" + dstname)
		case "logo3":
			payload.LogoURL3.Set(prefix + "/" + dstname)
		}

		finalizers = append(finalizers, func() (err error) {
			if pdebug.Enabled {
				g := pdebug.Marker("Finalizer for service.Sponsor.UploadImagesFromPayload").BindError(&err)
				defer g.End()
			}

			src := storagecl.Bucket(bucketName).Object(tmpname)
			dst := storagecl.Bucket(bucketName).Object(dstname)

			if pdebug.Enabled {
				pdebug.Printf("Copying %s to %s", tmpname, dstname)
			}
			if _, err = src.CopyTo(ctx, dst, nil); err != nil {
				return errors.Wrapf(err, "failed to copy from '%s' to '%s'", tmpname, dstname)
			}
			if pdebug.Enabled {
				pdebug.Printf("Deleting %s", tmpname)
			}
			if err := src.Delete(ctx); err != nil {
				return errors.Wrapf(err, "failed to delete '%s'", tmpname)
			}

			return nil
		})
	}

	if len(finalizers) == 0 {
		return nil
	}

	return finalizeFunc(func() error {
		var g errgroup.Group
		for _, f := range finalizers {
			g.Go(f)
		}
		return g.Wait()
	})
}

func (v *Sponsor) CreateFromPayload(ctx context.Context, tx *db.Tx, payload model.AddSponsorRequest, result *model.Sponsor) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Sponsor.CreateFromPayload").BindError(&err)
		defer g.End()
	}

	su := User{}
	if err := su.IsConferenceAdministrator(tx, payload.ConferenceID, payload.UserID); err != nil {
		return errors.Wrap(err, "creating a featured sponsor requires conference administrator privilege")
	}

	vdb := db.Sponsor{}
	if err := v.Create(tx, &vdb, model.CreateSponsorRequest{payload}); err != nil {
		return errors.Wrap(err, "failed to store in database")
	}

	c := model.Sponsor{}
	if err := c.FromRow(vdb); err != nil {
		return errors.Wrap(err, "failed to populate model from database")
	}

	*result = c
	return nil
}

func (v *Sponsor) UpdateFromPayload(ctx context.Context, tx *db.Tx, payload model.UpdateSponsorRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Sponsor.UpdateFromPayload").BindError(&err)
		defer g.End()
	}

	vdb := db.Sponsor{}
	if err := vdb.LoadByEID(tx, payload.ID); err != nil {
		return errors.Wrap(err, "failed to load featured sponsor from database")
	}

	su := User{}
	if err := su.IsConferenceAdministrator(tx, vdb.ConferenceID, payload.UserID); err != nil {
		return errors.Wrap(err, "updating a featured sponsor requires conference administrator privilege")
	}

	var uploadErr error
	if uploadErr := v.UploadImagesFromPayload(ctx, tx, &vdb, &payload); !errors.IsIgnorable(uploadErr) {
		return errors.Wrap(uploadErr, "failed to process image uploads")
	}

	if err := v.Update(tx, &vdb, payload); err != nil {
		return errors.Wrap(err, "failed to load featured sponsor from database")
	}

	if cb, ok := errors.IsFinalizationRequired(uploadErr); ok {
		return errors.Wrap(cb(), "failed to finalize image upload")
	}
	return nil

}

func (v *Sponsor) DeleteFromPayload(tx *db.Tx, payload model.DeleteSponsorRequest) error {
	var m db.Sponsor
	if err := m.LoadByEID(tx, payload.ID); err != nil {
		return errors.Wrap(err, "failed to load featured sponsor from database")
	}

	su := User{}
	if err := su.IsConferenceAdministrator(tx, m.ConferenceID, payload.UserID); err != nil {
		return errors.Wrap(err, "deleting venues require administrator privileges")
	}

	return errors.Wrap(v.Delete(tx, m.EID), "failed to delete from database")
}

func (v *Sponsor) ListFromPayload(tx *db.Tx, result *model.SponsorList, payload model.ListSponsorsRequest) error {
	var vdbl db.SponsorList
	if err := vdbl.LoadByConferenceSinceEID(tx, payload.ConferenceID, payload.Since.String, int(payload.Limit.Int)); err != nil {
		return errors.Wrap(err, "failed to load featured sponsor from database")
	}

	l := make(model.SponsorList, len(vdbl))
	for i, vdb := range vdbl {
		if err := (l[i]).FromRow(vdb); err != nil {
			return errors.Wrap(err, "failed to populate model from database")
		}

		if err := v.Decorate(tx, &l[i], payload.Lang.String); err != nil {
			return errors.Wrap(err, "failed to decorate venue with associated data")
		}
	}

	*result = l
	return nil
}

func (v *Sponsor) Decorate(tx *db.Tx, sponsor *model.Sponsor, lang string) error {
	if sponsor.LogoURL1 == "" {
		sponsor.LogoURL1 = "http://storage.googleapis.com/media-builderscon-1248/system/nophoto_600.png"
	}

	if sponsor.LogoURL2 == "" {
		sponsor.LogoURL2 = "http://storage.googleapis.com/media-builderscon-1248/system/nophoto_400.png"
	}

	if sponsor.LogoURL3 == "" {
		sponsor.LogoURL3 = "http://storage.googleapis.com/media-builderscon-1248/system/nophoto_200.png"
	}

	if lang == "" {
		return nil
	}

	if err := v.ReplaceL10NStrings(tx, sponsor, lang); err != nil {
		return errors.Wrap(err, "failed to replace L10N strings")
	}

	return nil
}
