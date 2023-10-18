package galaxy

import (
	"encoding/json"
	"errors"
	"github.com/yeencloud/ServiceCore/src/adapters/reflect/decompose"
	"github.com/yeencloud/ServiceCore/src/config"
	"github.com/yeencloud/ServiceCore/src/domain/galaxy"
	errorDomain "github.com/yeencloud/ServiceCore/src/domain/serviceError"
)

type RegisterRequest struct {
	Address    string
	Port       int
	Hostname   string
	Components decompose.Module
}

type RegisterResponse struct {
	Success bool
}

type GalaxyClient struct {
	serverConfig config.GalaxyServer
}

func (gc *GalaxyClient) Register(address string, port int, hostname string, module decompose.Module) error {
	registerRequest := RegisterRequest{
		address,
		port,
		hostname,
		module,
	}
	m, _ := json.Marshal(registerRequest)
	data := gc.CallWithAddress(string(gc.serverConfig.Host), int(gc.serverConfig.Port), "Galaxy", "Register", registerRequest)

	m, _ = json.Marshal(data)
	var Response RegisterResponse
	_ = json.Unmarshal(m, &Response)

	if data.Error != nil && data.Error.HttpCode != 200 {
		return errors.New(data.Error.String)
	}

	return nil
}

func (gc *GalaxyClient) LookUp(service string, method string) (galaxy.LookUpResponse, *errorDomain.Error) {
	data := gc.Call("Galaxy", "LookUp", galaxy.LookUpRequest{
		Service: service,
		Method:  method,
	})

	if data.Error != nil && data.Error.HttpCode != 200 {
		return galaxy.LookUpResponse{}, data.Error
	}

	var response galaxy.LookUpResponse
	marshal, _ := json.Marshal(data.Data)
	_ = json.Unmarshal(marshal, &response)

	return response, nil
}

func NewGalaxyClient(serverConfig config.GalaxyServer) *GalaxyClient {
	return &GalaxyClient{
		serverConfig: serverConfig,
	}
}
