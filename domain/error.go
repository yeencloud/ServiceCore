package domain

import (
	"github.com/yeencloud/ServiceCore/serviceError"
	"net/http"
)

var ErrVersionMismatch = serviceError.ErrorDescription{HttpCode: http.StatusNotAcceptable, String: "version mismatch"}