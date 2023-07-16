package servicecore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/yeencloud/ServiceCore/domain"
	"github.com/yeencloud/ServiceCore/serviceError"
	"github.com/yeencloud/ServiceCore/tools"
	"io"
	"net/http"
	"os"
)

func (sh *ServiceHost) callWithAddress(address string, port int, service string, method string, args any) (*domain.ServiceReply, *serviceError.Error) {
	requestURL := fmt.Sprintf("http://%s:%d/rpc/", address, port)

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

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, serviceError.Trace(ErrCallCouldNotReadResponseBody)
	}

	potentialErrorBody := resBody
	potentialResponseBody := potentialErrorBody

	var serverr serviceError.Error
	err = json.Unmarshal(potentialErrorBody, &serverr)

	if err == nil && serverr.HttpCode >= 300 {
		return nil, &serverr
	}

	var response domain.ServiceReply
	err = json.Unmarshal(potentialResponseBody, &response)

	if err != nil {
		return nil, serviceError.Trace(ErrCallCouldNotUnmarshalResponseBody)
	}

	return &response, nil
}

// Call calls a service method with the given arguments (preferably a struct).
func (sh *ServiceHost) Call(service string, method string, args any) (*domain.ServiceReply, *serviceError.Error) {
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

	return sh.callWithAddress(address, port, service, method, args)
}