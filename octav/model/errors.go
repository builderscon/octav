package model

import "reflect"

func (e ErrInvalidJSONFieldType) Error() string {
	rv := reflect.ValueOf(e.Value)
	var typ string
	if rv == (reflect.Value{}) {
		typ = "invalid"
	} else {
		typ = rv.Type().String()
	}
	return "invalid JSON field type for property '" + e.Field + "' (type: " + typ + ")"
}

func (e ErrInvalidFieldType) Error() string {
	return "invalid field type for property " + e.Field
}
