package service

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/cloud/storage"

	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/internal/errors"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	"github.com/lestrrat/go-pdebug"
)

func (v *Conference) populateRowForCreate(vdb *db.Conference, payload model.CreateConferenceRequest) error {
	vdb.EID = tools.UUID()
	vdb.Slug = payload.Slug
	vdb.Title = payload.Title
	vdb.SeriesID = payload.SeriesID
	vdb.Status = "private"

	if payload.SubTitle.Valid() {
		vdb.SubTitle.Valid = true
		vdb.SubTitle.String = payload.SubTitle.String
	}
	return nil
}

func (v *Conference) populateRowForUpdate(vdb *db.Conference, payload model.UpdateConferenceRequest) error {
	if payload.SeriesID.Valid() {
		vdb.SeriesID = payload.SeriesID.String
	}

	if payload.CoverURL.Valid() {
		vdb.CoverURL.Valid = true
		vdb.CoverURL.String = payload.CoverURL.String
	}

	if payload.Slug.Valid() {
		vdb.Slug = payload.Slug.String
	}

	if payload.Title.Valid() {
		vdb.Title = payload.Title.String
	}

	if payload.Status.Valid() {
		vdb.Status = payload.Status.String
	}

	if payload.SubTitle.Valid() {
		vdb.SubTitle.Valid = true
		vdb.SubTitle.String = payload.SubTitle.String
	}
	return nil
}

func (v *Conference) CreateDefaultSessionTypes(tx *db.Tx, c *model.Conference) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Conference.CreateDefaultSessionTypes").BindError(&err)
		defer g.End()
	}
	sst := SessionType{}

	var stocktypes []model.AddSessionTypeRequest

	for _, dur := range []int{3600, 1800} {
		r := model.AddSessionTypeRequest{
			Name:     fmt.Sprintf("%d min", dur/60),
			Abstract: fmt.Sprintf("%d minute session", dur/60),
			Duration: dur,
		}
		r.L10N.Set("ja", "abstract", fmt.Sprintf("%d分枠", dur/60))
		stocktypes = append(stocktypes, r)
	}
	r := model.AddSessionTypeRequest{
		Name:     "Lightning Talk",
		Abstract: "5 minute session about anything you want",
		Duration: 300,
	}
	r.L10N.Set("ja", "abstract", "5分間で強制終了、初心者も安心枠")
	stocktypes = append(stocktypes, r)

	for _, r := range stocktypes {
		var vdb db.SessionType
		r.ConferenceID = c.ID
		if err := sst.Create(tx, &vdb, model.CreateSessionTypeRequest{r}); err != nil {
			return errors.Wrap(err, "failed to create default session type")
		}
	}
	return nil
}

func (v *Conference) CreateFromPayload(tx *db.Tx, payload model.CreateConferenceRequest, result *model.Conference) error {
	su := User{}
	if err := su.IsConferenceSeriesAdministrator(tx, payload.SeriesID, payload.UserID); err != nil {
		return errors.Wrap(err, "creating a conference requires conference administrator privilege")
	}

	vdb := db.Conference{}
	if err := v.Create(tx, &vdb, payload); err != nil {
		return errors.Wrap(err, "failed to store in database")
	}

	// Description, CFPLead, CFPPresubmitInstructions, CFPPostsubmitInstruction
	// must be created
	cc := db.ConferenceComponent{
		ConferenceID: vdb.EID,
		CreatedOn: time.Now(),
	}
	if payload.Description.Valid() && payload.Description.String != "" {
		cc.EID = tools.UUID()
		cc.Name = "description"
		cc.Value = payload.Description.String
		if err := cc.Create(tx); err != nil {
			return errors.Wrap(err, "failed to insert description")
		}
	}

	if payload.CFPLeadText.Valid() && payload.CFPLeadText.String != "" {
		cc.EID = tools.UUID()
		cc.Name = "cfp_lead_text"
		cc.Value = payload.CFPLeadText.String
		if err := cc.Create(tx); err != nil {
			return errors.Wrap(err, "failed to insert description")
		}
	}

	if payload.CFPPostSubmitInstructions.Valid() && payload.CFPPostSubmitInstructions.String != "" {
		cc.EID = tools.UUID()
		cc.Name = "cfp_post_submit_instructions"
		cc.Value = payload.CFPPostSubmitInstructions.String
		if err := cc.Create(tx); err != nil {
			return errors.Wrap(err, "failed to insert description")
		}
	}

	if payload.CFPPreSubmitInstructions.Valid() && payload.CFPPreSubmitInstructions.String != "" {
		cc.EID = tools.UUID()
		cc.Name = "cfp_pre_submit_instructions"
		cc.Value = payload.CFPPreSubmitInstructions.String
		if err := cc.Create(tx); err != nil {
			return errors.Wrap(err, "failed to insert description")
		}
	}

	if err := v.AddAdministrator(tx, vdb.EID, payload.UserID); err != nil {
		return errors.Wrap(err, "failed to associate administrators to conference")
	}

	c := model.Conference{}
	if err := c.FromRow(vdb); err != nil {
		return errors.Wrap(err, "failed to populate model from database")
	}

	if err := v.CreateDefaultSessionTypes(tx, &c); err != nil {
		return errors.Wrap(err, "failed to create default session types")
	}

	*result = c
	return nil
}

