package servicecore

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
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
	data, err := sh.Call("Galaxy", "LookUp", LookUpRequest{
		Service: service,
		Method:  method,
	})

	if err != nil {
		spew.Dump("Call error", err)
		return LookUpResponse{}, err
	}

	spew.Dump("LookUp response body", data)

	var response LookUpResponse
	marshal, _ := json.Marshal(data.Data)
	_ = json.Unmarshal(marshal, &response)

	spew.Dump("LookUp response", response)

	return response, nil
}