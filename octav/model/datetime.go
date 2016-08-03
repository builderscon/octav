package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
	"unicode"

	"github.com/pkg/errors"
)

func NewDate(y, m, d int) Date {
	return Date{Year: y, Month: m, Day: d}
}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Date) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	return d.Parse(s)
}

func (d *Date) Parse(s string) error {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
		return errors.New("failed to parse model.Date: " + err.Error())
	}

	d.Year = t.Year()
	d.Month = int(t.Month())
	d.Day = t.Day()
	return nil
}

func (dl *DateList) Extract(v interface{}) error {
	var ret []Date
	switch v.(type) {
	case []string:
		ret = make([]Date, len(v.([]string)))
		for i, s := range v.([]string) {
			var dt Date
			if err := dt.Parse(s); err != nil {
				return err
			}
			ret[i] = dt
		}
	case []interface{}:
		vl := v.([]interface{})
		ret = make([]Date, len(vl))
		for i, s := range vl {
			switch s.(type) {
			case string:
				if err := ret[i].Parse(s.(string)); err != nil {
					return err
				}
			case map[string]interface{}:
				m := s.(map[string]interface{})
				v, ok := m["date"]
				if !ok {
					return errors.New("missing required field 'date'")
				}
				s2, ok := v.(string)
				if !ok {
					return errors.New("invalid type for required field 'date'")
				}
				if err := ret[i].Parse(s2); err != nil {
					return err
				}
			default:
				return errors.New("invalid value to parse for conference date")
			}
		}
	default:
		return errors.New("invaid value to parse for date")
	}

	*dl = ret
	return nil
}

func NewWallClock(h, m int) WallClock {
	return WallClock{hour: h, minute: m, Valid: true}
}

func (w WallClock) String() string {
	if !w.Valid {
		return ""
	}
	return fmt.Sprintf("%02d:%02d", w.hour, w.minute)
}

func (w WallClock) MarshalJSON() ([]byte, error) {
	if !w.Valid {
		return nil, nil
	}
	return json.Marshal(w.String())
}

func (w *WallClock) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	return w.Parse(s)
}

func (w *WallClock) Parse(s string) error {
	t, err := time.Parse("15:04", s)
	if err != nil {
		return errors.New("failed to parse WallClock: " + err.Error())
	}

	w.hour = t.Hour()
	w.minute = t.Minute()
	w.Valid = true
	return nil
}

func (cd ConferenceDate) String() string {
	buf := bytes.Buffer{}
	buf.WriteString(cd.Date.String())
	if cd.Open.Valid {
		buf.WriteByte('[')
		buf.WriteString(cd.Open.String())
		if cd.Close.Valid {
			buf.WriteByte('-')
			buf.WriteString(cd.Close.String())
		}
		buf.WriteByte(']')
	}
	return buf.String()
}

func (cd ConferenceDate) MarshalJSON() ([]byte, error) {
	// conference dates are represented as:
	// { date: "YYYY-MM-DD", open: "HH:MM", close: "HH:MM" }
	m := map[string]string{
		"encoded": cd.String(),
		"date":    cd.Date.String(),
	}
	if s := cd.Open.String(); len(s) > 0 {
		m["open"] = cd.Open.String()
	}
	if s := cd.Close.String(); len(s) > 0 {
		m["close"] = cd.Close.String()
	}
	return json.Marshal(m)
}

func (cd *ConferenceDate) extractMap(m map[string]interface{}) error {
	di, ok := m["date"]
	if !ok {
		return errors.New("missing required field 'date'")
	}
	ds, ok := di.(string)
	if !ok {
		return errors.New("invalid type for field 'date'")
	}

	if err := cd.Date.Parse(ds); err != nil {
		return err
	}

	if v, ok := m["open"]; ok {
		s, ok := v.(string)
		if !ok {
			return errors.New("invalid type for field 'open'")
		}
		if err := cd.Open.Parse(s); err != nil {
			return err
		}
	}

	if v, ok := m["close"]; ok {
		s, ok := v.(string)
		if !ok {
			return errors.New("invalid type for field 'close'")
		}
		if err := cd.Close.Parse(s); err != nil {
			return err
		}
	}
	return nil
}

