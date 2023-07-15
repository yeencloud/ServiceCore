package servicecore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/yeencloud/ServiceCore/domain"
	"github.com/yeencloud/ServiceCore/serviceError"
	"github.com/yeencloud/ServiceCore/tools"
	"io/ioutil"
	"net/http"
	"os"
)

func (sh *ServiceHost) callWithAddress(address string, port int, service string, method string, args any) (map[string]interface{}, *serviceError.Error) {
	requestURL := fmt.Sprintf("http://%s:%d/rpc/", address, port)

	spew.Dump("Call at ", requestURL)

	callData := tools.AnyToMap(args)

	j, _ := json.Marshal(domain.ServiceRequest{
		Module:     service,
		Method:     method,
		ApiVersion: domain.APIVersion,
		Data:       callData,
	})

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(j))
	if err != nil {
		spew.Dump("Call error", err)
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		spew.Dump("Call error", err)
		return nil, serviceError.Trace(ErrCallCouldNotReadResponseBody) //.Embed(err)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		spew.Dump("Call error", err)
		return nil, serviceError.Trace(ErrCallCouldNotReadResponseBody)
	}

	var response map[string]interface{}
	err = json.Unmarshal(resBody, &response)

	if err != nil {
		spew.Dump("Call error", err)
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
		if err != nil {
			return nil, err
		}
		address = lookup.Address
		port = lookup.Port
	}

	spew.Dump("service", service, "method", method, " is at ", address, ":", port)

	return sh.callWithAddress(address, port, service, method, args)
}