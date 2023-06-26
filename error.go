package servicecore

import "errors"

var (
	ErrVersionMismatch  = errors.New("version mismatch")
	ErrRequestIsMissing = errors.New("request is missing")
	ErrUnknownMethod    = errors.New("unknown method")
	ErrorUnknownModule  = errors.New("unknown module")
	ErrRequestMalformed = errors.New("request is malformed")
	ErrMethodIsInvalid  = errors.New("method is invalid")
	ErrServiceIsInvalid = errors.New("service is invalid")
	ErrValidationFailed = errors.New("validation failed")
)