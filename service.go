package ServiceCore

import (
	"github.com/rs/zerolog/log"
	"github.com/yeencloud/ServiceCore/decompose"
)

type ServiceHost struct {
	service           any
	serviceContent    *decompose.Module
	Config            *Config
	ServiceHttpServer *ServiceHttpServer
	GalaxyClient      *GalaxyClient
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
	s.GalaxyClient, err = newGalaxyClient(s.Config, s.serviceContent)

	if err != nil {
		log.Err(err).Msg("Failed to connect to galaxy")
	}

	return &s
}