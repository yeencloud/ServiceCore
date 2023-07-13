package domain

import (
	"github.com/yeencloud/ServiceCore/serviceError"
)

func CheckRequestVersion(request *ServiceRequest) *serviceError.Error {
	if request.ApiVersion != APIVersion {
		return serviceError.Trace(ErrVersionMismatch)
	}

	return nil
}