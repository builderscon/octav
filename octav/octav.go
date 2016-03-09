//go:generate ./gendb -t Conference -t Room -t Session -t User -t Venue -t LocalizedString -d db
//go:generate ./genmodel -t Room -t User -t Venue -d .

package octav

import (
	"errors"

	"github.com/lestrrat/go-pdebug"
)

func (r *ListRoomRequest) SetPropValue(s string, v interface{}) error {
	if pdebug.Enabled {
		pdebug.Printf("ListRoomRequest.SetPropValue(%s)", s)
	}

	switch s {
	case "venue_id":
		if jv, ok := v.(string); ok {
			r.VenueID = jv
			return nil
		}
		return ErrInvalidFieldType
	case "since":
		if jv, ok := v.(string); ok {
			r.Since = jv
			return nil
		}
		return ErrInvalidFieldType
	case "lang":
		if jv, ok := v.(string); ok {
			r.Lang = jv
			return nil
		}
		return ErrInvalidFieldType
	case "limit":
		if jv, ok := v.(int); ok {
			r.Limit = jv
			return nil
		}
		return ErrInvalidFieldType
	default:
		return errors.New("unknown column '" + s + "'")
	}
	return nil
}
