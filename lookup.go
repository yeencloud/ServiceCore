package servicecore

import (
	"encoding/json"
	"github.com/yeencloud/ServiceCore/domain"
)

type LookUpRequest struct {
	Service string
	Method  string
}

type LookUpResponse struct {
	Address string
}

func (sh *ServiceHost) LookUp(service string, method string) (string, *domain.ServiceError) {
	data, err := sh.Call("Galaxy", "LookUp", LookUpRequest{
		Service: service,
		Method:  method,
	})
	if err != nil {
		return "", err
	}

	var response LookUpResponse
	marshal, _ := json.Marshal(data)
	_ = json.Unmarshal(marshal, &response)
	return response.Address, nil
}