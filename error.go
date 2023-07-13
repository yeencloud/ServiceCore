package servicecore

import (
	"github.com/yeencloud/ServiceCore/serviceError"
	"net/http"
)

var (
	ErrServiceNotFound            = serviceError.ErrorDescription{HttpCode: http.StatusNotFound, String: "service not found"}
	ErrMethodNotFound             = serviceError.ErrorDescription{HttpCode: http.StatusNotFound, String: "method not found"}
	ErrCouldNotGetMethodParameter = serviceError.ErrorDescription{HttpCode: http.StatusInternalServerError, String: "could not get method parameter"}
	ErrRequestDataIsMissing       = serviceError.ErrorDescription{HttpCode: http.StatusBadRequest, String: "request data is missing"}
	ErrInvalidRequestID           = serviceError.ErrorDescription{HttpCode: http.StatusBadRequest, String: "request id is invalid"}

	ErrCallCouldNotMakeHttpRequest       = serviceError.ErrorDescription{HttpCode: http.StatusInternalServerError, String: "could not make http request"}
	ErrCallCouldNotReadResponseBody      = serviceError.ErrorDescription{HttpCode: http.StatusInternalServerError, String: "could not read response body"}
	ErrCallCouldNotUnmarshalResponseBody = serviceError.ErrorDescription{HttpCode: http.StatusInternalServerError, String: "could not unmarshal response body"}

	ErrRequestCouldNotBindRequest = serviceError.ErrorDescription{HttpCode: http.StatusBadRequest, String: "could not bind request because it is missing"}
	ErrRequestCouldNotBeCast      = serviceError.ErrorDescription{HttpCode: http.StatusBadRequest, String: "could not cast request to service request"}

	ErrValidationFailed = serviceError.ErrorDescription{HttpCode: http.StatusBadRequest, String: "validation failed"}
)