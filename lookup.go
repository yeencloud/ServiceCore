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
	Port    int
}

func (sh *ServiceHost) LookUp(service string, method string) (LookUpResponse, *serviceError.Error) {
	data := sh.Call("Galaxy", "LookUp", LookUpRequest{
		Service: service,
		Method:  method,
	})

	if data.Error != nil && data.Error.HttpCode != 200 {
		return LookUpResponse{}, data.Error
	}

	var response LookUpResponse
	marshal, _ := json.Marshal(data.Data)
	_ = json.Unmarshal(marshal, &response)

	return response, nil
}