package service

import (
	"github.com/yeencloud/ServiceCore/src/domain"
	galaxyDomain "github.com/yeencloud/ServiceCore/src/domain/galaxy"
	"github.com/yeencloud/ServiceCore/src/domain/serviceError"
)

type Logger interface {
	Debug(message string, fields ...domain.LogField)
	Info(message string, fields ...domain.LogField)
	Warn(message string, fields ...domain.LogField)
	Error(message string, fields ...domain.LogField)
}

type GalaxyClient interface {
	LookUp(service string, method string) (galaxyDomain.LookUpResponse, *serviceError.Error)
	Register(address string, port int, hostname string) error
}

type RPCServer interface {
	Listen() error
}
