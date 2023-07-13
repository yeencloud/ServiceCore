package servicecore

import (
	"encoding/json"
	"github.com/yeencloud/ServiceCore/serviceError"
)

type LookUpRequest struct {
	Service string
	Method  string
}

type LookUpResponse struct {
	Address string
}

func (sh *ServiceHost) LookUp(service string, method string) (string, *serviceError.Error) {
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