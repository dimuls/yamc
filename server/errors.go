package server

import "errors"

func e(s string) error {
	return errors.New(s)
}

var (
	errKeyRequired     = e("key query param required")
	errTTLRrequired    = e("ttl query param required")
	errInvalidTTL      = e("err invalid ttl")
	errInvalidListYAML = e("invalid list yaml")
	errInvalidDictYAML = e("invalid dict yaml")
)
