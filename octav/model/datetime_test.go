package model_test

import (
	"encoding/json"
	"testing"

	"github.com/builderscon/octav/octav/model"
	"github.com/stretchr/testify/assert"
)

func TestDateJSON(t *testing.T) {
	const str = `"2016-03-18"`
	var dt model.Date
	if !assert.NoError(t, json.Unmarshal([]byte(str), &dt), "JSON unmarshal of model.Date should succeed") {
		return
	}

	buf, err := json.Marshal(dt)
	if !assert.NoError(t, err, "JSON marshal of model.Date should succeed") {
		return
	}

	if !assert.Equal(t, string(buf), str, "result of marshal produces the same result") {
		return
	}
}

func TestWallClockJSON(t *testing.T) {
	const str = `"14:42"`
	var dt model.WallClock
	if !assert.NoError(t, json.Unmarshal([]byte(str), &dt), "JSON unmarshal of model.WallClock should succeed") {
		return
	}

	buf, err := json.Marshal(dt)
	if !assert.NoError(t, err, "JSON marshal of model.WallClock should succeed") {
		return
	}

	if !assert.Equal(t, string(buf), str, "result of marshal produces the same result") {
		return
	}
}

func TestConferenceDateJSON(t *testing.T) {
	in := map[string]model.ConferenceDate{
		`"2016-03-18"`: model.ConferenceDate{
			Date: model.NewDate(2016, 3, 18),
		},
		`"2016-03-18[14:42]"`: model.ConferenceDate{
			Date: model.NewDate(2016, 3, 18),
			Open: model.NewWallClock(14, 42),
		},
		`"2016-03-18[14:42-15:19]"`: model.ConferenceDate{
			Date:  model.NewDate(2016, 3, 18),
			Open:  model.NewWallClock(14, 42),
			Close: model.NewWallClock(15, 19),
		},
	}

	for pat, expected := range in {
		t.Logf("Testing pattern '%s' (should PASS)", pat)
		var dt model.ConferenceDate
		if !assert.NoError(t, json.Unmarshal([]byte(pat), &dt), "JSON unmarshal of model.ConferenceDate should succeed") {
			return
		}

		if !assert.Equal(t, dt, expected, "Unmarshaled result should match expected date") {
			return
		}

		buf, err := json.Marshal(dt)
		if !assert.NoError(t, err, "JSON marshal of model.ConferenceDate should succeed") {
			return
		}

		var dt2 model.ConferenceDate
		if !assert.NoError(t, json.Unmarshal(buf, &dt2), "Unmarshaling newly marshaled data should succeed") {
			t.Logf("Failed to unmarshal '%s'", buf)
			return
		}

		if !assert.Equal(t, dt2, dt, "Roundtrip should create same object") {
			return
		}
	}

	failures := []string{
		`"2006-1-2"`,
		`"2006-13-42[15:45"`,
	}
	for _, pat := range failures {
		t.Logf("Testing pattern '%s' (should FAIL)", pat)
		var dt model.ConferenceDate
		if !assert.Error(t, json.Unmarshal([]byte(pat), &dt), "JSON unmarshal should fail") {
			return
		}
	}
}