var slugSplitRx = regexp.MustCompile(`^/([^/]+)/(.+)$`)

func (v *Conference) LookupBySlug(tx *db.Tx, c *model.Conference, payload model.LookupConferenceBySlugRequest) error {
	matches := slugSplitRx.FindStringSubmatch(payload.Slug)
	if matches == nil {
		return errors.New("invalid slug pattern")
	}
	seriesSlug := matches[1]
	confSlug := matches[2]

	// XXX cache this later!!!
	// This is in two steps so we can leverage existing vdb.LoadByEID()
	row := tx.QueryRow(`SELECT `+db.ConferenceTable+`.eid FROM `+db.ConferenceTable+` JOIN `+db.ConferenceSeriesTable+` ON `+db.ConferenceSeriesTable+`.eid = `+db.ConferenceTable+`.series_id WHERE `+db.ConferenceSeriesTable+`.slug = ? AND `+db.ConferenceTable+`.slug = ?`, seriesSlug, confSlug)

	var eid string
	if err := row.Scan(&eid); err != nil {
		return errors.Wrap(err, "failed to select conference id from slug")
	}

	return v.LookupFromPayload(tx, c, model.LookupConferenceRequest{ID: eid, Lang: payload.Lang})
}

func (v *Conference) AddAdministrator(tx *db.Tx, cid, uid string) error {
	c := db.ConferenceAdministrator{
		ConferenceID: cid,
		UserID:       uid,
	}
	return c.Create(tx, db.WithInsertIgnore(true))
}

func (v *Conference) AddAdministratorFromPayload(tx *db.Tx, payload model.AddConferenceAdminRequest) error {
	su := User{}
	if err := su.IsConferenceAdministrator(tx, payload.ConferenceID, payload.UserID); err != nil {
		return errors.Wrap(err, "adding a conference administrator requires conference administrator privilege")
	}

	return errors.Wrap(v.AddAdministrator(tx, payload.ConferenceID, payload.AdminID), "failed to add administrator")
}

const datefmt = `2006-01-02`

func (v *Conference) LoadByRange(tx *db.Tx, vdbl *db.ConferenceList, since, rangeStart, rangeEnd string, limit int) error {
	var rs time.Time
	var re time.Time
	var err error

	if rangeStart != "" {
		rs, err = time.Parse(datefmt, rangeStart)
		if err != nil {
			return err
		}
	}

	if rangeEnd != "" {
		re, err = time.Parse(datefmt, rangeEnd)
		if err != nil {
			return err
		}
	}

	if err := vdbl.LoadByRange(tx, since, rs, re, limit); err != nil {
		return err
	}

	return nil
}

func (v *Conference) AddDatesFromPayload(tx *db.Tx, payload model.AddConferenceDatesRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Conference.AddDatesFromPayload").BindError(&err)
		defer g.End()
	}

	su := User{}
	if err := su.IsConferenceAdministrator(tx, payload.ConferenceID, payload.UserID); err != nil {
		return errors.Wrap(err, "adding conference dates requires conference administrator privilege")
	}

	for _, date := range payload.Dates {
		if pdebug.Enabled {
			pdebug.Printf("Adding conference date %s", date)
		}
		cd := db.ConferenceDate{
			ConferenceID: payload.ConferenceID,
			Date:         date.Date.String(),
			Open:         sql.NullString{String: date.Open.String(), Valid: true},
			Close:        sql.NullString{String: date.Close.String(), Valid: true},
		}
		if err := cd.Create(tx, db.WithInsertIgnore(true)); err != nil {
			return err
		}
	}

	return nil
}

