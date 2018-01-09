package server

type ServerError struct {
	Err     string
	Details string
}

// detailed adds details to error
func (e ServerError) detailed(details string) ServerError {
	e.Details = details
	return e
}

func (e ServerError) causedBy(err error) ServerError {
	return e.detailed(err.Error())
}

func (e ServerError) Error() string {
	err := e.Err
	if e.Details != "" {
		err += ": " + e.Details
	}
	return err
}

// e is a shortcut for string error constructor
func e(s string) ServerError {
	return ServerError{Err: s}
}

var (
	errKeyRequired   = e("key query param required")
	errTTLRequired   = e("ttl query param required")
	errIndexRequired = e("index query param required")
	errDKeyRequired  = e("dkey query param required")

	errInvalidTTL   = e("invalid ttl")
	errInvalidIndex = e("invalid index")

	errFailToReadAllBody = e("fail to read all body")

	errInvalidListYAML = e("invalid list yaml")
	errInvalidDictYAML = e("invalid dict yaml")

	errStoreError = e("store error")
)
