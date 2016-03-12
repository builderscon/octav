package service

func (e ErrInvalidJSONFieldType) Error() string {
	return "invalid JSON field type for property " +e.Field
}

func (e ErrInvalidFieldType) Error() string {
	return "invalid field type for property " +e.Field
}
