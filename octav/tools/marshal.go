package tools

import (
	"bytes"
	"encoding/json"

	"github.com/lestrrat/go-pdebug"
	"github.com/lestrrat/go-urlenc"
)

func MarshalJSONWithL10N(buf []byte, lf LocalizedFields) (ret []byte, err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("tools.MarshalJSONWithL10N").BindError(&err)
		defer g.End()
	}

	if lf.Len() == 0 {
		return buf, nil
	}

	l10buf, err := json.Marshal(lf)
pdebug.Printf("l10buf = %s", l10buf)
pdebug.Printf("err = %s", err)
	if err != nil {
		return nil, err
	}
	b := bytes.NewBuffer(buf[:len(buf)-1])
	b.WriteRune(',') // Replace closing '}'
	b.Write(l10buf[1:])

	return b.Bytes(), nil
}

func MarshalURLWithL10N(buf []byte, lf LocalizedFields) ([]byte, error) {
	if lf.Len() == 0 {
		return buf, nil
	}

	l10buf, err := urlenc.Marshal(lf)
	if err != nil {
		return nil, err
	}

	return append(append(buf, '&'), l10buf...), nil
}
