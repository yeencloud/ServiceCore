package servicecore

import (
	"github.com/AliceDiNunno/KubernetesUtil"
	"github.com/yeencloud/ServiceCore/config"
	"github.com/yeencloud/ServiceCore/decompose"
	"github.com/yeencloud/ServiceCore/rpc"
)

type ServiceHost struct {
	service           any
	serviceContent    *decompose.Module
	Config            *config.Config
	ServiceHttpServer *ServiceHTTPServer
	RPC               *rpc.RPC
}

func (sh *ServiceHost) RegisterService(svc any, name string) {
	sh.service = svc
	serviceContent, err := decompose.DecomposeModule(svc, name)
	if err != nil {
		return
	}
	sh.serviceContent = serviceContent
}

func NewServiceHost(service any, modulename string, registerToGalaxy bool) (*ServiceHost, error) {
	s := ServiceHost{}

	decomposed, err := decompose.DecomposeModule(service, modulename)

	if err != nil {
		return nil, err
	}

	rpc := rpc.NewRPC(modulename)

	s.Config = config.NewConfig()
	s.setupLogging()
	s.service = service
	s.serviceContent = decomposed
	s.RPC = &rpc
	s.ServiceHttpServer = newServiceHttpServer(s.Config, s.service, s.serviceContent)
	s.ServiceHttpServer.rpc = s.RPC

	if !registerToGalaxy {
		return &s, nil
	}

	host := "127.0.0.1"
	if KubernetesUtil.IsRunningInKubernetes() {
		host = KubernetesUtil.GetInternalServiceIP()
	}
	port := s.Config.GetRPCPort()

	err = s.register(host, port)

	if err != nil {
		return nil, err
	}

	return &s, nil
}