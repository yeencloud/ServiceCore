package config

import (
	"github.com/AliceDiNunno/KubernetesUtil"
	types2 "github.com/yeencloud/ServiceCore/src/domain/types"
)

type RpcServer struct {
	Host types2.Host
	Port types2.Port
}

func (config *Config) getRPC() {
	if config.RpcServer != nil {
		return
	}

	var port int
	if KubernetesUtil.IsRunningInKubernetes() {
		port = KubernetesUtil.GetInternalServicePort()
	} else {
		port = config.GetEnvIntOrDefault("RPC_PORT", 0)
	}

	rpcServer := &RpcServer{
		Port: types2.Port(port),
	}

	config.RpcServer = rpcServer
}
