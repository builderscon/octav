package model

import (
	"bytes"
	"encoding/json"
	"strings"
)

func (t TagString) MarshalJSON() ([]byte, error) {
	if t == "" {
		return []byte{'[',']'}, nil
	}
	return json.Marshal(strings.Split(string(t), ","))
}

func (t *TagString) UnmarshalJSON(data []byte) error {
	var ts []string
	if err := json.Unmarshal(data, &ts); err != nil {
		return err
	}
	buf := bytes.Buffer{}
	for _, tag := range ts {
		buf.WriteString(tag)
		buf.WriteString(",")
	}
	*t = TagString(strings.TrimSuffix(buf.String(), ","))
	return nil
}