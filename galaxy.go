package servicecore

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/yeencloud/ServiceCore/decompose"
)

type RegisterRequest struct {
	Address    string
	Components decompose.Module

	Version Version
}

type RegisterResponse struct {
	Success bool
}

func (sh *ServiceHost) register(host string, port int) error {
	data, err := sh.callWithAddress(sh.Config.GetGalaxyAddress(), sh.Config.GetGalaxyPort(), "Galaxy", "Register", RegisterRequest{
		fmt.Sprintf("%s:%d", host, port),
		*sh.serviceContent,
		APIVersion,
	})

	m, _ := json.Marshal(data)
	var Response RegisterResponse
	_ = json.Unmarshal(m, &Response)

	if err != nil {
		log.Err(err). /*.Strs("services", exported)*/ Msg("Failed to register service(s)")
		return err
	} else {
		log.Info(). /*.Strs("services", exported)*/ Msg("Registered service(s)")
	}
	return nil
}