package servicecore

import (
	"encoding/json"
	"errors"
	"github.com/yeencloud/ServiceCore/decompose"
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

func (sh *ServiceHost) register(address string, port int, hostname string) error {
	registerRequest := RegisterRequest{
		address,
		port,
		hostname,
		*sh.serviceContent,
	}
	m, _ := json.Marshal(registerRequest)
	data := sh.callWithAddress(sh.Config.GetGalaxyAddress(), sh.Config.GetGalaxyPort(), "Galaxy", "Register", registerRequest)

	m, _ = json.Marshal(data)
	var Response RegisterResponse
	_ = json.Unmarshal(m, &Response)

	if data.Error != nil && data.Error.HttpCode != 200 {
		return errors.New(data.Error.String)
	}

	return nil
}