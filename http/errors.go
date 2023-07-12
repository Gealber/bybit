package http

import "errors"

var (
	ErrorUnexpectedStatus = errors.New("unexpected http status code")
	ErrorUnavailableInformation = errors.New("unavailable information")
	ErrorInsuficcientBalance = errors.New("insuficcient balance")
)
