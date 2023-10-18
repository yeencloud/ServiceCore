package galaxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	serviceError "github.com/yeencloud/ServiceCore/src/adapters/error"
	domain2 "github.com/yeencloud/ServiceCore/src/domain"
	errorDomain "github.com/yeencloud/ServiceCore/src/domain/serviceError"
	"github.com/yeencloud/ServiceCore/src/helpers"
	"io"
	"net/http"
)

func localReply(request domain2.ServiceRequest, err error, code int) domain2.ServiceReply {
	return domain2.ServiceReply{
		RequestID:  request.RequestID,
		Module:     request.Module,
		Method:     request.Method,
		ApiVersion: request.ApiVersion,
		Error: serviceError.Trace(errorDomain.ErrorDescription{
			HttpCode: code,
			String:   err.Error(),
		}),
		ValidationErrors: nil,
		Data:             nil,
	}
}

func localReplyWithDescription(request domain2.ServiceRequest, description *errorDomain.Error) domain2.ServiceReply {
	return domain2.ServiceReply{
		RequestID:  request.RequestID,
		Module:     request.Module,
		Method:     request.Method,
		ApiVersion: request.ApiVersion,
		Error:      description,
		Data:       nil,
	}
}

func (cl *GalaxyClient) CallWithAddress(address string, port int, service string, method string, args any) *domain2.ServiceReply {
	requestURL := fmt.Sprintf("http://%s:%d/rpc/", address, port)

	callData := helpers.AnyToMap(args)

	request := domain2.ServiceRequest{
		Module:     service,
		Method:     method,
		ApiVersion: domain2.APIVersion,
		Data:       callData,
	}

	j, _ := json.Marshal(request)

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(j))
	if err != nil {
		reply := localReply(request, err, http.StatusInternalServerError)
		return &reply //.Embed(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		reply := localReply(request, err, http.StatusInternalServerError)
		return &reply //.Embed(err)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		reply := localReply(request, err, http.StatusInternalServerError)
		return &reply //.Embed(err)
	}

	potentialResponseBody := resBody

	var response domain2.ServiceReply
	err = json.Unmarshal(potentialResponseBody, &response)

	if err != nil {
		reply := localReply(request, err, http.StatusInternalServerError)
		return &reply
	}

	return &response
}

// Call calls a service method with the given arguments (preferably a struct).
func (cl *GalaxyClient) Call(service string, method string, args any) *domain2.ServiceReply {
	var address string
	var port int

	request := domain2.ServiceRequest{
		Module:     service,
		Method:     method,
		ApiVersion: domain2.APIVersion,
		Data:       helpers.AnyToMap(args),
	}

	if service == "Galaxy" {
		address = string(cl.serverConfig.Host)
		port = int(cl.serverConfig.Port)
	} else {
		lookup, err := cl.LookUp(service, method)
		if err != nil {
			reply := localReplyWithDescription(request, err)
			return &reply
		}
		address = lookup.Address
		port = lookup.Port
	}

	return cl.CallWithAddress(address, port, service, method, args)
}
