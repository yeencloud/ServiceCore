package servicecore

import (
	"encoding/json"
	"errors"
	"fmt"
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
	data, err := sh.callWithAddress(sh.Config.GetGalaxyAddress(), sh.Config.GetGalaxyPort(), "Galaxy", "Register", RegisterRequest{
		fmt.Sprintf("%s:%d", host, port),
		*sh.serviceContent,
		domain.APIVersion,
	})

	m, _ := json.Marshal(data)
	var Response RegisterResponse
	_ = json.Unmarshal(m, &Response)

	if err != nil {
		return errors.New(err.String)
	}
	return nil
}