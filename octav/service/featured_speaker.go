package service

import (
	"database/sql"
	"time"

	"github.com/builderscon/octav/octav/cache"
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/internal/context"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	pdebug "github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

func (v *FeaturedSpeakerSvc) Init() {}

func (v *FeaturedSpeakerSvc) populateRowForCreate(ctx context.Context, vdb *db.FeaturedSpeaker, payload *model.CreateFeaturedSpeakerRequest) error {
	vdb.EID = tools.UUID()
	vdb.ConferenceID = payload.ConferenceID
	vdb.DisplayName = payload.DisplayName
	vdb.Description = payload.Description

	if payload.AvatarURL.Valid() {
		vdb.AvatarURL.Valid = true
		vdb.AvatarURL.String = payload.AvatarURL.String
	}

	if payload.SpeakerID.Valid() {
		vdb.SpeakerID.Valid = true
		vdb.SpeakerID.String = payload.SpeakerID.String
	}

	return nil
}

func (v *FeaturedSpeakerSvc) populateRowForUpdate(ctx context.Context, vdb *db.FeaturedSpeaker, payload *model.UpdateFeaturedSpeakerRequest) error {
	if payload.DisplayName.Valid() {
		vdb.DisplayName = payload.DisplayName.String
	}

	if payload.Description.Valid() {
		vdb.Description = payload.Description.String
	}

	if payload.SpeakerID.Valid() {
		vdb.SpeakerID.Valid = true
		vdb.SpeakerID.String = payload.SpeakerID.String
	}

	if payload.AvatarURL.Valid() {
		vdb.AvatarURL.Valid = true
		vdb.AvatarURL.String = payload.AvatarURL.String
	}

	return nil
}

func (v *FeaturedSpeakerSvc) CreateFromPayload(ctx context.Context, tx *sql.Tx, payload *model.AddFeaturedSpeakerRequest, result *model.FeaturedSpeaker) error {
	su := User()
	if err := su.IsConferenceAdministrator(ctx, tx, payload.ConferenceID, context.GetUserID(ctx)); err != nil {
		return errors.Wrap(err, "creating a featured speaker requires conference administrator privilege")
	}

	var vdb db.FeaturedSpeaker
	if err := v.Create(ctx, tx, &vdb, &model.CreateFeaturedSpeakerRequest{payload}); err != nil {
		return errors.Wrap(err, "failed to store in database")
	}

	var m model.FeaturedSpeaker
	if err := m.FromRow(&vdb); err != nil {
		return errors.Wrap(err, "failed to populate model from database")
	}

	c := Cache()
	key := c.Key("FeaturedSpeaker", "LoadByConferenceID", payload.ConferenceID)
	c.Delete(key)
	if pdebug.Enabled {
		pdebug.Printf("CACHE DEL %s", key)
	}

	*result = m
	return nil
}

func (v *FeaturedSpeakerSvc) PreUpdateFromPayloadHook(ctx context.Context, tx *sql.Tx, vdb *db.FeaturedSpeaker, payload *model.UpdateFeaturedSpeakerRequest) (err error) {
	su := User()
	if err := su.IsConferenceAdministrator(ctx, tx, vdb.ConferenceID, context.GetUserID(ctx)); err != nil {
		return errors.Wrap(err, "updating a featured speaker requires conference administrator privilege")
	}
	return nil
}

func (v *FeaturedSpeakerSvc) DeleteFromPayload(ctx context.Context, tx *sql.Tx, payload *model.DeleteFeaturedSpeakerRequest) error {
	var m db.FeaturedSpeaker
	if err := m.LoadByEID(tx, payload.ID); err != nil {
		return errors.Wrap(err, "failed to load featured speaker from database")
	}

	su := User()
	if err := su.IsConferenceAdministrator(ctx, tx, m.ConferenceID, context.GetUserID(ctx)); err != nil {
		return errors.Wrap(err, "deleting venues require administrator privileges")
	}

	if err := v.Delete(tx, m.EID); err != nil {
		return errors.Wrap(err, "failed to delete from database")
	}

	c := Cache()
	key := c.Key("FeaturedSpeaker", "LoadByConferenceID", m.ConferenceID)
	c.Delete(key)
	if pdebug.Enabled {
		pdebug.Printf("CACHE DEL %s", key)
	}
	return nil
}

func (v *FeaturedSpeakerSvc) ListFromPayload(ctx context.Context, tx *sql.Tx, result *model.FeaturedSpeakerList, payload *model.ListFeaturedSpeakersRequest) error {
	var vdbl db.FeaturedSpeakerList
	if err := vdbl.LoadByConferenceSinceEID(tx, payload.ConferenceID, payload.Since.String, int(payload.Limit.Int)); err != nil {
		return errors.Wrap(err, "failed to load featured speakers from database")
	}

	l := make(model.FeaturedSpeakerList, len(vdbl))
	for i, vdb := range vdbl {
		if err := (l[i]).FromRow(&vdb); err != nil {
			return errors.Wrap(err, "failed to populate model from database")
		}

		if err := v.Decorate(ctx, tx, &l[i], payload.TrustedCall, payload.Lang.String); err != nil {
			return errors.Wrap(err, "failed to decorate venue with associated data")
		}
	}

	*result = l
	return nil
}

func (v *FeaturedSpeakerSvc) Decorate(ctx context.Context, tx *sql.Tx, speaker *model.FeaturedSpeaker, trustedCall bool, lang string) error {
	if lang == "" {
		return nil
	}

	if err := v.ReplaceL10NStrings(tx, speaker, lang); err != nil {
		return errors.Wrap(err, "failed to replace L10N strings")
	}

	return nil
}

func (v *FeaturedSpeakerSvc) LoadByConferenceID(ctx context.Context, tx *sql.Tx, cdl *model.FeaturedSpeakerList, cid string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("serviec.FeaturedSpeaker.LoadByConferenceID %s", cid).BindError(&err)
		defer g.End()
	}

	var ids []string
	c := Cache()
	key := c.Key("FeaturedSpeaker", "LoadByConferenceID", cid)
	if err := c.Get(key, &ids); err == nil {
		if pdebug.Enabled {
			pdebug.Printf("CACHE HIT: %s", key)
		}
		m := make(model.FeaturedSpeakerList, len(ids))
		for i, id := range ids {
			if err := v.Lookup(ctx, tx, &m[i], id); err != nil {
				return errors.Wrap(err, "failed to load from database")
			}
		}

		*cdl = m
		return nil
	}

	if pdebug.Enabled {
		pdebug.Printf("CACHE MISS: %s", key)
	}
	var vdbl db.FeaturedSpeakerList
	if err := db.LoadFeaturedSpeakers(tx, &vdbl, cid); err != nil {
		return err
	}

	ids = make([]string, len(vdbl))
	res := make(model.FeaturedSpeakerList, len(vdbl))
	for i, vdb := range vdbl {
		var u model.FeaturedSpeaker
		if err := u.FromRow(&vdb); err != nil {
			return err
		}
		ids[i] = vdb.EID
		res[i] = u
	}
	*cdl = res

	c.Set(key, ids, cache.WithExpires(15*time.Minute))
	return nil
}
