package rpc

import (
	"github.com/yeencloud/ServiceCore/domain"
)

func (rpc *RPC) CheckRequestModule(request *domain.ServiceRequest) *domain.ServiceError {
	if request.Module != rpc.Module {
		return &ErrUnknownModule
	}

	return nil
}