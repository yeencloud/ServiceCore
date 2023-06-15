package ServiceCore

import (
	"fmt"
	"github.com/AliceDiNunno/KubernetesUtil"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

type RegisterRequest struct {
	Address    string
	Components []string

	Version int
}

type RegisterResponse struct {
	Success bool
}

func GetRPCPort() int {
	port := 0
	if KubernetesUtil.IsRunningInKubernetes() {
		port = KubernetesUtil.GetInternalServicePort()
	} else {
		portStr := os.Getenv("RPC_PORT")
		if portStr == "" {
			log.Warn().Msg("RPC_PORT not set, defaulting to a random available port")
			portStr = "0"
		}
		envport, err := strconv.Atoi(portStr)
		if err != nil {
			log.Err(err).Str("port", portStr).Msg("RPC_PORT is invalid, defaulting to a random available port")
			envport = 0
		}
		port = envport
	}

	return port
}

func (galaxy *GalaxyClient) RegisterToGalaxy(service any, host string, port int) {
	var response RegisterResponse

	galaxy.ClientHost = host
	galaxy.ClientPort = port

	exported := ExportList(service)

	err := galaxy.client.Call("Galaxy.Register", RegisterRequest{
		fmt.Sprintf("%s:%d", host, port),
		exported,
		galaxy.Version,
	}, &response)
	if err != nil {
		log.Err(err).Strs("services", exported).Msg("Failed to register service(s)")
	} else {
		log.Info().Strs("services", exported).Msg("Registered service(s)")
	}
}