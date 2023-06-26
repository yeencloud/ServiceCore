package servicecore

import (
	"fmt"
	"github.com/AliceDiNunno/KubernetesUtil"
	"github.com/rs/zerolog/log"
	"github.com/yeencloud/ServiceCore/decompose"
	"os"
)

type ServiceHost struct {
	service           any
	serviceContent    *decompose.Module
	Config            *Config
	ServiceHttpServer *ServiceHTTPServer
}

func (sh *ServiceHost) RegisterService(svc any, name string) {
	sh.service = svc
	serviceContent, err := decompose.DecomposeModule(svc, name)
	if err != nil {
		return
	}
	sh.serviceContent = serviceContent
}

func NewServiceHost(service any, name string) *ServiceHost {
	s := ServiceHost{}

	decomposed, err := decompose.DecomposeModule(service, name)

	if err != nil {
		return nil
	}

	s.Config = newConfig()
	s.setupLogging()
	s.service = service
	s.serviceContent = decomposed
	s.ServiceHttpServer = newServiceHttpServer(s.Config, s.service, s.serviceContent)

	//If the service is galaxy, we don't want to connect it to itself
	if name == "Galaxy" {
		return &s
	}

	host := "127.0.0.1"
	if KubernetesUtil.IsRunningInKubernetes() {
		host = KubernetesUtil.GetInternalServiceIP()
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}
	port := s.Config.GetRPCPort()

	err = s.register(host, port)

	if err != nil {
		log.Err(err).Msg("Failed to register service(s) to Galaxy")
		return nil
	}

	return &s
}