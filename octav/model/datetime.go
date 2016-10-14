package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

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
	if !cd.Open.IsZero() {
		buf.WriteString(cd.Open.Format(time.RFC3339))
	}
	buf.WriteByte('-')
	if !cd.Close.IsZero() {
		buf.WriteString(cd.Close.Format(time.RFC3339))
	}
	return buf.String()
}

func (cd ConferenceDate) MarshalJSON() ([]byte, error) {
	// conference dates are represented as:
	// { open: "rfc3339", close: "rfc3339" }
	m := map[string]string{}
	if len(cd.ID) > 0 {
		m["id"] = cd.ID
	}

	if !cd.Open.IsZero() {
		m["open"] = cd.Open.Format(time.RFC3339)
	}
	if !cd.Close.IsZero() {
		m["close"] = cd.Close.Format(time.RFC3339)
	}
	return json.Marshal(m)
}

func (cd *ConferenceDate) extractMap(m map[string]interface{}) error {
	if v, ok := m["id"]; ok {
		s, ok := v.(string)
		if !ok {
			return errors.New("invalid type for field 'id'")
		}
		cd.ID = s
	}

	if v, ok := m["open"]; ok {
		s, ok := v.(string)
		if !ok {
			return errors.New("invalid type for field 'open'")
		}
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
		cd.Open = t
	}

	if v, ok := m["close"]; ok {
		s, ok := v.(string)
		if !ok {
			return errors.New("invalid type for field 'close'")
		}
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
		cd.Close = t
	}
	return nil
}

func (cd *ConferenceDate) UnmarshalJSON(data []byte) error {
	m := map[string]interface{}{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := cd.extractMap(m); err != nil {
		return err
	}
	return nil
}

func (t JSONTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format(time.RFC3339))
}

func (t *JSONTime) parse(s string) error {
	x, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return errors.Wrap(err, "failed to parse time using RFC3339")
	}
	*t = JSONTime(x)
	return nil
}
func (t *JSONTime) UnmarshalJSON(buf []byte) error {
	var s string
	if err := json.Unmarshal(buf, &s); err != nil {
		return err
	}

	return t.parse(s)
}

func (tl *JSONTimeList) Extract(v interface{}) error {
	var ret []JSONTime
	switch v.(type) {
	case []interface{}:
	default:
		return errors.New("invalid value for JSONTimeList")
	}

	vl := v.([]interface{})
	ret = make([]JSONTime, len(vl))
	for i, s := range vl {
		switch s.(type) {
		case string:
			if err := ret[i].parse(s.(string)); err != nil {
				return err
			}
		default:
			return errors.New("invalid value for JSONTime")
		}
	}

	*tl = ret
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

func (cdl *ConferenceDate) Extract(v interface{}) error {
	switch v.(type) {
	case map[string]interface{}:
		var cd ConferenceDate
		if err := cd.extractMap(v.(map[string]interface{})); err != nil {
			return err
		}
		*cdl = cd
	default:
		return errors.Errorf("Invalid conference date type: %s", reflect.TypeOf(v).Name)
	}
	return nil
}
