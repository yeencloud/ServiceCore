package ServiceCore

import (
	"fmt"
	"github.com/AliceDiNunno/KubernetesUtil"
	"github.com/rs/zerolog/log"
	"github.com/yeencloud/ServiceCore/decompose"
	"net"
	"net/http"
	"net/rpc"
)

type RegisterRequest struct {
	Address    string
	Components []decompose.Module

	Version int
}

type RegisterResponse struct {
	Success bool
}

type GalaxyClient struct {
	Version int

	ClientHost string
	ClientPort int
}

func newGalaxyClientWithAddress(galaxyAddress string) (*GalaxyClient, error) {
	/*log.Info().Str("address", galaxyAddress).Msg("Connecting to Galaxy")
	client, err := rpc.DialHTTP("tcp", galaxyAddress)
	if err != nil {
		return nil, err
	}

	return &GalaxyClient{
		Version: 1,
		client:  client,
	}, nil*/

	return nil, nil
}

func newGalaxyClient(config *Config, module *decompose.Module) (*GalaxyClient, error) {
	nameEnv := config.GetGalaxyAddress()
	portEnv := config.GetGalaxyPort()

	address := fmt.Sprintf("%s:%s", nameEnv, portEnv)

	return newGalaxyClientWithAddress(address)
}

func (galaxy *GalaxyClient) register(host string, port int) error {
	var response RegisterResponse

	galaxy.ClientHost = host
	galaxy.ClientPort = port

	err := galaxy.client.Call("Galaxy.Register", RegisterRequest{
		fmt.Sprintf("%s:%d", host, port),
		exported,
		galaxy.Version,
	}, &response)
	if err != nil {
		log.Err(err).Strs("services", exported).Msg("Failed to register service(s)")
		return err
	} else {
		log.Info().Strs("services", exported).Msg("Registered service(s)")
	}
	return nil
}

func PublishMicroService(receiver any, registerToGalaxy bool) error {
	err := rpc.Register(receiver)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to register service")
		return err
	}

	host := "127.0.0.1"
	if KubernetesUtil.IsRunningInKubernetes() {
		host = KubernetesUtil.GetInternalServiceIP()
	}
	port := GetRPCPort()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal().Err(err).Int("port", port).Msg("Failed to listen on port")
		return err
	}
	port = listener.Addr().(*net.TCPAddr).Port

	log.Info().Str("host", host).Int("port", port).Msg("Microservice Listening...")

	if registerToGalaxy {
		galaxy, err := NewGalaxyClient()

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to Galaxy")
			return err
		}

		galaxy.RegisterToGalaxy(receiver, host, port)
	}

	rpc.HandleHTTP()
	err = http.Serve(listener, nil)
	if err != nil {
		log.Err(err).Msg("Failed to start HTTP server")
		return err
	}

	return nil
}