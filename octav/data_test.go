package octav_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/builderscon/octav/octav"
	"github.com/stretchr/testify/assert"
)

func TestVenueJSONL10NKeys(t *testing.T) {
	const name = `東京ビッグサイト`
	const address = `〒135-0063 東京都江東区有明３丁目１０−１`

	src := fmt.Sprintf(
		`{"name#ja": "%s", "address#ja": "%s"}`,
		name,
		address,
	)

	v := octav.Venue{}
	if !assert.NoError(t, json.Unmarshal([]byte(src), &v), "Unmarshal should succeed") {
		return
	}

t.Logf("%#v", v)
{
	buf, _ := json.Marshal(&v)
t.Logf("%s", buf)
}

	var lv string
	var ok bool
	lv, ok = v.L10N.Get("ja", "name")
	if !assert.True(t, ok, "L10N.Get('name') should be present") {
		return
	}
	if !assert.Equal(t, lv, name, "L10N.Get('name') should yield expected value") {
		return
	}
	lv, ok = v.L10N.Get("ja", "address")
	if !assert.True(t, ok, "L10N.Get('address') should be present") {
		return
	}
	if !assert.Equal(t, lv, address, "L10N.Get('address') should yield expected value") {
		return
	}
}