package service

import (
	"time"

	"context"

	"github.com/builderscon/octav/octav/cache"
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	pdebug "github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

func (v *SessionTypeSvc) Init() {}

func (v *SessionTypeSvc) populateRowForCreate(vdb *db.SessionType, payload *model.CreateSessionTypeRequest) error {
	vdb.EID = tools.UUID()
	vdb.Name = payload.Name
	vdb.ConferenceID = payload.ConferenceID
	vdb.Abstract = payload.Abstract
	vdb.Duration = payload.Duration

	if payload.SubmissionStart.Valid() {
		t, err := time.Parse(time.RFC3339, payload.SubmissionStart.String)
		if err != nil {
			return errors.Wrap(err, "failed to parse submission_start for session type")
		}
		vdb.SubmissionStart.Valid = true
		vdb.SubmissionStart.Time = t.UTC()
	}

	if payload.SubmissionEnd.Valid() {
		t, err := time.Parse(time.RFC3339, payload.SubmissionEnd.String)
		if err != nil {
			return errors.Wrap(err, "failed to parse submission_end for session type")
		}
		vdb.SubmissionEnd.Valid = true
		vdb.SubmissionEnd.Time = t.UTC()
	}

	return nil
}

func (v *SessionTypeSvc) populateRowForUpdate(vdb *db.SessionType, payload *model.UpdateSessionTypeRequest) error {
	if payload.Name.Valid() {
		vdb.Name = payload.Name.String
	}

	if payload.Abstract.Valid() {
		vdb.Abstract = payload.Abstract.String
	}

	if payload.Duration.Valid() {
		vdb.Duration = int(payload.Duration.Int)
	}

	if payload.IsDefault.Valid() {
		vdb.IsDefault = payload.IsDefault.Bool
	}

	if payload.SubmissionStart.Valid() {
		t, err := time.Parse(time.RFC3339, payload.SubmissionStart.String)
		if err != nil {
			return errors.Wrap(err, "failed to parse submission_start for session type")
		}
		vdb.SubmissionStart.Valid = true
		vdb.SubmissionStart.Time = t.UTC()
	}

	if payload.SubmissionEnd.Valid() {
		t, err := time.Parse(time.RFC3339, payload.SubmissionEnd.String)
		if err != nil {
			return errors.Wrap(err, "failed to parse submission_end for session type")
		}
		vdb.SubmissionEnd.Valid = true
		vdb.SubmissionEnd.Time = t.UTC()
	}

	return nil
}

func (v *SessionTypeSvc) IsAcceptingSubmissions(tx *db.Tx, id string) error {
	return db.IsAcceptingSubmissions(tx, id)
}

func (v *SessionTypeSvc) CreateFromPayload(ctx context.Context, tx *db.Tx, payload *model.AddSessionTypeRequest, result *model.SessionType) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.SessionType.CreateFromPayload").BindError(&err)
		defer g.End()
	}

	su := User()
	if err := su.IsConferenceAdministrator(tx, payload.ConferenceID, payload.UserID); err != nil {
		return errors.Wrap(err, "creating a featured sponsor requires conference administrator privilege")
	}

	vdb := db.SessionType{}
	if err := v.Create(tx, &vdb, &model.CreateSessionTypeRequest{payload}); err != nil {
		return errors.Wrap(err, "failed to store in database")
	}

	c := model.SessionType{}
	if err := c.FromRow(&vdb); err != nil {
		return errors.Wrap(err, "failed to populate model from database")
	}

	*result = c
	return nil
}

func (v *SessionTypeSvc) PreUpdateFromPayloadHook(ctx context.Context, tx *db.Tx, vdb *db.SessionType, payload *model.UpdateSessionTypeRequest) (err error) {
	su := User()
	if err := su.IsConferenceAdministrator(tx, vdb.ConferenceID, payload.UserID); err != nil {
		return errors.Wrap(err, "updating a featured sponsor requires conference administrator privilege")
	}
	return nil
}

