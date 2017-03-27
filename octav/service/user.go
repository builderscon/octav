package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/gettext"
	"github.com/builderscon/octav/octav/internal/errors"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	"github.com/dghubble/go-twitter/twitter"
	pdebug "github.com/lestrrat/go-pdebug"
)

func (v *UserSvc) Init() {
	v.EnableVerify = true
}

func (v *UserSvc) populateRowForCreate(vdb *db.User, payload *model.CreateUserRequest) error {
	vdb.EID = tools.UUID()

	vdb.Nickname = payload.Nickname
	vdb.AuthVia = payload.AuthVia
	vdb.AuthUserID = payload.AuthUserID
	vdb.Lang = "en"

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

func (v *UserSvc) populateRowForUpdate(vdb *db.User, payload *model.UpdateUserRequest) error {
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

	if payload.Lang.Valid() {
		vdb.Lang = payload.Lang.String
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

// IsClaimedUser verifies if the claimed identity is indeed the
// user that we have received from the payload
//
// In order for this to work, the access token must be sent to
// us once. there after, we shall use sessions to keep state.
func (v *UserSvc) IsClaimedUser(ctx context.Context, tx *sql.Tx, token, uid string) error {
	// Load the user by ID, and check where we authed it from
	var u model.User
	if err := v.Lookup(ctx, tx, &u, uid); err != nil {
		return errors.Wrap(err, `failed to lookup user`)
	}

	switch u.AuthVia {
	case "github":
		r, err := http.NewRequest("GET", "https://api.github.com/users", nil)
		if err != nil {
			return errors.Wrap(err, `failed to generate HTTP request`)
		}

		r.Header.Set("Authorization", "Bearer "+token)
		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			return errors.Wrap(err, `failed to make request to github`)
		}
		var data struct {
			ID int `json:"id"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return errors.Wrap(err, `failed to decode JSON`)
		}
		if strconv.Itoa(data.ID) != u.AuthUserID {
			return errors.New(`authentication information mismatch`)
		}
	default:
		return errors.New(`unimplemented`)
	}

	return nil
}

func (v *UserSvc) IsAdministrator(ctx context.Context, tx *sql.Tx, id string) error {
	// TODO: cache
	return db.IsAdministrator(tx, id)
}

func (v *UserSvc) IsSystemAdmin(ctx context.Context, tx *sql.Tx, id string) error {
	// TODO: cache
	var u model.User
	if err := v.Lookup(ctx, tx, &u, id); err != nil {
		return errors.Wrap(err, "failed to load user from database")
	}

	if !u.IsAdmin {
		return errors.Errorf("user %s lacks system administrator privileges", id)
	}
	return nil
}

func (v *UserSvc) IsConferenceSeriesAdministrator(ctx context.Context, tx *sql.Tx, seriesID, userID string) error {
	// TODO: cache
	if err := db.IsConferenceSeriesAdministrator(tx, seriesID, userID); err == nil {
		return nil
	}

	if err := v.IsSystemAdmin(ctx, tx, userID); err == nil {
		return nil
	}
	return errors.Errorf("user %s lacks conference series administrator privileges for %s", userID, seriesID)
}

func (v *UserSvc) IsConferenceAdministrator(ctx context.Context, tx *sql.Tx, confID, userID string) error {
	// TODO: cache
	if err := db.IsConferenceAdministrator(tx, confID, userID); err == nil {
		return nil
	}

	var c model.Conference
	sc := Conference()
	if err := sc.Lookup(ctx, tx, &c, confID); err != nil {
		return errors.Wrap(err, "failed to load conference from database")
	}

	if err := v.IsConferenceSeriesAdministrator(ctx, tx, c.SeriesID, userID); err == nil {
		return nil
	}

	return errors.Errorf("user %s lacks conference administrator privileges for %s", userID, confID)
}

func (v *UserSvc) IsOwnerUser(ctx context.Context, tx *sql.Tx, targetID, userID string) error {
	if targetID == userID {
		return nil
	}

	return v.IsSystemAdmin(ctx, tx, userID)
}

func (v *UserSvc) ListFromPayload(ctx context.Context, tx *sql.Tx, result *model.UserList, payload *model.ListUserRequest) error {
	var vdbl db.UserList
	if err := vdbl.LoadFromQuery(tx, payload.Pattern.String, payload.Since.String, int(payload.Limit.Int)); err != nil {
		return errors.Wrap(err, "failed to load from database")
	}

	l := make(model.UserList, len(vdbl))
	for i, vdb := range vdbl {
		if err := l[i].FromRow(&vdb); err != nil {
			return errors.Wrap(err, "failed to populate model from database")
		}

		if err := v.Decorate(ctx, tx, &l[i], payload.TrustedCall, payload.Lang.String); err != nil {
			return errors.Wrap(err, "failed to decorate user with associated data")
		}
	}

	*result = l
	return nil
}

func (v *UserSvc) CreateFromPayload(ctx context.Context, tx *sql.Tx, result *model.User, payload *model.CreateUserRequest) error {
	// Normally we would like to limit who can create users, but this
	// is done via OAuth login, so anybody must be able to do it.
	var vdb db.User
	if err := v.Create(ctx, tx, &vdb, payload); err != nil {
		return errors.Wrap(err, "failed to create new user in database")
	}

	var m model.User
	if err := m.FromRow(&vdb); err != nil {
		return errors.Wrap(err, "failed to populate model from database")
	}

	*result = m
	return nil
}

func (v *UserSvc) DeleteFromPayload(ctx context.Context, tx *sql.Tx, payload *model.DeleteUserRequest) error {
	if err := v.IsOwnerUser(ctx, tx, payload.ID, payload.UserID); err != nil {
		return errors.Wrap(err, "deleting a user requires to be the user themselves, or a system administrator")
	}

	return errors.Wrap(v.Delete(tx, payload.ID), "failed to delete user")
}

func (v *UserSvc) PreUpdateFromPayloadHook(ctx context.Context, tx *sql.Tx, _ *db.User, payload *model.UpdateUserRequest) error {
	return errors.Wrap(
		v.IsOwnerUser(ctx, tx, payload.ID, payload.UserID),
		"updating a user requires to be the user themselves, or a system administrator",
	)
}

func (v *UserSvc) IsSessionOwner(ctx context.Context, tx *sql.Tx, sessionID, userID string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("session.User.IsSessionOwner session ID = %s, user ID = %s", sessionID, userID).BindError(&err)
		defer g.End()
	}

	if err := db.IsSessionOwner(tx, sessionID, userID); err == nil {
		return nil
	}

	ss := Session()
	var m model.Session
	if err := ss.Lookup(ctx, tx, &m, sessionID); err != nil {
		return errors.Wrap(err, "failed to load session")
	}

	if err := v.IsConferenceAdministrator(ctx, tx, m.ConferenceID, userID); err == nil {
		return nil
	}

	return errors.Errorf("user %s lacks session owner privileges for %s", userID, sessionID)
}

func (v *UserSvc) Decorate(ctx context.Context, tx *sql.Tx, user *model.User, trustedCall bool, lang string) error {
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

func (v *UserSvc) LookupUserByAuthUserIDFromPayload(ctx context.Context, tx *sql.Tx, result *model.User, payload *model.LookupUserByAuthUserIDRequest) error {
	var vdb db.User
	if err := vdb.LoadByAuthUserID(tx, payload.AuthVia, payload.AuthUserID); err != nil {
		return errors.Wrap(err, "failed to load from database")
	}

	var r model.User
	if err := r.FromRow(&vdb); err != nil {
		return errors.Wrap(err, "failed to populate mode from database")
	}

	if err := v.Decorate(ctx, tx, &r, payload.TrustedCall, payload.Lang.String); err != nil {
		return errors.Wrap(err, "failed to decorate with assocaited data")
	}

	if err := v.PostLookupFromPayloadHook(ctx, tx, &r); err != nil {
		return errors.Wrap(err, "failed to execute PostLookupFromPayloadHook")
	}

	*result = r

	return nil
}

func (v *UserSvc) CreateTemporaryEmailFromPayload(tx *sql.Tx, key *string, payload *model.CreateTemporaryEmailRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.User.CreateTemporaryEmailFromPayload").BindError(&err)
		defer g.End()
	}

	var row db.TemporaryEmail

	row.UserID = payload.TargetID
	row.Email = payload.Email
	row.ConfirmationKey = tools.UUID()
	row.ExpiresOn = time.Now().Add(time.Duration(24 * time.Hour))

	if err := row.Upsert(tx); err != nil {
		return errors.Wrap(err, "failed to create temporary email")
	}

	*key = row.ConfirmationKey
	gettext.SetLocale(payload.Lang.String)

	t, err := Template().Get("templates/en/eml/confirm_registration.eml")
	if err != nil {
		return errors.Wrap(err, "failed to fetch email template")
	}

	var txt bytes.Buffer
	if err := t.Execute(&txt, row); err != nil {
		return errors.Wrap(err, "failed to execute template")
	}

	mg := Mailgun()
	if pdebug.Enabled {
		pdebug.Printf("Got mailgun %v", mg)
	}
	mm := MailMessage{
		Recipients: []string{payload.Email},
		Subject:    gettext.Get("Confirm Your Email Registration"),
		Text:       txt.String(),
	}

	if pdebug.Enabled {
		pdebug.Printf("Sending via mailgun: %#v", mm)
	}

	if err := mg.Send(&mm); err != nil {
		return errors.Wrap(err, "failed to send message")
	}
	return nil
}

func (v *UserSvc) ConfirmTemporaryEmailFromPayload(tx *sql.Tx, payload *model.ConfirmTemporaryEmailRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.User.ConfirmTemporaryEmailFromPayload").BindError(&err)
		defer g.End()
	}

	var row db.TemporaryEmail
	if err := row.LoadByUserIDAndConfirmationKey(tx, payload.TargetID, payload.ConfirmationKey); err != nil {
		return errors.Wrap(err, "failed to load temporary email")
	}

	var u db.User
	if err := u.LoadByEID(tx, payload.TargetID); err != nil {
		return errors.Wrap(err, "failed to load user")
	}

	u.Email.String = row.Email
	u.Email.Valid = true
	if err := u.Update(tx); err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	if err := row.Delete(tx); err != nil {
		return errors.Wrap(err, "failed to delete temporary email")
	}

	return nil
}

func (v *UserSvc) ShouldVerify(_ *model.User) bool {
	if v.EnableVerify {
		return tools.RandFloat64() < 0.1
	}
	return false
}

func (v *UserSvc) Verify(ctx context.Context, m *model.User) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.User.Verify").BindError(&err)
		defer g.End()
	}

	// Check if the avatar URL is valid
	if len(m.AvatarURL) > 0 {
		res, err := http.Head(m.AvatarURL)
		if err != nil {
			return errors.Wrap(err, "failed to make HEAD request")
		}

		if res.StatusCode == http.StatusOK {
			if pdebug.Enabled {
				pdebug.Printf("AvatarURL verified")
			}

			return nil
		}
	}

	if pdebug.Enabled {
		pdebug.Printf("AvatarURL %s is invalid", m.AvatarURL)
	}

	// Dangit, got to update it
	var newAvatarURL string
	switch m.AuthVia {
	case "twitter":
		c := Twitter()
		id, err := strconv.ParseInt(m.AuthUserID, 10, 64)
		if err != nil {
			return errors.Wrap(err, "failed to convert user id to int64")
		}
		u, _, err := c.client.Users.Show(&twitter.UserShowParams{UserID: id})
		if err != nil {
			return errors.Wrap(err, "failed to fetch twitter user information via users/show")
		}
		newAvatarURL = u.ProfileImageURLHttps
	case "github":
		newAvatarURL = "https://avatars.githubusercontent.com/u/" + m.AuthUserID
	}

	if len(newAvatarURL) == 0 {
		return errors.New("failed to fetch a new avatar url")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start db transaction")
	}

	payload := model.UpdateUserRequest{
		ID: m.ID,
	}
	payload.AvatarURL.Set(newAvatarURL)

	var vdb db.User
	if err := vdb.LoadByEID(tx, payload.ID); err != nil {
		return errors.Wrap(err, "failed to load from database")
	}

	if err := v.Update(tx, &vdb); err != nil {
		return errors.Wrap(err, "failed to update database")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit data to database")
	}
	return nil
}

func (v *UserSvc) PostLookupFromPayloadHook(ctx context.Context, tx *sql.Tx, m *model.User) error {
	if v.ShouldVerify(m) {
		go v.Verify(ctx, m)
	}
	return nil
}
