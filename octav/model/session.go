package model

import (
	"bytes"
	"strconv"
	"strings"
	"time"
)

func (v *Session) unmarshalFromMap(m map[string]interface{}) error {
	if jv, ok := m["id"]; ok {
		switch jv.(type) {
		case string:
			v.ID = jv.(string)
		default:
			return ErrInvalidJSONFieldType{Field: "id"}
		}
	}

	if jv, ok := m["conference_id"]; ok {
		switch jv.(type) {
		case string:
			v.ConferenceID = jv.(string)
		default:
			return ErrInvalidJSONFieldType{Field: "conference_id"}
		}
	}

	if jv, ok := m["room_id"]; ok {
		switch jv.(type) {
		case string:
			v.RoomID = jv.(string)
		default:
			return ErrInvalidJSONFieldType{Field: "room_id"}
		}
	}

	if jv, ok := m["speaker_id"]; ok {
		switch jv.(type) {
		case string:
			v.SpeakerID = jv.(string)
		default:
			return ErrInvalidJSONFieldType{Field: "speaker_id"}
		}
	}

	if jv, ok := m["title"]; ok {
		switch jv.(type) {
		case string:
			v.Title = jv.(string)
		default:
			return ErrInvalidJSONFieldType{Field: "title"}
		}
	}

	if jv, ok := m["abstract"]; ok {
		switch jv.(type) {
		case string:
			v.Abstract = jv.(string)
		default:
			return ErrInvalidJSONFieldType{Field: "abstract"}
		}
	}

	if jv, ok := m["memo"]; ok {
		switch jv.(type) {
		case string:
			v.Memo = jv.(string)
		default:
			return ErrInvalidJSONFieldType{Field: "memo"}
		}
	}

	if jv, ok := m["starts_on"]; ok {
		switch jv.(type) {
		case string:
			t, err := time.Parse(time.RFC3339, jv.(string))
			if err != nil {
				return err
			}
			v.StartsOn = t
		default:
			return ErrInvalidJSONFieldType{Field: "starts_on"}
		}
	}

	if jv, ok := m["duration"]; ok {
		switch jv.(type) {
		case float64:
			v.Duration = int(jv.(float64))
		default:
			return ErrInvalidJSONFieldType{Field: "duration"}
		}
	}

	if jv, ok := m["material_level"]; ok {
		switch jv.(type) {
		case string:
			v.MaterialLevel = jv.(string)
		default:
			return ErrInvalidJSONFieldType{Field: "material_level", Value: jv}
		}
	}

	if jv, ok := m["tags"]; ok {
		buf := bytes.Buffer{}
		switch jv.(type) {
		case []interface{}:
			l := len(jv.([]interface{}))
			for i, e := range jv.([]interface{}) {
				switch e.(type) {
				case string:
					buf.WriteString(strings.TrimSpace(e.(string)))
					if i < l {
						buf.WriteString(",")
					}
				default:
					return ErrInvalidJSONFieldType{Field: "tags[" + strconv.Itoa(i) + "]", Value: e}
				}
			}
			v.Tags = TagString(buf.String())
		default:
			return ErrInvalidJSONFieldType{Field: "tags", Value: jv}
		}
	}

	if jv, ok := m["category"]; ok {
		switch jv.(type) {
		case string:
			v.Category = jv.(string)
		default:
			return ErrInvalidJSONFieldType{Field: "category"}
		}
	}

	if jv, ok := m["spoken_language"]; ok {
		switch jv.(type) {
		case string:
			v.SpokenLanguage = jv.(string)
		default:
			return ErrInvalidJSONFieldType{Field: "spoken_language"}
		}
	}

	if jv, ok := m["slide_language"]; ok {
		switch jv.(type) {
		case string:
			v.SlideLanguage = jv.(string)
		default:
			return ErrInvalidJSONFieldType{Field: "slide_language"}
		}
	}

	if jv, ok := m["slide_subtitles"]; ok {
		switch jv.(type) {
		case string:
			v.SlideSubtitles = jv.(string)
		default:
			return ErrInvalidJSONFieldType{Field: "slide_subtitles"}
		}
	}

	if jv, ok := m["slide_url"]; ok {
		switch jv.(type) {
		case string:
			v.SlideURL = jv.(string)
		default:
			return ErrInvalidJSONFieldType{Field: "slide_url"}
		}
	}

	if jv, ok := m["video_url"]; ok {
		switch jv.(type) {
		case string:
			v.VideoURL = jv.(string)
		default:
			return ErrInvalidJSONFieldType{Field: "video_url"}
		}
	}

	if jv, ok := m["photo_permission"]; ok {
		switch jv.(type) {
		case string:
			v.PhotoPermission = jv.(string)
		default:
			return ErrInvalidJSONFieldType{Field: "photo_permission"}
		}
	}

	if jv, ok := m["video_permission"]; ok {
		switch jv.(type) {
		case string:
			v.VideoURL = jv.(string)
		default:
			return ErrInvalidJSONFieldType{Field: "video_permission"}
		}
	}

	if jv, ok := m["has_interpretation"]; ok {
		switch jv.(type) {
		case bool:
			v.HasInterpretation = jv.(bool)
		default:
			return ErrInvalidJSONFieldType{Field: "has_interpretation"}
		}
	}

	if jv, ok := m["status"]; ok {
		switch jv.(type) {
		case string:
			v.Status = jv.(string)
		default:
			return ErrInvalidJSONFieldType{Field: "status"}
		}
	}

	if jv, ok := m["confirmed"]; ok {
		switch jv.(type) {
		case bool:
			v.Confirmed = jv.(bool)
		default:
			return ErrInvalidJSONFieldType{Field: "confirmed"}
		}
	}

	return nil
}
