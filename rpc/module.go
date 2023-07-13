package rpc

import (
	"github.com/yeencloud/ServiceCore/domain"
	"github.com/yeencloud/ServiceCore/serviceError"
)

func (rpc *RPC) CheckRequestModule(request *domain.ServiceRequest) *serviceError.Error {
	if request.Module != rpc.Module {
		return serviceError.Trace(ErrUnknownModule)
	}

	return nil
}