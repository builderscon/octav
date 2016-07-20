package service

import (
	"database/sql"
	"regexp"
	"strings"
	"time"

	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	"github.com/lestrrat/go-pdebug"
	"github.com/pkg/errors"
)

func (v *Conference) populateRowForCreate(vdb *db.Conference, payload model.CreateConferenceRequest) error {
	vdb.EID = tools.UUID()
	vdb.Slug = payload.Slug
	vdb.Title = payload.Title
	vdb.SeriesID = payload.SeriesID

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

	if payload.Slug.Valid() {
		vdb.Slug = payload.Slug.String
	}

	if payload.Title.Valid() {
		vdb.Title = payload.Title.String
	}

	if payload.SubTitle.Valid() {
		vdb.SubTitle.Valid = true
		vdb.SubTitle.String = payload.SubTitle.String
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

	if err := v.AddAdministrator(tx, vdb.EID, payload.UserID); err != nil {
		return errors.Wrap(err, "failed to associate administrators to conference")
	}

	c := model.Conference{}
	if err := c.FromRow(vdb); err != nil {
		return errors.Wrap(err, "failed to populate model from database")
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

func (v *Conference) Decorate(tx *db.Tx, c *model.Conference, lang string) error {
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

	if lang != "" {
		sv := Venue{}
		for i := range c.Venues {
			if err := sv.Decorate(tx, &c.Venues[i], lang); err != nil {
				return errors.Wrap(err, "failed to decorate venue with associated data")
			}
		}
		if err := v.ReplaceL10NStrings(tx, c, lang); err != nil {
			return errors.Wrap(err, "failed to replace L10N strings")
		}
	}

	return nil
}

func (v *Conference) UpdateFromPayload(tx *db.Tx, payload model.UpdateConferenceRequest) (err error) {
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

	return errors.Wrap(v.Update(tx, &vdb, payload), "failed to load conference from database")
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

	status := "public";
	if payload.Status.Valid() {
		status = payload.Status.String
	}

	vdbl := db.ConferenceList{}
	if err := vdbl.LoadByStatusAndRange(tx, status, payload.Since.String, rs, re, int(payload.Limit.Int)); err != nil {
		return errors.Wrap(err, "failed to load list from database")
	}

	r := make(model.ConferenceList, len(vdbl))
	for i, vdb := range vdbl {
		if err := (r[i]).FromRow(vdb); err != nil {
			return errors.Wrap(err, "failed populate model from database")
		}
		if err := v.Decorate(tx, &r[i], payload.Lang.String); err != nil {
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


