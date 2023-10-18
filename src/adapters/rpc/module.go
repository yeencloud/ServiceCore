package rpc

import (
	"github.com/davecgh/go-spew/spew"
	serviceError "github.com/yeencloud/ServiceCore/src/adapters/error"
	"github.com/yeencloud/ServiceCore/src/domain"
	errorDomain "github.com/yeencloud/ServiceCore/src/domain/serviceError"
)

func (rpc *RPC) CheckRequestModule(request *domain.ServiceRequest) *errorDomain.Error {
	spew.Dump(request)
	spew.Dump(rpc)

	if request.Module != rpc.Module {
		return serviceError.Trace(ErrUnknownModule)
	}

	return nil
}
