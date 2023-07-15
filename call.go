package servicecore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/yeencloud/ServiceCore/domain"
	"github.com/yeencloud/ServiceCore/serviceError"
	"github.com/yeencloud/ServiceCore/tools"
	"io"
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

	spew.Dump("Call request", string(j))

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
	spew.Dump("Call response", res)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		spew.Dump("Call error", err)
		return nil, serviceError.Trace(ErrCallCouldNotReadResponseBody)
	}
	spew.Dump("Call response body", string(resBody))

	potentialErrorBody := resBody
	potentialResponseBody := potentialErrorBody

	spew.Dump("Call potential error body", string(potentialErrorBody))
	spew.Dump("Call potential response body", string(potentialResponseBody))

	var serverr serviceError.Error
	err = json.Unmarshal(potentialErrorBody, &serverr)

	if err != nil {
		spew.Dump("Call error 62", err)
	}

	if err == nil && serverr.HttpCode >= 300 {
		spew.Dump("Call error 65", serverr)
		return nil, &serverr
	}

	var response map[string]interface{}
	err = json.Unmarshal(potentialResponseBody, &response)

	if err != nil {
		spew.Dump("Call error 71", err.Error())
	}

	if err != nil {
		spew.Dump("Call error", err)
		return nil, serviceError.Trace(ErrCallCouldNotUnmarshalResponseBody)
	}
	spew.Dump("Call response", response)

	return response, nil
}

// Call calls a service method with the given arguments (preferably a struct).
func (sh *ServiceHost) Call(service string, method string, args any) (map[string]interface{}, *serviceError.Error) {
	var address string
	var port int

	if service == "Galaxy" {
		address = sh.Config.GetGalaxyAddress()
		port = sh.Config.GetGalaxyPort()

		spew.Dump("Galaxy address", address, "port", port)
	} else {
		lookup, err := sh.LookUp(service, method)
		if err != nil {
			return nil, err
		}
		address = lookup.Address
		port = lookup.Port

		spew.Dump("lookup", lookup)
	}

	spew.Dump("service", service, "method", method, " is at ", address, ":", port)

	return sh.callWithAddress(address, port, service, method, args)
}