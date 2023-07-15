package servicecore

import (
	"github.com/AliceDiNunno/KubernetesUtil"
	"github.com/davecgh/go-spew/spew"
	"github.com/yeencloud/ServiceCore/config"
	"github.com/yeencloud/ServiceCore/decompose"
	"github.com/yeencloud/ServiceCore/rpc"
	"os"
)

type ServiceHost struct {
	service           any
	serviceContent    *decompose.Module
	Config            *config.Config
	ServiceHttpServer *ServiceHTTPServer
	RPC               *rpc.RPC
}

type ServiceClient struct {
	Config *config.Config
}

func (sh *ServiceHost) RegisterService(svc any, name string) {
	sh.service = svc
	serviceContent, err := decompose.DecomposeModule(svc, name)
	if err != nil {
		return
	}
	sh.serviceContent = serviceContent
}

func NewServiceClient() *ServiceHost {
	s := ServiceHost{}
	s.Config = config.NewConfig()
	s.setupLogging()

	return &s
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

	address := "127.0.0.1"
	if KubernetesUtil.IsRunningInKubernetes() {
		address = KubernetesUtil.GetInternalServiceIP()
	}
	port := s.Config.GetRPCPort()

	hostname, _ := os.Hostname()

	spew.Dump("registering", address, port, hostname)
	err = s.register(address, port, hostname)

	if err != nil {
		return nil, err
	}

	return &s, nil
}