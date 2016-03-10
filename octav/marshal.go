package octav

import (
	"bytes"
	"encoding/json"

	"github.com/lestrrat/go-urlenc"
)

func marshalJSONWithL10N(buf []byte, lf LocalizedFields) ([]byte, error) {
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

func marshalURLWithL10N(buf []byte, lf LocalizedFields) ([]byte, error) {
	if lf.Len() == 0 {
		return buf, nil
	}

	l10buf, err := urlenc.Marshal(lf)
	if err != nil {
		return nil, err
	}

	return append(append(buf, '&'), l10buf...), nil
}

