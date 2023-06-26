package servicecore

import "encoding/json"

type LookUpRequest struct {
	Service string
	Method  string
}

type LookUpResponse struct {
	Address string
}

func (sh *ServiceHost) LookUp(service string, method string) (string, error) {
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