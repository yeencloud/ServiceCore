package rpc

import (
	"errors"
	"github.com/yeencloud/ServiceCore/domain"
	"net/http"
)

var errVersionMismatch = errors.New("version mismatch")
var ErrVersionMismatch = domain.ServiceError{Code: http.StatusNotAcceptable, Error: errVersionMismatch}

var errUnknownModule = errors.New("unknown module")
var ErrUnknownModule = domain.ServiceError{Code: http.StatusNotFound, Error: errUnknownModule}