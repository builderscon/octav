package model

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"github.com/builderscon/octav/octav/db"
	"github.com/lestrrat/go-pdebug"
	urlenc "github.com/lestrrat/go-urlenc"
	"golang.org/x/text/language"
)

func (lf LocalizedFields) GetPropNames() ([]string, error) {
	lf.lock.RLock()
	defer lf.lock.RUnlock()

	var ret []string
	for lang, kv := range lf.fields {
		for k := range kv {
			ret = append(ret, k+"#"+lang)
		}
	}
	return ret, nil
}

func (lf LocalizedFields) GetPropValue(s string) (interface{}, error) {
	if v, ok := lf.GetFQKey(s); ok {
		return v, nil
	}

	return nil, errors.New("not found")
}

func (lf LocalizedFields) MarshalJSON() ([]byte, error) {
	lf.lock.RLock()
	defer lf.lock.RUnlock()

	buf := bytes.Buffer{}
	buf.WriteString("{")
	for lang, kv := range lf.fields {
		for k, v := range kv {
			jk, err := json.Marshal(k + "#" + lang)
			if err != nil {
				return nil, err
			}
			jv, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			buf.Write(jk)
			buf.WriteRune(':')
			buf.Write(jv)
			buf.WriteRune(',')
		}
	}

	b := buf.Bytes()
	b[len(b)-1] = '}' // replace trailing "," with a "}"
	return b, nil
}

func (lf LocalizedFields) Len() int {
	return len(lf.fields)
}

func (lf LocalizedFields) Languages() []string {
	lf.lock.Lock()
	defer lf.lock.Unlock()

	l := make([]string, 0, len(lf.fields))
	for k := range lf.fields {
		l = append(l, k)
	}
	return l
}

func (lf LocalizedFields) Foreach(cb func(string, string, string) error) error {
	lf.lock.RLock()
	defer lf.lock.RUnlock()

	for lang, kv := range lf.fields {
		for k, v := range kv {
			if err := cb(lang, k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (lf LocalizedFields) GetFQKey(key string) (string, bool) {
	lang, skey, err := splitFQKey(key)
	if err != nil {
		return "", false
	}

	return lf.Get(lang, skey)
}

func (lf LocalizedFields) Get(lang, key string) (string, bool) {
	lf.lock.RLock()
	defer lf.lock.RUnlock()

	kv, ok := lf.fields[lang]
	if !ok {
		return "", false
	}

	v, ok := kv[key]
	return v, ok
}

func (lf *LocalizedFields) Set(lang, key, value string) error {
	lf.lock.Lock()
	defer lf.lock.Unlock()

	if lf.fields == nil {
		lf.fields = map[string]map[string]string{}
	}

	kv, ok := lf.fields[lang]
	if !ok {
		kv = map[string]string{}
		lf.fields[lang] = kv
	}
	kv[key] = value

	return nil
}

func (lf *LocalizedFields) CreateLocalizedStrings(tx *sql.Tx, parentType, parentID string) error {
	if pdebug.Enabled {
		g := pdebug.Marker("LocalizedFields.CreateLocalizedStrings (%s)", parentType)
		defer g.End()
	}

	if lf.Len() <= 0 {
		if pdebug.Enabled {
			pdebug.Printf("Nothing to register, bailing out")
		}
		return nil
	}
	err := lf.Foreach(func(lang, key, val string) error {
		if len(val) == 0 {
			return nil
		}

		if pdebug.Enabled {
			pdebug.Printf("Creating l10n string for '%s' (%s)", key, lang)
		}
		ldb := db.LocalizedString{
			ParentType: parentType,
			ParentID:   parentID,
			Language:   lang,
			Localized:  val,
			Name:       key,
		}
		return ldb.Create(tx)
	})
	if err != nil {
		return err
	}
	return nil
}

func splitFQKey(k string) (string, string, error) {
	sp := strings.SplitN(k, "#", 2)
	if len(sp) != 2 {
		return "", "", errors.New("not a localized name")
	}

	t, err := language.Default.Parse(sp[1])
	if err != nil {
		return "", "", err
	}
	return t.String(), sp[0], nil
}

var l10nKeyRx = regexp.MustCompile(`^[^#]+#[^#]+$`)

func ExtractL10NFields(m map[string]interface{}, lf *LocalizedFields, keys []string) error {
	km := make(map[string]struct{})
	for _, k := range keys {
		km[k] = struct{}{}
	}

	for lk, lv := range m {
		switch lv.(type) {
		case string:
		default:
			continue
		}

		// Ignore names that do not match key#lang pattern
		if !l10nKeyRx.MatchString(lk) {
			continue
		}

		lang, name, err := splitFQKey(lk)
		if err != nil {
			return err
		}
		if _, ok := km[name]; !ok {
			continue
		}

		lf.Set(lang, name, lv.(string))
	}
	return nil
}

func MarshalJSONWithL10N(buf []byte, lf LocalizedFields) (ret []byte, err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("MarshalJSONWithL10N").BindError(&err)
		defer g.End()
	}

	if lf.Len() == 0 {
		return buf, nil
	}

	l10buf, err := json.Marshal(lf)
	if err != nil {
		return nil, err
	}
	b := bytes.NewBuffer(buf[:len(buf)-1])
	b.WriteRune(',') // Replace closing '}'
	b.Write(l10buf[1:])

	return b.Bytes(), nil
}

func MarshalURLWithL10N(buf []byte, lf LocalizedFields) ([]byte, error) {
	if lf.Len() == 0 {
		return buf, nil
	}

	l10buf, err := urlenc.Marshal(lf)
	if err != nil {
		return nil, err
	}

	return append(append(buf, '&'), l10buf...), nil
}
