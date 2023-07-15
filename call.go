package servicecore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/yeencloud/ServiceCore/domain"
	"github.com/yeencloud/ServiceCore/serviceError"
	"github.com/yeencloud/ServiceCore/tools"
	"io/ioutil"
	"net/http"
	"os"
)

func (sh *ServiceHost) callWithAddress(hostname string, port int, service string, method string, args any) (map[string]interface{}, *serviceError.Error) {
	requestURL := fmt.Sprintf("http://%s:%d/rpc/", hostname, port)

	callData := tools.AnyToMap(args)

	j, _ := json.Marshal(domain.ServiceRequest{
		Module:     service,
		Method:     method,
		ApiVersion: domain.APIVersion,
		Data:       callData,
	})

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(j))
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, serviceError.Trace(ErrCallCouldNotReadResponseBody) //.Embed(err)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, serviceError.Trace(ErrCallCouldNotReadResponseBody)
	}

	var response map[string]interface{}
	err = json.Unmarshal(resBody, &response)

	if err != nil {
		return nil, serviceError.Trace(ErrCallCouldNotUnmarshalResponseBody)
	}

	return response, nil
}

// Call calls a service method with the given arguments (preferably a struct).
func (sh *ServiceHost) Call(service string, method string, args any) (map[string]interface{}, *serviceError.Error) {
	var address string
	var port int

	if service == "Galaxy" {
		address = sh.Config.GetGalaxyAddress()
		port = sh.Config.GetGalaxyPort()
	} else {
		lookup, err := sh.LookUp(service, method)
		address = lookup.Address
		if err != nil {
			return nil, err
		}
		port = lookup.Port
	}

	return sh.callWithAddress(address, port, service, method, args)
}