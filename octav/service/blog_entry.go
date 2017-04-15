package service

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"io"
	"net/url"
	"sort"
	"time"

	"github.com/builderscon/octav/octav/cache"
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/internal/context"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	pdebug "github.com/lestrrat/go-pdebug"
	urlenc "github.com/lestrrat/go-urlenc"
	"github.com/pkg/errors"
)

func (v *BlogEntrySvc) Init() {}

func (v *BlogEntrySvc) populateRowForCreate(ctx context.Context, vdb *db.BlogEntry, payload *model.CreateBlogEntryRequest) error {
	vdb.EID = tools.UUID()
	vdb.ConferenceID = payload.ConferenceID
	vdb.Title = payload.Title
	vdb.Status = payload.Status

	// Parse the URL, and do away with the URL fragment, if any
	u, err := url.Parse(payload.URL)
	if err != nil {
		return errors.Wrap(err, "failed to parse URL")
	}
	u.Fragment = ""
	vdb.URL = u.String()

	h := sha1.New()
	io.WriteString(h, payload.URL)
	vdb.URLHash = fmt.Sprintf("%x", (h.Sum(nil)))
	return nil
}

func (v *BlogEntrySvc) populateRowForUpdate(ctx context.Context, vdb *db.BlogEntry, payload *model.UpdateBlogEntryRequest) error {
	if payload.Title.Valid() {
		vdb.Title = payload.Title.String
	}
	if payload.URL.Valid() {
		vdb.URL = payload.URL.String
	}
	return nil
}

func invalidateBlogEntryLoadByConferenceID(confID string) error {
	c := Cache()

	var r model.ListBlogEntriesRequest
	r.ConferenceID = confID

	var keys []string
	for _, status := range [][]string{{"private"}, {"public"}, {"private", "public"}} {
		for _, verifiedCall := range []bool{true, false} {
			for _, lang := range []string{"ja", "en", ""} {
				if lang == "" {
					r.Lang.ValidFlag = false
					r.Lang.String = ""
				} else {
					r.Lang.Set(lang)
				}

				r.Status = status
				r.VerifiedCall = verifiedCall
				keybytes, err := urlenc.Marshal(r)
				if err != nil {
					return errors.Wrap(err, "failed to marshal payload")
				}

				key := c.Key("BlogEntry", "ListFromPayload", string(keybytes))
				keys = append(keys, key)
			}
		}
	}

	for _, key := range keys {
		c.Delete(key)
	}
	return nil
}

func (v *BlogEntrySvc) PostCreateHook(ctx context.Context, _ *sql.Tx, vdb *db.BlogEntry) error {
	return invalidateBlogEntryLoadByConferenceID(vdb.ConferenceID)
}

func (v *BlogEntrySvc) PostUpdateHook(_ *sql.Tx, vdb *db.BlogEntry) error {
	return invalidateBlogEntryLoadByConferenceID(vdb.ConferenceID)
}

func (v *BlogEntrySvc) PostDeleteHook(_ *sql.Tx, vdb *db.BlogEntry) error {
	return invalidateBlogEntryLoadByConferenceID(vdb.ConferenceID)
}

func (v *BlogEntrySvc) Decorate(tx *sql.Tx, m *model.BlogEntry, verifiedCall bool, lang string) (err error) {
	// If this is not a verifiedCall, we don't want to send the conference_id, status, or the url_hash
	if !verifiedCall {
		m.ConferenceID = ""
		m.Status = ""
		m.URLHash = ""
	}
	return nil
}

func (v *BlogEntrySvc) CreateFromPayload(ctx context.Context, tx *sql.Tx, result *model.BlogEntry, payload *model.CreateBlogEntryRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("BlogEntrySvc.CreateFromPayload").BindError(&err)
		defer g.End()
	}

	su := User()
	if err := su.IsConferenceAdministrator(ctx, tx, payload.ConferenceID, context.GetUserID(ctx)); err != nil {
		return errors.Wrap(err, "creating blog entries require conference administrator privileges")
	}

	var vdb db.BlogEntry
	if err := v.Create(ctx, tx, &vdb, payload); err != nil {
		return errors.Wrap(err, "failed to insert into database")
	}

	if result != nil {
		var m model.BlogEntry
		if err := m.FromRow(&vdb); err != nil {
			return errors.Wrap(err, "failed to populate model from database")
		}

		*result = m
	}
	return nil
}

func (v *BlogEntrySvc) DeleteFromPayload(ctx context.Context, tx *sql.Tx, payload *model.DeleteBlogEntryRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("BlogEntrySvc.DeleteFromPayload").BindError(&err)
		defer g.End()
	}

	var m model.BlogEntry
	if err := v.Lookup(ctx, tx, &m, payload.ID); err != nil {
		return errors.Wrap(err, "failed to look up blog entry")
	}

	su := User()
	if err := su.IsConferenceAdministrator(ctx, tx, m.ConferenceID, context.GetUserID(ctx)); err != nil {
		return errors.Wrap(err, "deleting blog entries requires conference administrator privilege")
	}

	if err := v.Delete(tx, payload.ID); err != nil {
		return errors.Wrap(err, "failed to delete session")
	}
	return nil
}

func (v *BlogEntrySvc) ListFromPayload(ctx context.Context, tx *sql.Tx, result *model.BlogEntryList, payload *model.ListBlogEntriesRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("BlogEntrySvc.ListFromPayload").BindError(&err)
		defer g.End()
	}

	status := payload.Status
	if len(status) == 0 {
		status = append(status, model.StatusPublic)
	}

	sort.Strings(payload.Status) // normalize
	keybytes, err := urlenc.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "failed to marshal payload")
	}
	c := Cache()
	key := c.Key("BlogEntry", "ListFromPayload", string(keybytes))
	x, err := c.GetOrSet(key, result, func() (interface{}, error) {
		if pdebug.Enabled {
			pdebug.Printf("CACHE MISS: Re-generating")
		}

		var vdbl db.BlogEntryList
		if err := vdbl.LoadByConference(tx, payload.ConferenceID, status); err != nil {
			return nil, errors.Wrap(err, "failed to load from database")
		}

		l := make(model.BlogEntryList, len(vdbl))
		for i, vdb := range vdbl {
			if err := l[i].FromRow(&vdb); err != nil {
				return nil, errors.Wrap(err, "failed to populate model from database")
			}

			if err := v.Decorate(tx, &l[i], payload.VerifiedCall, payload.Lang.String); err != nil {
				return nil, errors.Wrap(err, "failed to decorate session with associated data")
			}
		}

		return &l, nil
	}, cache.WithExpires(10*time.Minute))

	if err != nil {
		return err
	}

	*result = *(x.(*model.BlogEntryList))
	return nil

}
