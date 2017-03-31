package service

import (
	"database/sql"
	"time"

	"github.com/builderscon/octav/octav/cache"
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/internal/errors"
	"github.com/builderscon/octav/octav/model"
	pdebug "github.com/lestrrat/go-pdebug"
)

func (v *LocalizedStringSvc) Init() {}

func (v *LocalizedStringSvc) LookupFields(tx *sql.Tx, parentType, parentID, lang string, list *[]db.LocalizedString) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.LocalizedString.LookupFields (%s, %s, %s)", parentType, parentID, lang).BindError(&err)
		defer g.End()
	}

	if parentType == "" {
		return errors.New("missing parent type")
	}
	if parentID == "" {
		return errors.New("missing parent ID")
	}
	if lang == "" {
		return errors.New("missing language")
	}

	c := Cache()
	key := c.Key(`LocalizedString`, parentType, parentID, lang)
	if err := c.Get(key, list); err == nil {
		if pdebug.Enabled {
			pdebug.Printf("CACHE HIT: %s", key)
		}
		return nil
	}

	if pdebug.Enabled {
		pdebug.Printf("CACHE MISS: %s", key)
	}
	rows, err := tx.Query(`SELECT name, localized FROM localized_strings WHERE parent_type = ? AND parent_id = ? AND language = ?`, parentType, parentID, lang)
	if err != nil {
		if errors.IsSQLNoRows(err) {
			return nil
		}
		return errors.Wrap(err, `failed to excute query`)
	}

	var l db.LocalizedString
	for rows.Next() {
		if err := rows.Scan(&l.Name, &l.Localized); err != nil {
			return errors.Wrap(err, "failed to scan row from localized strings")
		}
		if len(l.Localized) == 0 {
			continue
		}
		l.ParentType = parentType
		l.ParentID = parentID
		l.Language = lang
		*list = append(*list, l)
	}

	c.Set(key, *list, cache.WithExpires(time.Hour))
	return nil
}
func (v *LocalizedStringSvc) UpdateFields(tx *sql.Tx, parentType, parentID string, fields model.LocalizedFields) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.LocalizedString.UpdateFields (%s, %s)", parentType, parentID).BindError(&err)
		defer g.End()
	}
	// We cache the entire list of localizations, so here we
	// invalidate per parent type/id
	var updated bool
	langs := make(map[string]struct{})
	defer func() {
		if !updated {
			return
		}

		c := Cache()
		for l := range langs {
			key := c.Key("LocalizedString", parentType, parentID, l)
			c.Delete(key)
			if pdebug.Enabled {
				pdebug.Printf("CACHE DEL: %s", key)
			}
		}
	}()

	return fields.Foreach(func(l, k, x string) error {
		if pdebug.Enabled {
			pdebug.Printf("Updating l10n string for '%s' (%s)", k, l)
		}
		updated = true
		langs[l] = struct{}{}
		ls := db.LocalizedString{
			ParentType: parentType,
			ParentID:   parentID,
			Language:   l,
			Name:       k,
			Localized:  x,
		}
		return ls.Upsert(tx)
	})
}
