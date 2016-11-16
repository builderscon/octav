package service

// Automatically generated by genmodel utility. DO NOT EDIT!

import (
	"sync"
	"time"

	"github.com/builderscon/octav/octav/cache"

	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/internal/errors"
	"github.com/builderscon/octav/octav/model"
	"github.com/lestrrat/go-pdebug"
)

var _ = time.Time{}

var sponsorSvc SponsorSvc
var sponsorOnce sync.Once

func Sponsor() *SponsorSvc {
	sponsorOnce.Do(sponsorSvc.Init)
	return &sponsorSvc
}

func (v *SponsorSvc) LookupFromPayload(tx *db.Tx, m *model.Sponsor, payload model.LookupSponsorRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Sponsor.LookupFromPayload").BindError(&err)
		defer g.End()
	}
	if err = v.Lookup(tx, m, payload.ID); err != nil {
		return errors.Wrap(err, "failed to load model.Sponsor from database")
	}
	if err := v.Decorate(tx, m, payload.TrustedCall, payload.Lang.String); err != nil {
		return errors.Wrap(err, "failed to load associated data for model.Sponsor from database")
	}
	return nil
}

func (v *SponsorSvc) Lookup(tx *db.Tx, m *model.Sponsor, id string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Sponsor.Lookup").BindError(&err)
		defer g.End()
	}

	var r model.Sponsor
	key := `api.Sponsor.` + id
	c := Cache()
	_, err = c.GetOrSet(key, &r, func() (interface{}, error) {
		if pdebug.Enabled {
			pdebug.Printf(`CACHE MISS: %s`, key)
		}
		if err = r.Load(tx, id); err != nil {
			return nil, errors.Wrap(err, "failed to load model.Sponsor from database")
		}
		return &r, nil
	}, cache.WithExpires(time.Hour))
	*m = r
	return nil
}

// Create takes in the transaction, the incoming payload, and a reference to
// a database row. The database row is initialized/populated so that the
// caller can use it afterwards.
func (v *SponsorSvc) Create(tx *db.Tx, vdb *db.Sponsor, payload model.CreateSponsorRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Sponsor.Create").BindError(&err)
		defer g.End()
	}

	if err := v.populateRowForCreate(vdb, payload); err != nil {
		return err
	}

	if err := vdb.Create(tx); err != nil {
		return err
	}

	if err := payload.LocalizedFields.CreateLocalizedStrings(tx, "Sponsor", vdb.EID); err != nil {
		return err
	}
	return nil
}

func (v *SponsorSvc) Update(tx *db.Tx, vdb *db.Sponsor, payload model.UpdateSponsorRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Sponsor.Update (%s)", vdb.EID).BindError(&err)
		defer g.End()
	}

	if vdb.EID == `` {
		return errors.New("vdb.EID is required (did you forget to call vdb.Load(tx) before hand?)")
	}

	if err := v.populateRowForUpdate(vdb, payload); err != nil {
		return err
	}

	if err := vdb.Update(tx); err != nil {
		return err
	}
	key := `api.Sponsor.` + vdb.EID
	c := Cache()
	c.Delete(key)
	if pdebug.Enabled {
		pdebug.Printf(`CACHE DEL %s`, key)
	}

	return payload.LocalizedFields.Foreach(func(l, k, x string) error {
		if pdebug.Enabled {
			pdebug.Printf("Updating l10n string for '%s' (%s)", k, l)
		}
		ls := db.LocalizedString{
			ParentType: "Sponsor",
			ParentID:   vdb.EID,
			Language:   l,
			Name:       k,
			Localized:  x,
		}
		return ls.Upsert(tx)
	})
}

func (v *SponsorSvc) ReplaceL10NStrings(tx *db.Tx, m *model.Sponsor, lang string) error {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Sponsor.ReplaceL10NStrings lang = %s", lang)
		defer g.End()
	}
	switch lang {
	case "", "en":
		if len(m.Name) > 0 {
			return nil
		}
		for _, extralang := range []string{`ja`} {
			rows, err := tx.Query(`SELECT oid, parent_id, parent_type, name, language, localized FROM localized_strings WHERE parent_type = ? AND parent_id = ? AND language = ?`, "Sponsor", m.ID, extralang)
			if err != nil {
				if errors.IsSQLNoRows(err) {
					break
				}
				return errors.Wrap(err, `failed to excute query`)
			}

			var l db.LocalizedString
			for rows.Next() {
				if err := l.Scan(rows); err != nil {
					return err
				}
				if len(l.Localized) == 0 {
					continue
				}
				switch l.Name {
				case "name":
					if len(m.Name) == 0 {
						if pdebug.Enabled {
							pdebug.Printf("Replacing for key 'name' (fallback en -> %s", l.Language)
						}
						m.Name = l.Localized
					}
				}
			}
		}
		return nil
	case "all":
		rows, err := tx.Query(`SELECT oid, parent_id, parent_type, name, language, localized FROM localized_strings WHERE parent_type = ? AND parent_id = ?`, "Sponsor", m.ID)
		if err != nil {
			return err
		}

		var l db.LocalizedString
		for rows.Next() {
			if err := l.Scan(rows); err != nil {
				return err
			}
			if len(l.Localized) == 0 {
				continue
			}
			if pdebug.Enabled {
				pdebug.Printf("Adding key '%s#%s'", l.Name, l.Language)
			}
			m.LocalizedFields.Set(l.Language, l.Name, l.Localized)
		}
	default:
		rows, err := tx.Query(`SELECT oid, parent_id, parent_type, name, language, localized FROM localized_strings WHERE parent_type = ? AND parent_id = ? AND language = ?`, "Sponsor", m.ID, lang)
		if err != nil {
			return err
		}

		var l db.LocalizedString
		for rows.Next() {
			if err := l.Scan(rows); err != nil {
				return err
			}
			if len(l.Localized) == 0 {
				continue
			}

			switch l.Name {
			case "name":
				if pdebug.Enabled {
					pdebug.Printf("Replacing for key 'name'")
				}
				m.Name = l.Localized
			}
		}
	}
	return nil
}

func (v *SponsorSvc) Delete(tx *db.Tx, id string) error {
	if pdebug.Enabled {
		g := pdebug.Marker("Sponsor.Delete (%s)", id)
		defer g.End()
	}

	vdb := db.Sponsor{EID: id}
	if err := vdb.Delete(tx); err != nil {
		return err
	}
	key := `api.Sponsor.` + id
	c := Cache()
	c.Delete(key)
	if pdebug.Enabled {
		pdebug.Printf(`CACHE DEL %s`, key)
	}
	if err := db.DeleteLocalizedStringsForParent(tx, id, "Sponsor"); err != nil {
		return err
	}
	return nil
}

func (v *SponsorSvc) LoadList(tx *db.Tx, vdbl *db.SponsorList, since string, limit int) error {
	return vdbl.LoadSinceEID(tx, since, limit)
}
