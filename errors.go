package togglr

import (
	"errors"
)

var (
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrTooManyRequests     = errors.New("too many requests")
	ErrInvalidConfig       = errors.New("invalid configuration")
	ErrFeatureNotFound     = errors.New("feature not found")
	ErrBadRequest          = errors.New("bad request")
	ErrInternalServerError = errors.New("internal server error")
)

type APIError struct {
	Code       string
	Message    string
	StatusCode int
	Err        error
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return e.Message
}

func (e *APIError) Unwrap() error {
	return e.Err
}

func IsAPIError(err error) bool {
	var apiErr *APIError
	ok := errors.As(err, &apiErr)

	return ok
}

func GetAPIErrorCode(err error) (string, bool) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code, true
	}

	return "", false
}
