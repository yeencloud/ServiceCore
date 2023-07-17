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
)

func localReply(request domain.ServiceRequest, err error, code int) domain.ServiceReply {
	return domain.ServiceReply{
		RequestID:  request.RequestID,
		Module:     request.Module,
		Method:     request.Method,
		ApiVersion: request.ApiVersion,
		Error: serviceError.Trace(serviceError.ErrorDescription{
			HttpCode: code,
			String:   err.Error(),
		}),
		ValidationErrors: nil,
		Data:             nil,
	}
}

func localReplyWithDescription(request domain.ServiceRequest, description *serviceError.Error) domain.ServiceReply {
	return domain.ServiceReply{
		RequestID:  request.RequestID,
		Module:     request.Module,
		Method:     request.Method,
		ApiVersion: request.ApiVersion,
		Error:      description,
		Data:       nil,
	}
}

func (sh *ServiceHost) callWithAddress(address string, port int, service string, method string, args any) *domain.ServiceReply {
	requestURL := fmt.Sprintf("http://%s:%d/rpc/", address, port)

	callData := tools.AnyToMap(args)

	request := domain.ServiceRequest{
		Module:     service,
		Method:     method,
		ApiVersion: domain.APIVersion,
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

	var response domain.ServiceReply
	err = json.Unmarshal(potentialResponseBody, &response)

	if err != nil {
		reply := localReply(request, err, http.StatusInternalServerError)
		return &reply
	}

	return &response
}

// Call calls a service method with the given arguments (preferably a struct).
func (sh *ServiceHost) Call(service string, method string, args any) *domain.ServiceReply {
	var address string
	var port int

	request := domain.ServiceRequest{
		Module:     service,
		Method:     method,
		ApiVersion: domain.APIVersion,
		Data:       tools.AnyToMap(args),
	}

	if service == "Galaxy" {
		address = sh.Config.GetGalaxyAddress()
		port = sh.Config.GetGalaxyPort()
	} else {
		lookup, err := sh.LookUp(service, method)
		if err != nil {
			reply := localReplyWithDescription(request, err)
			return &reply
		}
		address = lookup.Address
		port = lookup.Port
	}

	return sh.callWithAddress(address, port, service, method, args)
}