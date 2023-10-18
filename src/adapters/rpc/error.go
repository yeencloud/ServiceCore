package rpc

import (
	"github.com/yeencloud/ServiceCore/src/domain/serviceError"
	"net/http"
)

var ErrUnknownModule = serviceError.ErrorDescription{HttpCode: http.StatusNotFound, String: "unknown module"}
