package octav

import (
	"bytes"
	"encoding/json"
	"strings"

	"golang.org/x/text/language"
)

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

func (lf LocalizedFields) Keys() []string {
	lf.lock.Lock()
	defer lf.lock.Unlock()

	l := make([]string, 0, len(lf.fields))
	for k := range lf.fields {
		l = append(l, k)
	}
	return l
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

func (lf *LocalizedFields) Set(lang, key, value string) {
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
}

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

		sp := strings.SplitN(lk, "#", 2)
		if _, ok := km[sp[0]]; !ok {
			continue
		}

		t, err := language.Default.Parse(sp[1])
		if err != nil {
			return err
		}

		lf.Set(t.String(), sp[0], lv.(string))
	}
	return nil
}
