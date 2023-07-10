package servicecore

import (
	"errors"
	"github.com/yeencloud/ServiceCore/domain"
	"net/http"
)

var (
	ErrServiceNotFound            = domain.ServiceError{Code: http.StatusNotFound, Error: errors.New("service not found")}
	ErrMethodNotFound             = domain.ServiceError{Code: http.StatusNotFound, Error: errors.New("method not found")}
	ErrCouldNotGetMethodParameter = domain.ServiceError{Code: http.StatusInternalServerError, Error: errors.New("could not get method parameter")}
	ErrRequestDataIsMissing       = domain.ServiceError{Code: http.StatusBadRequest, Error: errors.New("request data is missing")}

	ErrCallCouldNotMakeHttpRequest       = domain.ServiceError{Code: http.StatusInternalServerError, Error: errors.New("could not make http request")}
	ErrCallCouldNotReadResponseBody      = domain.ServiceError{Code: http.StatusInternalServerError, Error: errors.New("could not read response body")}
	ErrCallCouldNotUnmarshalResponseBody = domain.ServiceError{Code: http.StatusInternalServerError, Error: errors.New("could not unmarshal response body")}

	ErrRequestCouldNotBindRequest = domain.ServiceError{Code: http.StatusBadRequest, Error: errors.New("could not bind request because it is missing")}
	ErrRequestCouldNotBeCast      = domain.ServiceError{Code: http.StatusBadRequest, Error: errors.New("could not cast request to service request")}

	ErrValidationFailed = domain.ServiceError{Code: http.StatusBadRequest, Error: errors.New("validation failed")}
)