func (v *Conference) DeleteDatesFromPayload(tx *db.Tx, payload model.DeleteConferenceDatesRequest) error {
	su := User{}
	if err := su.IsConferenceAdministrator(tx, payload.ConferenceID, payload.UserID); err != nil {
		return errors.Wrap(err, "deleting conference dates requires conference administrator privilege")
	}

	vdb := db.ConferenceDate{}
	sdatelist := make([]string, len(payload.Dates))
	for i, dt := range payload.Dates {
		sdatelist[i] = dt.String()
	}
	return vdb.DeleteDates(tx, payload.ConferenceID, sdatelist...)
}

func (v *Conference) LoadDates(tx *db.Tx, cdl *model.ConferenceDateList, cid string) error {
	vdbl := db.ConferenceDateList{}
	if err := vdbl.LoadByConferenceID(tx, cid); err != nil {
		return err
	}

	res := make(model.ConferenceDateList, len(vdbl))
	for i, vdb := range vdbl {
		dt := vdb.Date
		if i := strings.IndexByte(dt, 'T'); i > -1 { // Cheat. Loading from DB contains time....!!!!
			dt = dt[:i]
		}
		if err := res[i].Date.Parse(dt); err != nil {
			return err
		}

		if vdb.Open.Valid {
			t := vdb.Open.String
			if len(t) > 5 {
				t = t[:5]
			}
			if err := res[i].Open.Parse(t); err != nil {
				return err
			}
		}

		if vdb.Close.Valid {
			t := vdb.Close.String
			if len(t) > 5 {
				t = t[:5]
			}
			if err := res[i].Close.Parse(t); err != nil {
				return err
			}
		}
	}
	*cdl = res
	return nil
}

func (v *Conference) DeleteAdministratorFromPayload(tx *db.Tx, payload model.DeleteConferenceAdminRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Conference.DeleteAdministratorFromPayload").BindError(&err)
		defer g.End()
	}

	su := User{}
	if err := su.IsConferenceAdministrator(tx, payload.ConferenceID, payload.UserID); err != nil {
		return errors.Wrap(err, "deleting a conference administrator requires conference administrator privilege")
	}

	return db.DeleteConferenceAdministrator(tx, payload.ConferenceID, payload.AdminID)
}

func (v *Conference) LoadAdmins(tx *db.Tx, cdl *model.UserList, cid string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Conference.LoadAdmins").BindError(&err)
		defer g.End()
	}

	var vdbl db.UserList
	if err := db.LoadConferenceAdministrators(tx, &vdbl, cid); err != nil {
		return err
	}

	if pdebug.Enabled {
		pdebug.Printf("Loaded %d admins", len(vdbl))
	}

	res := make(model.UserList, len(vdbl))
	for i, vdb := range vdbl {
		var u model.User
		if err := u.FromRow(vdb); err != nil {
			return err
		}
		res[i] = u
	}
	*cdl = res
	return nil
}

func (v *Conference) AddVenueFromPayload(tx *db.Tx, payload model.AddConferenceVenueRequest) error {
	su := User{}
	if err := su.IsConferenceAdministrator(tx, payload.ConferenceID, payload.UserID); err != nil {
		return errors.Wrap(err, "adding a conference venue requires conference administrator privilege")
	}
	cd := db.ConferenceVenue{
		ConferenceID: payload.ConferenceID,
		VenueID:      payload.VenueID,
	}
	if err := cd.Create(tx, db.WithInsertIgnore(true)); err != nil {
		return errors.Wrap(err, "failed to insert new conference/venue relation")
	}

	return nil
}

func (v *Conference) DeleteVenueFromPayload(tx *db.Tx, payload model.DeleteConferenceVenueRequest) error {
	su := User{}
	if err := su.IsConferenceAdministrator(tx, payload.ConferenceID, payload.UserID); err != nil {
		return errors.Wrap(err, "deleting a conference venue requires conference administrator privilege")
	}
	return errors.Wrap(db.DeleteConferenceVenue(tx, payload.ConferenceID, payload.VenueID), "failed to delete conference venue")
}

