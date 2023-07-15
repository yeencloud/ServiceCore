package servicecore

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/yeencloud/ServiceCore/decompose"
	"github.com/yeencloud/ServiceCore/domain"
)

type RegisterRequest struct {
	Address    string
	Components decompose.Module

	Version domain.Version
}

type RegisterResponse struct {
	Success bool
}

func (sh *ServiceHost) register(host string, port int) error {

	registerRequest := RegisterRequest{
		fmt.Sprintf("%s:%d", host, port),
		*sh.serviceContent,
		domain.APIVersion,
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