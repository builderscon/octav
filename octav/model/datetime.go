package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"
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
			default:
				return errors.New("invalid value to parse for conference date")
			}

			var dt Date
			if err := dt.Parse(s.(string)); err != nil {
				return err
			}
			ret[i] = dt
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
	return json.Marshal(cd.String())
}

func (cd *ConferenceDate) UnmarshalJSON(data []byte) error {
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
			default:
				return errors.New("invalid value to parse for conference date")
			}

			var dt ConferenceDate
			if err := dt.Parse(s.(string)); err != nil {
				return err
			}
			ret[i] = dt
		}
	default:
		return errors.New("invaid value to parse for conference date")
	}

	*cdl = ret
	return nil
}