func (cd *ConferenceDate) UnmarshalJSON(data []byte) error {
	var isHash bool
	for i := 0; i < len(data); i++ {
		if unicode.IsSpace(rune(data[i])) {
			continue
		}
		if data[i] == '{' {
			isHash = true
		}
	}

	if isHash {
		m := map[string]interface{}{}
		if err := json.Unmarshal(data, &m); err != nil {
			return err
		}
		if err := cd.extractMap(m); err != nil {
			return err
		}
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	return cd.Parse(s)
}

func (cd *ConferenceDate) Parse(s string) error {
	// YYYY-MM-DD
	if len(s) < 10 {
		return errors.New("invalid conference date string")
	}

	if err := cd.Date.Parse(s[:10]); err != nil {
		return err
	}

	s = s[10:]
	switch remain := len(s); remain {
	case 0:
		return nil
	default:
		if remain < 7 { // "[HH:MM-" or "[HH:MM]"
			return ErrInvalidConferenceHour
		}
	}

	if s[0] != '[' {
		return ErrInvalidConferenceHour
	}
	s = s[1:]
	if err := cd.Open.Parse(s[:5]); err != nil {
		return err
	}

	switch s[5] {
	case '-':
		// continue to parse closing hour
	case ']':
		// done parsing
		return nil
	default:
		return ErrInvalidConferenceHour
	}

	s = s[6:]

	if remain := len(s); remain < 6 { // "HH:MM]"
		return ErrInvalidConferenceHour
	}

	if err := cd.Close.Parse(s[:5]); err != nil {
		return err
	}

	if s[5] != ']' {
		return ErrInvalidConferenceHour
	}

	return nil
}

type ErrInvalidConferenceDateType struct {
	Type reflect.Type
}

func (e ErrInvalidConferenceDateType) Error() string {
	buf := bytes.Buffer{}
	buf.WriteString("invalid value type to parse for conference date: ")

	var ts string
	if e.Type == nil {
		ts = "(nil)"
	} else {
		ts = e.Type.String()
	}
	buf.WriteString(ts)
	return buf.String()
}

func (cdl *ConferenceDateList) Extract(v interface{}) error {
	var ret []ConferenceDate
	switch v.(type) {
	case []string:
		ret = make([]ConferenceDate, len(v.([]string)))
		for i, s := range v.([]string) {
			var dt ConferenceDate
			if err := dt.Parse(s); err != nil {
				return err
			}
			ret[i] = dt
		}
	case []interface{}:
		vl := v.([]interface{})
		ret = make([]ConferenceDate, len(vl))
		for i, s := range vl {
			switch s.(type) {
			case string:
				if err := ret[i].Parse(s.(string)); err != nil {
					return err
				}
			case map[string]interface{}:
				if err := ret[i].extractMap(s.(map[string]interface{})); err != nil {
					return err
				}
			default:
				return ErrInvalidConferenceDateType{Type: reflect.TypeOf(s)}
			}

		}
	default:
		return ErrInvalidConferenceDateType{Type: reflect.TypeOf(v)}
	}

	*cdl = ret
	return nil
}

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format(time.RFC3339))
}

func (t *JSONTime) UnmarshalJSON(buf []byte) error {
	x, err := time.Parse(time.RFC3339, string(buf))
	if err != nil {
		return errors.Wrap(err, "failed to parse time using RFC3339")
	}
	*t = JSONTime(x)
	return nil
}

type MaybeJSONTime struct {
	JSONTime
	ValidFlag bool
}

func (t MaybeJSONTime) Valid() bool {
	return t.ValidFlag
}

func (t MaybeJSONTime) Value() interface{} {
	return t.JSONTime
}

func (t *MaybeJSONTime) Set(x interface{}) error {
	switch x.(type) {
	case time.Time:
		t.JSONTime = JSONTime(x.(time.Time))
		t.ValidFlag = true
	case JSONTime:
		t.JSONTime = x.(JSONTime)
		t.ValidFlag = true
	case string:
		v, err := time.Parse(time.RFC3339, x.(string))
		if err != nil {
			return errors.Wrap(err, "failed to parse string value for MaybeJSONTime")
		}
		t.JSONTime = JSONTime(v)
		t.ValidFlag = true
	default:
		return errors.Errorf("invalid type %s for MaybeJSONTime.Set", reflect.TypeOf(x))
	}
	return nil
}

func (t *MaybeJSONTime) Reset() {
	t.JSONTime = JSONTime{}
	t.ValidFlag = false
}

func (t *MaybeJSONTime) UnmarshalJSON(buf []byte) error {
	if err := json.Unmarshal(buf, &t.JSONTime); err != nil {
		return errors.Wrap(err, "failed to unmarshal JSONTime")
	}
	t.ValidFlag = true
	return nil
}
