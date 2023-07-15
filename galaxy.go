package servicecore

import (
	"encoding/json"
	"errors"
	"github.com/davecgh/go-spew/spew"
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
	spew.Dump("registering", address, port, hostname)

	registerRequest := RegisterRequest{
		address,
		port,
		hostname,
		*sh.serviceContent,
	}

	spew.Dump("register request", registerRequest)
	m, _ := json.Marshal(registerRequest)
	spew.Dump(m)

	data, err := sh.callWithAddress(sh.Config.GetGalaxyAddress(), sh.Config.GetGalaxyPort(), "Galaxy", "Register", registerRequest)

	m, _ = json.Marshal(data)
	var Response RegisterResponse
	_ = json.Unmarshal(m, &Response)

	spew.Dump("register reply", m)

	if err != nil {
		return errors.New(err.String)
	}
	return nil
}