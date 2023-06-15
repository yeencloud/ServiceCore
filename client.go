package ServiceCore

import (
	"fmt"
	"github.com/AliceDiNunno/KubernetesUtil"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type GalaxyClient struct {
	Version int

	ClientHost string
	ClientPort int

	client *rpc.Client
}

func NewGalaxyClientWithAddress(galaxyAddress string) (*GalaxyClient, error) {
	log.Info().Str("address", galaxyAddress).Msg("Connecting to Galaxy")
	client, err := rpc.DialHTTP("tcp", galaxyAddress)
	if err != nil {
		return nil, err
	}

	return &GalaxyClient{
		Version: 1,
		client:  client,
	}, nil
}

func NewGalaxyClient() (*GalaxyClient, error) {
	nameEnv := ""
	portEnv := ""
	//If we are running in kubernetes, galaxy's service (galaxy) should be exposed as an environment variable by kubernetes
	//Otherwise you can set the environment variables GALAXY_HOST and GALAXY_PORT
	//If you don't set them, it will default to localhost:3000
	if KubernetesUtil.IsRunningInKubernetes() {
		nameEnv = os.Getenv("GALAXY_SERVICE_HOST")
		portEnv = os.Getenv("GALAXY_SERVICE_PORT")
	} else {
		nameEnv = os.Getenv("GALAXY_HOST")
		portEnv = os.Getenv("GALAXY_PORT")

		if nameEnv == "" {
			nameEnv = "localhost"
		}

		if portEnv == "" {
			portEnv = "3000"
		}
	}
	address := fmt.Sprintf("%s:%s", nameEnv, portEnv)

	return NewGalaxyClientWithAddress(address)
}

func PublishMicroService(receiver any, registerToGalaxy bool) {
	err := rpc.Register(receiver)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to register service")
	}

	host := "127.0.0.1"
	if KubernetesUtil.IsRunningInKubernetes() {
		host = KubernetesUtil.GetInternalServiceIP()
	}
	port := GetRPCPort()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal().Err(err).Int("port", port).Msg("Failed to listen on port")
	}
	port = listener.Addr().(*net.TCPAddr).Port

	log.Info().Str("host", host).Int("port", port).Msg("Microservice Listening...")

	if registerToGalaxy {
		galaxy, err := NewGalaxyClient()

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to Galaxy")
		}

		galaxy.RegisterToGalaxy(receiver, host, port)
	}

	rpc.HandleHTTP()
	err = http.Serve(listener, nil)
	if err != nil {
		log.Err(err).Msg("Failed to start HTTP server")
	}
}