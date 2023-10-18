package service

import (
	"github.com/AliceDiNunno/KubernetesUtil"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/yeencloud/ServiceCore/src/adapters/environment"
	"github.com/yeencloud/ServiceCore/src/adapters/galaxy"
	"github.com/yeencloud/ServiceCore/src/adapters/http"
	"github.com/yeencloud/ServiceCore/src/adapters/logging/zerolog"
	"github.com/yeencloud/ServiceCore/src/adapters/persistence/postgres"
	"github.com/yeencloud/ServiceCore/src/adapters/reflect/decompose"
	"github.com/yeencloud/ServiceCore/src/adapters/rpc"
	"github.com/yeencloud/ServiceCore/src/config"
	"os"
	"reflect"
)

type Service struct {
	rpcServer RPCServer
	instance  any

	Name string
}

type ServiceClient struct {
	Config *config.Config
}

func (sh *Service) RegisterService(svc any, name string) {
	_, err := decompose.DecomposeModule(svc, name)
	if err != nil {
		return
	}
}

type ServiceConfig struct {
	ModuleName       string
	RegisterToGalaxy bool
	DatabaseItems    []interface{}
}

func processDatabaseItems(items []interface{}) {
	//for each item check if it has a field called ServiceDatabaseEntity

	for _, item := range items {
		itemtype := reflect.TypeOf(item)

		field, found := itemtype.FieldByName("ServiceDatabaseEntity")

		if !found {
			log.Warn().Msg("No ServiceDatabaseEntity field found in item " + itemtype.Name())
			return
		}

		if field.Type != reflect.TypeOf(postgres.ServiceDatabaseEntity{}) {
			log.Warn().Msg("A field called ServiceDatabaseEntity was found in item " + itemtype.Name() + " but it is not of type postgres.ServiceDatabaseEntity")
		}
	}
}

func NewServiceClient() *Service {
	s := Service{}

	return &s
}

func NewServiceHost(instance any, serviceConfig ServiceConfig) (*Service, error) {
	_ = godotenv.Load()

	s := Service{}

	_ = zerolog.NewLogger(environment.GetEnvironment())

	config := config.NewConfig()

	s.Name = serviceConfig.ModuleName

	moduleContent, err := decompose.DecomposeModule(instance, s.Name)
	if err != nil {
		return nil, err
	}

	_ = rpc.NewRPC(serviceConfig.ModuleName)

	if serviceConfig.DatabaseItems != nil {
		postgres.StartGormDatabase(config.Database, serviceConfig.ModuleName)
		processDatabaseItems(serviceConfig.DatabaseItems)
		//	s.Database.Migrate()
	}

	rpcServer := http.NewServiceHttpServer(*config.RpcServer, instance, moduleContent)
	s.rpcServer = rpcServer

	galaxyClient := galaxy.NewGalaxyClient(*config.GalaxyServer)
	if !serviceConfig.RegisterToGalaxy {
		return &s, nil
	}
	//

	serviceAddress := "127.0.0.1"
	if KubernetesUtil.IsRunningInKubernetes() {
		serviceAddress = KubernetesUtil.GetInternalServiceIP()
	}
	port := config.RpcServer.Port

	hostname, _ := os.Hostname()

	err = galaxyClient.Register(serviceAddress, int(port), hostname, *moduleContent)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (sh *Service) Listen() error {
	return sh.rpcServer.Listen()
}