func (v *Conference) LoadVenues(tx *db.Tx, cdl *model.VenueList, cid string) error {
	var vdbl db.VenueList
	if err := db.LoadConferenceVenues(tx, &vdbl, cid); err != nil {
		return err
	}

	res := make(model.VenueList, len(vdbl))
	for i, vdb := range vdbl {
		var u model.Venue
		if err := u.FromRow(vdb); err != nil {
			return err
		}
		res[i] = u
	}
	*cdl = res
	return nil
}

func (v *Conference) LoadTextComponents(tx *db.Tx, c *model.Conference) error {
	var ccl db.ConferenceComponentList

	if err := ccl.LoadByConferenceID(tx, c.ID); err != nil {
		return errors.Wrap(err, "failed to load text components for conference")
	}

	for _, cc := range ccl {
		switch cc.Name {
		case "description":
			c.Description = cc.Value
		case "cfp_lead_text":
			c.CFPLeadText = cc.Value
		case "cfp_pre_submit_instructions":
			c.CFPPreSubmitInstructions = cc.Value
		case "cfp_post_submit_instructions":
			c.CFPPostSubmitInstructions = cc.Value
		}
	}
	return nil
}

func (v *Conference) Decorate(tx *db.Tx, c *model.Conference, trustedCall bool, lang string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Conference.Decorate").BindError(&err)
		defer g.End()
	}

	if seriesID := c.SeriesID; seriesID != "" {
		sdb := db.ConferenceSeries{}
		if err := sdb.LoadByEID(tx, seriesID); err != nil {
			return errors.Wrapf(err, "failed to load conferences series '%s'", seriesID)
		}

		s := model.ConferenceSeries{}
		if err := s.FromRow(sdb); err != nil {
			return errors.Wrapf(err, "failed to load conferences series '%s'", seriesID)
		}
		c.Series = &s
		c.FullSlug = s.Slug + "/" + c.Slug
	}

	if c.CoverURL == "" {
		// TODO: fix later
		c.CoverURL = "https://builderscon.io/assets/images/heroimage.png"
	}

	if err := v.LoadTextComponents(tx, c); err != nil {
		return errors.Wrapf(err, "failed to load conference text components for '%s'", c.ID)
	}

	if err := v.LoadDates(tx, &c.Dates, c.ID); err != nil {
		return errors.Wrapf(err, "failed to load conference date for '%s'", c.ID)
	}

	if err := v.LoadAdmins(tx, &c.Administrators, c.ID); err != nil {
		return errors.Wrapf(err, "failed to load administrators for '%s'", c.ID)
	}

	if err := v.LoadVenues(tx, &c.Venues, c.ID); err != nil {
		return errors.Wrapf(err, "failed to load venues for '%s'", c.ID)
	}

	if err := v.LoadFeaturedSpeakers(tx, &c.FeaturedSpeakers, c.ID); err != nil {
		return errors.Wrapf(err, "failed to load featured speakers for '%s'", c.ID)
	}

	if err := v.LoadSponsors(tx, &c.Sponsors, c.ID); err != nil {
		return errors.Wrapf(err, "failed to load sponsors for '%s'", c.ID)
	}

	sv := Venue{}
	for i := range c.Venues {
		if err := sv.Decorate(tx, &c.Venues[i], trustedCall, lang); err != nil {
			return errors.Wrap(err, "failed to decorate venue with associated data")
		}
	}

	sfs := FeaturedSpeaker{}
	for i := range c.FeaturedSpeakers {
		if err := sfs.Decorate(tx, &c.FeaturedSpeakers[i], trustedCall, lang); err != nil {
			return errors.Wrap(err, "failed to decorate featured speakers with associated data")
		}
	}

	sps := Sponsor{}
	for i := range c.Sponsors {
		if err := sps.Decorate(tx, &c.Sponsors[i], trustedCall, lang); err != nil {
			return errors.Wrap(err, "failed to decorate sponsors with associated data")
		}
	}

	switch lang {
	case "", "en":
	default:
		if err := v.ReplaceL10NStrings(tx, c, lang); err != nil {
			return errors.Wrap(err, "failed to replace L10N strings")
		}
	}

	return nil
}

func (v *Conference) GetStorage() StorageClient {
	if cl := v.Storage; cl != nil {
		return cl
	}
	return DefaultStorage
}

