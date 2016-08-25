package service

import (
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	pdebug "github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

func (v *User) populateRowForCreate(vdb *db.User, payload model.CreateUserRequest) error {
	vdb.EID = tools.UUID()

	vdb.Nickname = payload.Nickname
	vdb.AuthVia = payload.AuthVia
	vdb.AuthUserID = payload.AuthUserID

	if payload.AvatarURL.Valid() {
		vdb.AvatarURL.Valid = true
		vdb.AvatarURL.String = payload.AvatarURL.String
	}

	if payload.FirstName.Valid() {
		vdb.FirstName.Valid = true
		vdb.FirstName.String = payload.FirstName.String
	}

	if payload.LastName.Valid() {
		vdb.LastName.Valid = true
		vdb.LastName.String = payload.LastName.String
	}

	if payload.Email.Valid() {
		vdb.Email.Valid = true
		vdb.Email.String = payload.Email.String
	}

	if payload.TshirtSize.Valid() {
		vdb.TshirtSize.Valid = true
		vdb.TshirtSize.String = payload.TshirtSize.String
	}

	return nil
}

func (v *User) populateRowForUpdate(vdb *db.User, payload model.UpdateUserRequest) error {
	if payload.Nickname.Valid() {
		vdb.Nickname = payload.Nickname.String
	}

	if payload.AuthVia.Valid() {
		vdb.AuthVia = payload.AuthVia.String
	}

	if payload.AuthUserID.Valid() {
		vdb.AuthUserID = payload.AuthUserID.String
	}

	if payload.AvatarURL.Valid() {
		vdb.AvatarURL.Valid = true
		vdb.AvatarURL.String = payload.AvatarURL.String
	}

	if payload.FirstName.Valid() {
		vdb.FirstName.Valid = true
		vdb.FirstName.String = payload.FirstName.String
	}

	if payload.LastName.Valid() {
		vdb.LastName.Valid = true
		vdb.LastName.String = payload.LastName.String
	}

	if payload.Email.Valid() {
		vdb.Email.Valid = true
		vdb.Email.String = payload.Email.String
	}

	if payload.TshirtSize.Valid() {
		vdb.TshirtSize.Valid = true
		vdb.TshirtSize.String = payload.TshirtSize.String
	}

	return nil
}

func (v *User) IsAdministrator(tx *db.Tx, id string) error {
	// TODO: cache
	return db.IsAdministrator(tx, id)
}

func (v *User) IsSystemAdmin(tx *db.Tx, id string) error {
	// TODO: cache
	u := model.User{}
	if err := v.Lookup(tx, &u, id); err != nil {
		return errors.Wrap(err, "failed to load user from database")
	}

	if !u.IsAdmin {
		return errors.Errorf("user %s lacks system administrator privileges", id)
	}
	return nil
}

func (v *User) IsConferenceSeriesAdministrator(tx *db.Tx, seriesID, userID string) error {
	// TODO: cache
	if err := db.IsConferenceSeriesAdministrator(tx, seriesID, userID); err == nil {
		return nil
	}

	if err := v.IsSystemAdmin(tx, userID); err == nil {
		return nil
	}
	return errors.Errorf("user %s lacks conference series administrator privileges for %s", userID, seriesID)
}

func (v *User) IsConferenceAdministrator(tx *db.Tx, confID, userID string) error {
	// TODO: cache
	if err := db.IsConferenceAdministrator(tx, confID, userID); err == nil {
		return nil
	}

	c := model.Conference{}
	sc := Conference{}
	if err := sc.Lookup(tx, &c, confID); err != nil {
		return errors.Wrap(err, "failed to load conference from database")
	}

	if err := v.IsConferenceSeriesAdministrator(tx, c.SeriesID, userID); err == nil {
		return nil
	}

	return errors.Errorf("user %s lacks conference administrator privileges for %s", userID, confID)
}

func (v *User) IsOwnerUser(tx *db.Tx, targetID, userID string) error {
	if targetID == userID {
		return nil
	}

	return v.IsSystemAdmin(tx, userID)
}

func (v *User) ListFromPayload(tx *db.Tx, result *model.UserList, payload model.ListUserRequest) error {
	var vdbl db.UserList
	if err := vdbl.LoadSinceEID(tx, payload.Since.String, int(payload.Limit.Int)); err != nil {
		return errors.Wrap(err, "failed to load from database")
	}

	l := make(model.UserList, len(vdbl))
	for i, vdb := range vdbl {
		if err := l[i].FromRow(vdb); err != nil {
			return errors.Wrap(err, "failed to populate model from database")
		}

		if err := v.Decorate(tx, &l[i], payload.TrustedCall, payload.Lang.String); err != nil {
			return errors.Wrap(err, "failed to decorate user with associated data")
		}
	}

	*result = l
	return nil
}

func (v *User) CreateFromPayload(tx *db.Tx, result *model.User, payload model.CreateUserRequest) error {
	// Normally we would like to limit who can create users, but this
	// is done via OAuth login, so anybody must be able to do it.
	vdb := db.User{}
	if err := v.Create(tx, &vdb, payload); err != nil {
		return errors.Wrap(err, "failed to create new user in database")
	}

	var m model.User
	if err := m.FromRow(vdb); err != nil {
		return errors.Wrap(err, "failed to populate model from database")
	}

	*result = m
	return nil
}

func (v *User) DeleteFromPayload(tx *db.Tx, payload model.DeleteUserRequest) error {
	if err := v.IsOwnerUser(tx, payload.ID, payload.UserID); err != nil {
		return errors.Wrap(err, "deleting a user requires to be the user themselves, or a system administrator")
	}

	return errors.Wrap(v.Delete(tx, payload.ID), "failed to delete user")
}

func (v *User) UpdateFromPayload(tx *db.Tx, payload model.UpdateUserRequest) error {
	if err := v.IsOwnerUser(tx, payload.ID, payload.UserID); err != nil {
		return errors.Wrap(err, "updating a user requires to be the user themselves, or a system administrator")
	}

	vdb := db.User{}
	if err := vdb.LoadByEID(tx, payload.ID); err != nil {
		return errors.Wrap(err, "failed to load from database")
	}

	return errors.Wrap(v.Update(tx, &vdb, payload), "failed to update database")
}

func (v *User) IsSessionOwner(tx *db.Tx, sessionID, userID string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("session.User.IsSessionOwner session ID = %s, user ID = %s", sessionID, userID).BindError(&err)
		defer g.End()
	}

	if err := db.IsSessionOwner(tx, sessionID, userID); err == nil {
		return nil
	}

	var ss Session
	var m model.Session
	if err := ss.Lookup(tx, &m, sessionID); err != nil {
		return errors.Wrap(err, "failed to load session")
	}

	if err := v.IsConferenceAdministrator(tx, m.ConferenceID, userID); err == nil {
		return nil
	}

	return errors.Errorf("user %s lacks session owner privileges for %s", userID, sessionID)
}

func (v *User) Decorate(tx *db.Tx, user *model.User, trustedCall bool, lang string) error {
	if !trustedCall {
		user.Email = ""
		user.TshirtSize = ""
		user.AuthVia = ""
		user.AuthUserID = ""
	}

	if lang != "" {
		if err := v.ReplaceL10NStrings(tx, user, lang); err != nil {
			return errors.Wrap(err, "failed to replace L10N strings")
		}
	}
	return nil
}

func (v *User) LookupUserByAuthUserIDFromPayload(tx *db.Tx, result *model.User, payload model.LookupUserByAuthUserIDRequest) error {
	var vdb db.User
	if err := vdb.LoadByAuthUserID(tx, payload.AuthVia, payload.AuthUserID); err != nil {
		return errors.Wrap(err, "failed to load from database")
	}

	var r model.User
	if err := r.FromRow(vdb); err != nil {
		return errors.Wrap(err, "failed to populate mode from database")
	}

	if err := v.Decorate(tx, &r, payload.TrustedCall, payload.Lang.String); err != nil {
		return errors.Wrap(err, "failed to decorate with assocaited data")
	}

	*result = r
	return nil
}
