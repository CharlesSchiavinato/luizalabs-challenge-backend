package usecase

// ErrParamValidate denotes failing validate param.
type ErrParamValidate struct {
	Message string
}

// ErrParamValidate returns the param validation error.
func (epv ErrParamValidate) Error() string {
	return epv.Message
}

// ErrRecordValidate denotes failing validate record.
type ErrRecordValidate struct {
	Message string
}

// ErrRecordValidate returns the record validation error.
func (erv ErrRecordValidate) Error() string {
	return erv.Message
}