func (v *Conference) UploadImagesFromPayload(ctx context.Context, tx *db.Tx, row *db.Conference, payload *model.UpdateConferenceRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Conference.UploadImagesFromPayload").BindError(&err)
		defer g.End()
	}

	// There's nothing to do
	if payload.MultipartForm == nil || payload.MultipartForm.File == nil {
		return nil
	}

	field := "cover"
	fhs := payload.MultipartForm.File[field]
	if len(fhs) == 0 {
		return nil
	}

	var imgf multipart.File
	imgf, err = fhs[0].Open()
	if err != nil {
		return errors.Wrap(err, "failed to open cover file from multipart form")
	}

	var imgbuf bytes.Buffer
	if _, err := io.Copy(&imgbuf, imgf); err != nil {
		return errors.Wrap(err, "failed to copy cover image data to memory")
	}
	ct := http.DetectContentType(imgbuf.Bytes())

	// Only work with image/png or image/jpeg
	var suffix string
	switch ct {
	case "image/png":
		suffix = "png"
	case "image/jpeg":
		suffix = "jpeg"
	default:
		return errors.Errorf("Unsupported image type %s", ct)
	}

	// TODO: Validate the image
	// TODO: Avoid Google Storage hardcoding?
	// Upload this to a temporary location, then upon successful write to DB
	// rename it to $conference_id/$sponsor_id
	tmpname := time.Now().UTC().Format("2006-01-02") + "/" + tools.RandomString(64) + "." + suffix
	cl := v.GetStorage()
	err = cl.Upload(ctx, tmpname, &imgbuf, WithObjectAttrs(storage.ObjectAttrs{
		ContentType: ct,
		ACL: []storage.ACLRule{
			{storage.AllUsers, storage.RoleReader},
		},
	}))
	if err != nil {
		return errors.Wrap(err, "failed to upload file")
	}

	if pdebug.Enabled {
		pdebug.Printf("Writing '%s' to %s", field, tmpname)
	}

	dstname := "conferences/" + row.EID + "/cover." + suffix
	payload.CoverURL.Set(cl.URLFor(dstname))

	return finalizeFunc(func() (err error) {
		if pdebug.Enabled {
			g := pdebug.Marker("Finalizer for service.Conference.UploadImagesFromPayload").BindError(&err)
			defer g.End()
		}
		return cl.Move(ctx, tmpname, dstname)
	})
}

func (v *Conference) UpdateFromPayload(ctx context.Context, tx *db.Tx, payload model.UpdateConferenceRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Conference.UpdateFromPayload").BindError(&err)
		defer g.End()
	}

	su := User{}
	if err := su.IsConferenceAdministrator(tx, payload.ID, payload.UserID); err != nil {
		return errors.Wrap(err, "updating a conference requires conference administrator privilege")
	}

	vdb := db.Conference{}
	if err := vdb.LoadByEID(tx, payload.ID); err != nil {
		return errors.Wrap(err, "failed to load conference from database")
	}

	var uploadErr error
	if uploadErr = v.UploadImagesFromPayload(ctx, tx, &vdb, &payload); !errors.IsIgnorable(uploadErr) {
		return errors.Wrap(uploadErr, "failed to process image uploads")
	}

	if err := v.Update(tx, &vdb, payload); err != nil {
		return errors.Wrap(err, "failed to update sponsor in database")
	}

	var ccs ConferenceComponent
	deletedTextComponents := []string{}
	addedTextComponents := map[string]string{}
	if payload.Description.Valid() {
		s := payload.Description.String
		if len(s) == 0 {
			deletedTextComponents = append(deletedTextComponents, s)
		} else {
			addedTextComponents["description"] = s
		}
	}

	if payload.CFPLeadText.Valid() {
		s := payload.CFPLeadText.String
		if len(s) == 0 {
			deletedTextComponents = append(deletedTextComponents, s)
		} else {
			addedTextComponents["cfp_lead_text"] = s
		}
	}

	if payload.CFPPreSubmitInstructions.Valid() {
		s := payload.CFPPreSubmitInstructions.String
		if len(s) == 0 {
			deletedTextComponents = append(deletedTextComponents, s)
		} else {
			addedTextComponents["cfp_pre_submit_instructions"] = s
		}
	}

	if payload.CFPPostSubmitInstructions.Valid() {
		s := payload.CFPPostSubmitInstructions.String
		if len(s) == 0 {
			deletedTextComponents = append(deletedTextComponents, s)
		} else {
			addedTextComponents["cfp_post_submit_instructions"] = s
		}
	}

	if len(deletedTextComponents) > 0 {
		if err := ccs.DeleteByConferenceIDAndName(tx, payload.ID, deletedTextComponents...); err != nil {
			return errors.Wrap(err, "failed to delete components")
		}
	}

	if len(addedTextComponents) > 0 {
		if err := ccs.UpsertByConferenceIDAndName(tx, payload.ID, addedTextComponents); err != nil {
			return errors.Wrap(err, "failed to register components")
		}
	}

	if _, ok := errors.IsFinalizationRequired(uploadErr); ok {
		return uploadErr
	}
	return nil
}