func (v *SessionTypeSvc) PostUpdateHook(tx *db.Tx, vdb *db.SessionType) error {
	c := Cache()
	key := c.Key("SessionType", "LoadByConferenceID", vdb.ConferenceID)
	c.Delete(key)
	if pdebug.Enabled {
		pdebug.Printf("CACHE DEL: %s", key)
	}
	return nil
}

func (v *SessionTypeSvc) DeleteFromPayload(ctx context.Context, tx *db.Tx, payload *model.DeleteSessionTypeRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.SessionType.DeleteFromPayload").BindError(&err)
		defer g.End()
	}

	var m db.SessionType
	if err := m.LoadByEID(tx, payload.ID); err != nil {
		return errors.Wrap(err, "failed to load session type from database")
	}

	su := User()
	if err := su.IsConferenceAdministrator(tx, m.ConferenceID, payload.UserID); err != nil {
		return errors.Wrap(err, "deleting venues require administrator privileges")
	}

	if err := v.Delete(tx, m.EID); err != nil {
		return errors.Wrap(err, "failed to delete from database")
	}


	return nil
}

func (v *SessionTypeSvc) ListFromPayload(tx *db.Tx, result *model.SessionTypeList, payload *model.ListSessionTypesByConferenceRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.SessionType.ListFromPayload").BindError(&err)
		defer g.End()
	}

	var vdbl db.SessionTypeList
	if err := db.LoadSessionTypes(tx, &vdbl, payload.ConferenceID); err != nil {
		return errors.Wrap(err, "failed to load session type from database")
	}

	l := make(model.SessionTypeList, len(vdbl))
	for i, vdb := range vdbl {
		if err := (l[i]).FromRow(&vdb); err != nil {
			return errors.Wrap(err, "failed to populate model from database")
		}

		if err := v.Decorate(tx, &l[i], payload.TrustedCall, payload.Lang.String); err != nil {
			return errors.Wrap(err, "failed to decorate venue with associated data")
		}
	}

	*result = l
	return nil
}

func (v *SessionTypeSvc) Decorate(tx *db.Tx, st *model.SessionType, trustedCall bool, lang string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.SessionType.Decorate").BindError(&err)
		defer g.End()
	}

	now := time.Now()
	ssvalid := !st.SubmissionStart.IsZero()
	sevalid := !st.SubmissionEnd.IsZero()

	if ssvalid {
		st.SubmissionStart = st.SubmissionStart.UTC()
	}
	if sevalid {
		st.SubmissionEnd = st.SubmissionEnd.UTC()
	}

	if ssvalid && sevalid {
		if now.After(st.SubmissionStart) && now.Before(st.SubmissionEnd) {
			st.IsAcceptingSubmission = true
		}
	}

	if lang == "" {
		return nil
	}

	if err := v.ReplaceL10NStrings(tx, st, lang); err != nil {
		return errors.Wrap(err, "failed to replace L10N strings")
	}

	return nil
}

func (v *SessionTypeSvc) LoadByConferenceID(tx *db.Tx, cdl *model.SessionTypeList, cid string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("serviec.SessionType.LoadByConferenceID %s", cid).BindError(&err)
		defer g.End()
	}

	var ids []string
	c := Cache()
	key := c.Key("SessionType", "LoadByConferenceID", cid)
	if err := c.Get(key, &ids); err == nil {
		if pdebug.Enabled {
			pdebug.Printf("CACHE HIT: %s", key)
		}
		m := make(model.SessionTypeList, len(ids))
		for i, id := range ids {
			if err := v.Lookup(tx, &m[i], id); err != nil {
				return errors.Wrap(err, "failed to load from database")
			}
		}

		*cdl = m
		return nil
	}

	if pdebug.Enabled {
		pdebug.Printf("CACHE MISS: %s", key)
	}
	var vdbl db.SessionTypeList
	if err := db.LoadSessionTypes(tx, &vdbl, cid); err != nil {
		return err
	}

	ids = make([]string, len(vdbl))
	res := make(model.SessionTypeList, len(vdbl))
	for i, vdb := range vdbl {
		var u model.SessionType
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
