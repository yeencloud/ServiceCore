package rpc

import "github.com/yeencloud/ServiceCore/domain"

func CheckRequestVersion(request *domain.ServiceRequest) *domain.ServiceError {
	if request.ApiVersion != domain.APIVersion {
		return &ErrVersionMismatch
	}

	return nil
}