func (v *Conference) ListFromPayload(tx *db.Tx, l *model.ConferenceList, payload model.ListConferenceRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Conference.ListFromPayload").BindError(&err)
		defer g.End()
	}

	var rs time.Time
	var re time.Time

	if payload.RangeStart.Valid() {
		if rs, err = time.Parse(datefmt, payload.RangeStart.String); err != nil {
			return errors.Wrap(err, "failed to parse range start")
		}
	}

	if payload.RangeEnd.Valid() {
		if re, err = time.Parse(datefmt, payload.RangeEnd.String); err != nil {
			return errors.Wrap(err, "failed to parse range end")
		}
	}

	status := "public"
	if payload.Status.Valid() {
		status = payload.Status.String
	}

	vdbl := db.ConferenceList{}
	if status == "any" {
		if err := vdbl.LoadByRange(tx, payload.Since.String, rs, re, int(payload.Limit.Int)); err != nil {
			return errors.Wrap(err, "failed to load list from database")
		}
	} else {
		if err := vdbl.LoadByStatusAndRange(tx, status, payload.Since.String, rs, re, int(payload.Limit.Int)); err != nil {
			return errors.Wrap(err, "failed to load list from database")
		}
	}

	r := make(model.ConferenceList, len(vdbl))
	for i, vdb := range vdbl {
		if err := (r[i]).FromRow(vdb); err != nil {
			return errors.Wrap(err, "failed populate model from database")
		}
		if err := v.Decorate(tx, &r[i], false, payload.Lang.String); err != nil {
			return errors.Wrap(err, "failed to decorate venue with associated data")
		}
	}

	*l = r
	return nil
}

func (v *Conference) LoadFeaturedSpeakers(tx *db.Tx, cdl *model.FeaturedSpeakerList, cid string) error {
	var vdbl db.FeaturedSpeakerList
	if err := db.LoadFeaturedSpeakers(tx, &vdbl, cid); err != nil {
		return err
	}

	res := make(model.FeaturedSpeakerList, len(vdbl))
	for i, vdb := range vdbl {
		var u model.FeaturedSpeaker
		if err := u.FromRow(vdb); err != nil {
			return err
		}
		res[i] = u
	}
	*cdl = res
	return nil
}

func (v *Conference) LoadSponsors(tx *db.Tx, cdl *model.SponsorList, cid string) error {
	var vdbl db.SponsorList
	if err := db.LoadSponsors(tx, &vdbl, cid); err != nil {
		return err
	}

	res := make(model.SponsorList, len(vdbl))
	for i, vdb := range vdbl {
		var u model.Sponsor
		if err := u.FromRow(vdb); err != nil {
			return err
		}
		res[i] = u
	}
	*cdl = res
	return nil
}

func (v *Conference) ListByOrganizerFromPayload(tx *db.Tx, l *model.ConferenceList, payload model.ListConferencesByOrganizerRequest) (err error) {
	var vdbl db.ConferenceList
	if err := db.ListConferencesByOrganizer(tx, &vdbl, payload.OrganizerID, payload.Since.String, int(payload.Limit.Int)); err != nil {
		return err
	}

	res := make(model.ConferenceList, len(vdbl))
	for i, vdb := range vdbl {
		if err := (res[i]).FromRow(vdb); err != nil {
			return errors.Wrap(err, "failed populate model from database")
		}
		if err := v.Decorate(tx, &res[i], false, payload.Lang.String); err != nil {
			return errors.Wrap(err, "failed to decorate conference with associated data")
		}
	}
	*l = res
	return nil

}

